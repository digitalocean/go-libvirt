// Copyright 2016 The go-libvirt Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package libvirt is a pure Go implementation of the libvirt RPC protocol.
// For more information on the protocol, see https://libvirt.org/internals/l.html
package libvirt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"

	"github.com/davecgh/go-xdr/xdr2"
	"github.com/digitalocean/go-libvirt/internal/constants"
)

// ErrEventsNotSupported is returned by Events() if event streams
// are unsupported by either QEMU or libvirt.
var ErrEventsNotSupported = errors.New("event monitor is not supported")

// Libvirt implements LibVirt's remote procedure call protocol.
type Libvirt struct {
	conn net.Conn
	r    *bufio.Reader
	w    *bufio.Writer

	// method callbacks
	cm        sync.Mutex
	callbacks map[uint32]chan response

	// event listeners
	em     sync.Mutex
	events map[uint32]chan *DomainEvent

	// next request serial number
	s uint32
}

// Domain represents a domain as seen by libvirt.
type Domain struct {
	Name string
	UUID [constants.UUIDSize]byte
	ID   int
}

// DomainEvent represents a libvirt domain event.
type DomainEvent struct {
	CallbackID   uint32
	Domain       Domain
	Event        string
	Seconds      uint64
	Microseconds uint32
	Padding      uint8
	Details      []byte
}

// Secret represents a secret managed by the libvirt daemon.
type Secret struct {
	UUID      [constants.UUIDSize]byte
	UsageType SecretUsageType
	UsageID   string
}

// StoragePool represents a storage pool as seen by libvirt.
type StoragePool struct {
	Name string
	UUID [constants.UUIDSize]byte
}

// qemuError represents a QEMU process error.
type qemuError struct {
	Error struct {
		Class       string `json:"class"`
		Description string `json:"desc"`
	} `json:"error"`
}

// DomainXMLFlags specifies options for dumping a domain's XML.
type DomainXMLFlags uint32

// DomainAffectFlags specifies options for whether an operation affects the
// running VM, or the persistent VM configuration on disk. See FlagDomain...
// consts for values.
type DomainAffectFlags uint32

// Consts used for flags

// virDomainModificationImpact and virTypedParameterFlags values. These are
// combined here because they are both used to set the same flags fields in the
// libvirt API.
const (
	// FlagDomainAffectCurrent means affect the current domain state
	FlagDomainAffectCurrent DomainAffectFlags = 0
	// FlagDomainAffectLive means affect the running domain state
	FlagDomainAffectLive = 1 << (iota - 1)
	// FlagDomainAffectConfig means affect the persistent domain state.
	FlagDomainAffectConfig
	// FlagTypedParamStringOkay tells the server that this client understands
	// TypedParamStrings.
	FlagTypedParamStringOkay
)

// Consts relating to TypedParams:
const (
	// TypeParamInt is a C int.
	TypeParamInt = iota + 1
	// TypeParamUInt is a C unsigned int.
	TypeParamUInt
	// TypeParamLLong is a C long long int.
	TypeParamLLong
	// TypeParamULLong is a C unsigned long long int.
	TypeParamULLong
	// TypeParamDouble is a C double.
	TypeParamDouble
	// TypeParamBoolean is a C char.
	TypeParamBoolean
	// TypeParamString is a C char*.
	TypeParamString

	// TypeParamLast is just an end-of-enum marker.
	TypeParamLast
)

// TypedParam represents libvirt's virTypedParameter, which is used to pass
// typed parameters to libvirt functions. This is defined as an interface, and
// implemented by TypedParam... concrete types.
type TypedParam interface {
	Field() string
	Value() interface{}
}

// TypedParamInt contains a 32-bit signed integer.
type TypedParamInt struct {
	Fld     string
	PType   int32
	Padding [4]byte
	Val     int32
}

// Field returns the field name, a string name for the parameter.
func (tp *TypedParamInt) Field() string {
	return tp.Fld
}

// Value returns the value for the typed parameter, as an empty interface.
func (tp *TypedParamInt) Value() interface{} {
	return tp.Val
}

// NewTypedParamInt returns a TypedParam encoding for an int.
func NewTypedParamInt(name string, v int32) *TypedParamInt {
	// Truncate the field name if it's longer than the limit.
	if len(name) > constants.TypedParamFieldLength {
		name = name[:constants.TypedParamFieldLength]
	}
	tp := TypedParamInt{
		Fld:   name,
		PType: TypeParamInt,
		Val:   v,
	}
	return &tp
}

// TypedParamULongLong contains a 64-bit unsigned integer.
type TypedParamULongLong struct {
	Fld   string
	PType int32
	Val   uint64
}

// Field returns the field name, a string name for the parameter.
func (tp *TypedParamULongLong) Field() string {
	return tp.Fld
}

// Value returns the value for the typed parameter, as an empty interface.
func (tp *TypedParamULongLong) Value() interface{} {
	return tp.Val
}

// NewTypedParamULongLong returns a TypedParam encoding for an unsigned long long.
func NewTypedParamULongLong(name string, v uint64) *TypedParamULongLong {
	// Truncate the field name if it's longer than the limit.
	if len(name) > constants.TypedParamFieldLength {
		name = name[:constants.TypedParamFieldLength]
	}
	tp := TypedParamULongLong{
		Fld:   name,
		PType: TypeParamULLong,
		Val:   v,
	}
	return &tp
}

// TypedParamString contains a string parameter.
type TypedParamString struct {
	Fld   string
	PType int
	Val   string
}

// Field returns the field name, a string name for the parameter.
func (tp *TypedParamString) Field() string {
	return tp.Fld
}

// Value returns the value for the typed parameter, as an empty interface.
func (tp *TypedParamString) Value() interface{} {
	return tp.Val
}

// NewTypedParamString returns a typed parameter containing the passed string.
func NewTypedParamString(name string, v string) TypedParam {
	if len(name) > constants.TypedParamFieldLength {
		name = name[:constants.TypedParamFieldLength]
	}
	tp := TypedParamString{
		Fld:   name,
		PType: TypeParamString,
		Val:   v,
	}
	return &tp
}

const (
	// DomainXMLFlagSecure dumps XML with sensitive information included.
	DomainXMLFlagSecure DomainXMLFlags = 1 << iota

	// DomainXMLFlagInactive dumps XML with inactive domain information.
	DomainXMLFlagInactive

	// DomainXMLFlagUpdateCPU dumps XML with guest CPU requirements according to the host CPU.
	DomainXMLFlagUpdateCPU

	// DomainXMLFlagMigratable dumps XML suitable for migration.
	DomainXMLFlagMigratable
)

// MigrateFlags specifies options when performing a migration.
type MigrateFlags uint32

const (
	// MigrateFlagLive performs a zero-downtime live migration.
	MigrateFlagLive MigrateFlags = 1 << iota

	// MigrateFlagPeerToPeer creates a direct source to destination control channel.
	MigrateFlagPeerToPeer

	// MigrateFlagTunneled tunnels migration data over the libvirtd connection.
	MigrateFlagTunneled

	// MigrateFlagPersistDestination will persist the VM on the destination host.
	MigrateFlagPersistDestination

	// MigrateFlagUndefineSource undefines the VM on the source host.
	MigrateFlagUndefineSource

	// MigrateFlagPaused will pause the remote side VM.
	MigrateFlagPaused

	// MigrateFlagNonSharedDisk migrate non-shared storage with full disk copy.
	MigrateFlagNonSharedDisk

	// MigrateFlagNonSharedIncremental migrate non-shared storage with incremental copy.
	MigrateFlagNonSharedIncremental

	// MigrateFlagChangeProtection prevents any changes to the domain configuration through the whole migration process.
	MigrateFlagChangeProtection

	// MigrateFlagUnsafe will force a migration even when it is considered unsafe.
	MigrateFlagUnsafe

	// MigrateFlagOffline is used to perform an offline migration.
	MigrateFlagOffline

	// MigrateFlagCompressed compresses data during migration.
	MigrateFlagCompressed

	// MigrateFlagAbortOnError will abort a migration on I/O errors encountered during migration.
	MigrateFlagAbortOnError

	// MigrateFlagAutoConverge forces convergence.
	MigrateFlagAutoConverge

	// MigrateFlagRDMAPinAll enables RDMA memory pinning.
	MigrateFlagRDMAPinAll
)

// UndefineFlags specifies options available when undefining a domain.
type UndefineFlags uint32

const (
	// UndefineFlagManagedSave removes all domain managed save data.
	UndefineFlagManagedSave UndefineFlags = 1 << iota

	// UndefineFlagSnapshotsMetadata removes all domain snapshot metadata.
	UndefineFlagSnapshotsMetadata

	// UndefineFlagNVRAM removes all domain NVRAM files.
	UndefineFlagNVRAM
)

// DomainDefineXMLFlags specifies options available when defining a domain.
type DomainDefineXMLFlags uint32

const (
	// DefineValidate validates the XML document against schema
	DefineValidate DomainDefineXMLFlags = 1
)

// DestroyFlags specifies options available when destroying a domain.
type DestroyFlags uint32

const (
	// DestroyFlagDefault default behavior, forcefully terminate the domain.
	DestroyFlagDefault DestroyFlags = 1 << iota

	// DestroyFlagGraceful only sends a SIGTERM no SIGKILL.
	DestroyFlagGraceful
)

// ShutdownFlags specifies options available when shutting down a domain.
type ShutdownFlags uint32

const (
	// ShutdownAcpiPowerBtn - send ACPI event
	ShutdownAcpiPowerBtn ShutdownFlags = 1 << iota

	// ShutdownGuestAgent - use guest agent
	ShutdownGuestAgent

	// ShutdownInitctl - use initctl
	ShutdownInitctl

	// ShutdownSignal - use signal
	ShutdownSignal

	// ShutdownParavirt - use paravirt guest control
	ShutdownParavirt
)

// DomainState specifies state of the domain
type DomainState uint32

const (
	// DomainStateNoState No state
	DomainStateNoState = iota
	// DomainStateRunning The domain is running
	DomainStateRunning
	// DomainStateBlocked The domain is blocked on resource
	DomainStateBlocked
	// DomainStatePaused The domain is paused by user
	DomainStatePaused
	// DomainStateShutdown The domain is being shut down
	DomainStateShutdown
	// DomainStateShutoff The domain is shut off
	DomainStateShutoff
	// DomainStateCrashed The domain is crashed
	DomainStateCrashed
	// DomainStatePMSuspended The domain is suspended by guest power management
	DomainStatePMSuspended
	// DomainStateLast This value will increase over time as new events are added to the libvirt
	// API. It reflects the last state supported by this version of the libvirt API.
	DomainStateLast
)

// SecretUsageType specifies the usage for a libvirt secret.
type SecretUsageType uint32

const (
	// SecretUsageTypeNone specifies no usage.
	SecretUsageTypeNone SecretUsageType = iota
	// SecretUsageTypeVolume specifies a volume secret.
	SecretUsageTypeVolume
	// SecretUsageTypeCeph specifies secrets for ceph devices.
	SecretUsageTypeCeph
	// SecretUsageTypeISCSI specifies secrets for ISCSI devices.
	SecretUsageTypeISCSI
)

// StoragePoolsFlags specifies storage pools to list.
type StoragePoolsFlags uint32

// These flags come in groups; if all bits from a group are 0,
// then that group is not used to filter results.
const (
	StoragePoolsFlagInactive = 1 << iota
	StoragePoolsFlagActive

	StoragePoolsFlagPersistent
	StoragePoolsFlagTransient

	StoragePoolsFlagAutostart
	StoragePoolsFlagNoAutostart

	// pools by type
	StoragePoolsFlagDir
	StoragePoolsFlagFS
	StoragePoolsFlagNETFS
	StoragePoolsFlagLogical
	StoragePoolsFlagDisk
	StoragePoolsFlagISCSI
	StoragePoolsFlagSCSI
	StoragePoolsFlagMPATH
	StoragePoolsFlagRBD
	StoragePoolsFlagSheepdog
	StoragePoolsFlagGluster
	StoragePoolsFlagZFS
)

// DomainCreateFlags specify options for starting domains
type DomainCreateFlags uint32

const (
	// DomainCreateFlagPaused creates paused domain.
	DomainCreateFlagPaused = 1 << iota

	// DomainCreateFlagAutoDestroy destoy domain after libvirt connection closed.
	DomainCreateFlagAutoDestroy

	// DomainCreateFlagBypassCache avoid file system cache pollution.
	DomainCreateFlagBypassCache

	// DomainCreateFlagStartForceBoot boot, discarding any managed save
	DomainCreateFlagStartForceBoot

	// DomainCreateFlagStartValidate validate the XML document against schema
	DomainCreateFlagStartValidate
)

// RebootFlags specifies domain reboot methods
type RebootFlags uint32

const (
	// RebootAcpiPowerBtn - send ACPI event
	RebootAcpiPowerBtn RebootFlags = 1 << iota

	// RebootGuestAgent - use guest agent
	RebootGuestAgent

	// RebootInitctl - use initctl
	RebootInitctl

	// RebootSignal - use signal
	RebootSignal

	// RebootParavirt - use paravirt guest control
	RebootParavirt
)

// DomainMemoryStatTag specifies domain memory tags
type DomainMemoryStatTag uint32

const (
	// DomainMemoryStatTagSwapIn - The total amount of data read from swap space (in kB).
	DomainMemoryStatTagSwapIn DomainMemoryStatTag = iota

	// DomainMemoryStatTagSwapOut - The total amount of memory written out to swap space (in kB).
	DomainMemoryStatTagSwapOut

	// DomainMemoryStatTagMajorFault - Page faults occur when a process makes a valid access to virtual memory
	// that is not available.  When servicing the page fault, if disk IO is
	// required, it is considered a major fault.
	// These are expressed as the number of faults that have occurred.
	DomainMemoryStatTagMajorFault

	// DomainMemoryStatTagMinorFault - If the page fault not require disk IO, it is a minor fault.
	DomainMemoryStatTagMinorFault

	// DomainMemoryStatTagUnused - The amount of memory left completely unused by the system (in kB).
	DomainMemoryStatTagUnused

	// DomainMemoryStatTagAvailable - The total amount of usable memory as seen by the domain (in kB).
	DomainMemoryStatTagAvailable

	// DomainMemoryStatTagActualBalloon - Current balloon value (in KB).
	DomainMemoryStatTagActualBalloon

	// DomainMemoryStatTagRss - Resident Set Size of the process running the domain (in KB).
	DomainMemoryStatTagRss

	// DomainMemoryStatTagUsable - How much the balloon can be inflated without pushing the guest system
	// to swap, corresponds to 'Available' in /proc/meminfo
	DomainMemoryStatTagUsable

	// DomainMemoryStatTagLastUpdate - Timestamp of the last update of statistics, in seconds.
	DomainMemoryStatTagLastUpdate

	// DomainMemoryStatTagNr - The number of statistics supported by this version of the interface.
	DomainMemoryStatTagNr
)

// DomainMemoryStat specifies memory stats of the domain
type DomainMemoryStat struct {
	Tag DomainMemoryStatTag
	Val uint64
}

// Capabilities returns an XML document describing the host's capabilties.
func (l *Libvirt) Capabilities() ([]byte, error) {
	resp, err := l.request(constants.ProcConnectGetCapabilties, constants.ProgramRemote, nil)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	caps, _, err := dec.DecodeString()

	return []byte(caps), err
}

// Connect establishes communication with the libvirt server.
// The underlying libvirt socket connection must be previously established.
func (l *Libvirt) Connect() error {
	return l.connect()
}

// Disconnect shuts down communication with the libvirt server
// and closes the underlying net.Conn.
func (l *Libvirt) Disconnect() error {
	// close event streams
	for id := range l.events {
		if err := l.removeStream(id); err != nil {
			return err
		}
	}

	// inform libvirt we're done
	if err := l.disconnect(); err != nil {
		return err
	}

	return l.conn.Close()
}

// Domains returns a list of all domains managed by libvirt.
func (l *Libvirt) Domains() ([]Domain, error) {
	// these are the flags as passed by `virsh`, defined in:
	// src/remote/remote_protocol.x # remote_connect_list_all_domains_args
	req := struct {
		NeedResults uint32
		Flags       uint32
	}{
		NeedResults: 1,
		Flags:       3,
	}

	buf, err := encode(&req)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcConnectListAllDomains, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	result := struct {
		Domains []Domain
		Count   uint32
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Domains, nil
}

// DomainCreateWithFlags starts specified domain with flags
func (l *Libvirt) DomainCreateWithFlags(dom string, flags DomainCreateFlags) error {
	d, err := l.lookup(dom)
	if err != nil {
		return err
	}
	req := struct {
		Domain Domain
		Flags  DomainCreateFlags
	}{
		Domain: *d,
		Flags:  flags,
	}

	buf, err := encode(&req)
	if err != nil {
		return err
	}
	resp, err := l.request(constants.ProcDomainCreateWithFlags, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}
	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}
	return nil
}

// DomainMemoryStats returns memory stats of the domain managed by libvirt.
func (l *Libvirt) DomainMemoryStats(dom string) ([]DomainMemoryStat, error) {

	d, err := l.lookup(dom)
	if err != nil {
		return nil, err
	}

	req := struct {
		Domain   Domain
		MaxStats uint32
		Flags    uint32
	}{
		Domain:   *d,
		MaxStats: 8,
		Flags:    0,
	}

	buf, err := encode(&req)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcDomainMemoryStats, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp

	result := struct {
		DomainMemoryStats []DomainMemoryStat
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.DomainMemoryStats, nil
}

// DomainState returns state of the domain managed by libvirt.
func (l *Libvirt) DomainState(dom string) (DomainState, error) {
	d, err := l.lookup(dom)
	if err != nil {
		return DomainStateNoState, err
	}

	req := struct {
		Domain Domain
		Flags  uint32
	}{
		Domain: *d,
		Flags:  0,
	}

	buf, err := encode(&req)
	if err != nil {
		return DomainStateNoState, err
	}

	resp, err := l.request(constants.ProcDomainGetState, constants.ProgramRemote, &buf)
	if err != nil {
		return DomainStateNoState, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return DomainStateNoState, decodeError(r.Payload)
	}

	result := struct {
		State  uint32
		Reason uint32
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return DomainStateNoState, err
	}

	return DomainState(result.State), nil
}

// Events streams domain events.
// If a problem is encountered setting up the event monitor connection
// an error will be returned. Errors encountered during streaming will
// cause the returned event channel to be closed.
func (l *Libvirt) Events(dom string) (<-chan DomainEvent, error) {
	d, err := l.lookup(dom)
	if err != nil {
		return nil, err
	}

	payload := struct {
		Padding [4]byte
		Domain  Domain
		Event   [2]byte
		Flags   [2]byte
	}{
		Padding: [4]byte{0x0, 0x0, 0x1, 0x0},
		Domain:  *d,
		Event:   [2]byte{0x0, 0x0},
		Flags:   [2]byte{0x0, 0x0},
	}

	buf, err := encode(&payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.QEMUConnectDomainMonitorEventRegister, constants.ProgramQEMU, &buf)
	if err != nil {
		return nil, err
	}

	res := <-resp
	if res.Status != StatusOK {
		err = decodeError(res.Payload)
		if err == ErrUnsupported {
			return nil, ErrEventsNotSupported
		}

		return nil, decodeError(res.Payload)
	}

	dec := xdr.NewDecoder(bytes.NewReader(res.Payload))

	cbID, _, err := dec.DecodeUint()
	if err != nil {
		return nil, err
	}

	stream := make(chan *DomainEvent)
	l.addStream(cbID, stream)
	c := make(chan DomainEvent)
	go func() {
		// process events
		for e := range stream {
			c <- *e
		}
	}()

	return c, nil
}

// Migrate synchronously migrates the domain specified by dom, e.g.,
// 'prod-lb-01', to the destination hypervisor specified by dest, e.g.,
// 'qemu+tcp://example.com/system'. The flags argument determines the
// type of migration and how it will be performed. For more information
// on available migration flags and their meaning, see MigrateFlag*.
func (l *Libvirt) Migrate(dom string, dest string, flags MigrateFlags) error {
	_, err := url.Parse(dest)
	if err != nil {
		return err
	}

	d, err := l.lookup(dom)
	if err != nil {
		return err
	}

	// Two unknowns remain here , Libvirt specifies RemoteParameters
	// and CookieIn. In testing both values are always set to 0 by virsh
	// and the source does not provide clear definitions of their purpose.
	// For now, using the same zero'd values as done by virsh will be Good Enough.
	payload := struct {
		Domain           Domain
		Padding          [4]byte
		DestinationURI   string
		RemoteParameters uint32
		CookieIn         uint32
		Flags            MigrateFlags
	}{
		Domain:           *d,
		Padding:          [4]byte{0x0, 0x0, 0x0, 0x1},
		DestinationURI:   dest,
		RemoteParameters: 0,
		CookieIn:         0,
		Flags:            flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcMigratePerformParams, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// MigrateSetMaxSpeed set the maximum migration bandwidth (in MiB/s) for a
// domain which is being migrated to another host. Specifying a negative value
// results in an essentially unlimited value being provided to the hypervisor.
func (l *Libvirt) MigrateSetMaxSpeed(dom string, speed int64) error {
	d, err := l.lookup(dom)
	if err != nil {
		return err
	}

	payload := struct {
		Padding   [4]byte
		Domain    Domain
		Bandwidth int64
		Flags     uint32
	}{
		Padding:   [4]byte{0x0, 0x0, 0x1, 0x0},
		Domain:    *d,
		Bandwidth: speed,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainMigrateSetMaxSpeed, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// Run executes the given QAPI command against a domain's QEMU instance.
// For a list of available QAPI commands, see:
//	http://git.qemu.org/?p=qemu.git;a=blob;f=qapi-schema.json;hb=HEAD
func (l *Libvirt) Run(dom string, cmd []byte) ([]byte, error) {
	d, err := l.lookup(dom)
	if err != nil {
		return nil, err
	}

	payload := struct {
		Domain  Domain
		Command []byte
		Flags   uint32
	}{
		Domain:  *d,
		Command: cmd,
		Flags:   0,
	}

	buf, err := encode(&payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.QEMUDomainMonitor, constants.ProgramQEMU, &buf)
	if err != nil {
		return nil, err
	}

	res := <-resp
	// check for libvirt errors
	if res.Status != StatusOK {
		return nil, decodeError(res.Payload)
	}

	// check for QEMU process errors
	if err = getQEMUError(res); err != nil {
		return nil, err
	}

	r := bytes.NewReader(res.Payload)
	dec := xdr.NewDecoder(r)
	data, _, err := dec.DecodeFixedOpaque(int32(r.Len()))
	if err != nil {
		return nil, err
	}

	// drop QMP control characters from start of line, and drop
	// any trailing NULL characters from the end
	return bytes.TrimRight(data[4:], "\x00"), nil
}

// Secrets returns all secrets managed by the libvirt daemon.
func (l *Libvirt) Secrets() ([]Secret, error) {
	req := struct {
		NeedResults uint32
		Flags       uint32
	}{
		NeedResults: 1,
		Flags:       0, // unused per libvirt source, callers should pass 0
	}

	buf, err := encode(&req)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcConnectListAllSecrets, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	result := struct {
		Secrets []Secret
		Count   uint32
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Secrets, nil
}

// StoragePool returns the storage pool associated with the provided name.
// An error is returned if the requested storage pool is not found.
func (l *Libvirt) StoragePool(name string) (*StoragePool, error) {
	req := struct {
		Name string
	}{
		Name: name,
	}

	buf, err := encode(&req)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcStoragePoolLookupByName, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	result := struct {
		Pool StoragePool
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Pool, nil
}

// StoragePoolRefresh refreshes the storage pool specified by name.
func (l *Libvirt) StoragePoolRefresh(name string) error {
	pool, err := l.StoragePool(name)
	if err != nil {
		return err
	}

	req := struct {
		Pool  StoragePool
		Flags uint32
	}{
		Pool:  *pool,
		Flags: 0, // unused per libvirt source, callers should pass 0
	}

	buf, err := encode(&req)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcStoragePoolRefresh, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// StoragePools returns a list of defined storage pools. Pools are filtered by
// the provided flags. See StoragePools*.
func (l *Libvirt) StoragePools(flags StoragePoolsFlags) ([]StoragePool, error) {
	req := struct {
		NeedResults uint32
		Flags       StoragePoolsFlags
	}{
		NeedResults: 1,
		Flags:       flags,
	}

	buf, err := encode(&req)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcConnectListAllStoragePools, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	result := struct {
		Pools []StoragePool
		Count uint32
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Pools, nil
}

// Undefine undefines the domain specified by dom, e.g., 'prod-lb-01'.
// The flags argument allows additional options to be specified such as
// cleaning up snapshot metadata. For more information on available
// flags, see UndefineFlag*.
func (l *Libvirt) Undefine(dom string, flags UndefineFlags) error {
	d, err := l.lookup(dom)
	if err != nil {
		return err
	}

	payload := struct {
		Domain Domain
		Flags  UndefineFlags
	}{
		Domain: *d,
		Flags:  flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainUndefineFlags, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// Destroy destroys the domain specified by dom, e.g., 'prod-lb-01'.
// The flags argument allows additional options to be specified such as
// allowing a graceful shutdown with SIGTERM than SIGKILL.
// For more information on available flags, see DestroyFlag*.
func (l *Libvirt) Destroy(dom string, flags DestroyFlags) error {
	d, err := l.lookup(dom)
	if err != nil {
		return err
	}

	payload := struct {
		Domain Domain
		Flags  DestroyFlags
	}{
		Domain: *d,
		Flags:  flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainDestroyFlags, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// XML returns a domain's raw XML definition, akin to `virsh dumpxml <domain>`.
// See DomainXMLFlag* for optional flags.
func (l *Libvirt) XML(dom string, flags DomainXMLFlags) ([]byte, error) {
	d, err := l.lookup(dom)
	if err != nil {
		return nil, err
	}

	payload := struct {
		Domain Domain
		Flags  DomainXMLFlags
	}{
		Domain: *d,
		Flags:  flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcDomainGetXMLDesc, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	pl := bytes.NewReader(r.Payload)
	dec := xdr.NewDecoder(pl)
	s, _, err := dec.DecodeString()
	if err != nil {
		return nil, err
	}

	return []byte(s), nil
}

// DefineXML defines a domain, but does not start it.
func (l *Libvirt) DefineXML(x []byte, flags DomainDefineXMLFlags) error {
	payload := struct {
		Domain []byte
		Flags  DomainDefineXMLFlags
	}{
		Domain: x,
		Flags:  flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainDefineXMLFlags, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// Version returns the version of the libvirt daemon.
func (l *Libvirt) Version() (string, error) {
	resp, err := l.request(constants.ProcConnectGetLibVersion, constants.ProgramRemote, nil)
	if err != nil {
		return "", err
	}

	r := <-resp
	if r.Status != StatusOK {
		return "", decodeError(r.Payload)
	}

	result := struct {
		Version uint64
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return "", err
	}

	// The version is provided as an int following this formula:
	// version * 1,000,000 + minor * 1000 + micro
	// See src/libvirt-host.c # virConnectGetLibVersion
	major := result.Version / 1000000
	result.Version %= 1000000
	minor := result.Version / 1000
	result.Version %= 1000
	micro := result.Version

	versionString := fmt.Sprintf("%d.%d.%d", major, minor, micro)
	return versionString, nil
}

// Shutdown shuts down a domain. Note that the guest OS may ignore the request.
// If flags is set to 0 then the hypervisor will choose the method of shutdown it considers best.
func (l *Libvirt) Shutdown(dom string, flags ShutdownFlags) error {
	d, err := l.lookup(dom)
	if err != nil {
		return err
	}

	payload := struct {
		Domain Domain
		Flags  ShutdownFlags
	}{
		Domain: *d,
		Flags:  flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainShutdownFlags, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// Reboot reboots the domain. Note that the guest OS may ignore the request.
// If flags is set to zero, then the hypervisor will choose the method of shutdown it considers best.
func (l *Libvirt) Reboot(dom string, flags RebootFlags) error {
	d, err := l.lookup(dom)
	if err != nil {
		return err
	}

	payload := struct {
		Domain Domain
		Flags  RebootFlags
	}{
		Domain: *d,
		Flags:  flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainReboot, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// Reset resets domain immediately without any guest OS shutdown
func (l *Libvirt) Reset(dom string) error {
	d, err := l.lookup(dom)
	if err != nil {
		return err
	}

	payload := struct {
		Domain Domain
		Flags  uint32
	}{
		Domain: *d,
		Flags:  0,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainReset, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// BlockLimit contains a name and value pair for a Get/SetBlockIOTune limit. The
// Name field is the name of the limit (to see a list of the limits that can be
// applied, execute the 'blkdeviotune' command on a VM in virsh). Callers can
// use the QEMUBlockIO... constants below for the Name value. The Value field is
// the limit to apply.
type BlockLimit struct {
	Name  string
	Value uint64
}

// BlockIOTune-able values. These tunables are different for different
// hypervisors; currently only the tunables for QEMU are defined here. These are
// not necessarily the only possible values; different libvirt versions may add
// or remove parameters from this list.
const (
	QEMUBlockIOGroupName              = "group_name"
	QEMUBlockIOTotalBytesSec          = "total_bytes_sec"
	QEMUBlockIOReadBytesSec           = "read_bytes_sec"
	QEMUBlockIOWriteBytesSec          = "write_bytes_sec"
	QEMUBlockIOTotalIOPSSec           = "total_iops_sec"
	QEMUBlockIOReadIOPSSec            = "read_iops_sec"
	QEMUBlockIOWriteIOPSSec           = "write_iops_sec"
	QEMUBlockIOTotalBytesSecMax       = "total_bytes_sec_max"
	QEMUBlockIOReadBytesSecMax        = "read_bytes_sec_max"
	QEMUBlockIOWriteBytesSecMax       = "write_bytes_sec_max"
	QEMUBlockIOTotalIOPSSecMax        = "total_iops_sec_max"
	QEMUBlockIOReadIOPSSecMax         = "read_iops_sec_max"
	QEMUBlockIOWriteIOPSSecMax        = "write_iops_sec_max"
	QEMUBlockIOSizeIOPSSec            = "size_iops_sec"
	QEMUBlockIOTotalBytesSecMaxLength = "total_bytes_sec_max_length"
	QEMUBlockIOReadBytesSecMaxLength  = "read_bytes_sec_max_length"
	QEMUBlockIOWriteBytesSecMaxLength = "write_bytes_sec_max_length"
	QEMUBlockIOTotalIOPSSecMaxLength  = "total_iops_sec_max_length"
	QEMUBlockIOReadIOPSSecMaxLength   = "read_iops_sec_max_length"
	QEMUBlockIOWriteIOPSSecMaxLength  = "write_iops_sec_max_length"
)

// SetBlockIOTune changes the per-device block I/O tunables within a guest.
// Parameters are the name of the VM, the name of the disk device to which the
// limits should be applied, and 1 or more BlockLimit structs containing the
// actual limits.
//
// The limits which can be applied here are enumerated in the QEMUBlockIO...
// constants above, and you can also see the full list by executing the
// 'blkdeviotune' command on a VM in virsh.
//
// Example usage:
//  SetBlockIOTune("vm-name", "vda", BlockLimit{libvirt.QEMUBlockIOWriteBytesSec, 1000000})
func (l *Libvirt) SetBlockIOTune(dom string, disk string, limits ...BlockLimit) error {
	d, err := l.lookup(dom)
	if err != nil {
		return err
	}

	// https://libvirt.org/html/libvirt-libvirt-domain.html#virDomainSetBlockIoTune
	payload := struct {
		Domain Domain
		Disk   string
		Params []TypedParam
		Flags  DomainAffectFlags
	}{
		Domain: *d,
		Disk:   disk,
		Flags:  FlagDomainAffectLive,
	}

	for _, limit := range limits {
		tp := NewTypedParamULongLong(limit.Name, limit.Value)
		payload.Params = append(payload.Params, tp)
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}
	resp, err := l.request(constants.ProcDomainSetBlockIOTune, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// GetBlockIOTune returns a slice containing the current block I/O tunables for
// a disk.
func (l *Libvirt) GetBlockIOTune(dom string, disk string) ([]BlockLimit, error) {
	d, err := l.lookup(dom)
	if err != nil {
		return nil, err
	}

	payload := struct {
		Domain     Domain
		Disk       []string
		ParamCount uint32
		Flags      DomainAffectFlags
	}{
		Domain:     *d,
		Disk:       []string{disk},
		ParamCount: 32,
		Flags:      FlagTypedParamStringOkay,
	}

	buf, err := encode(&payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcDomainGetBlockIOTune, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	var limits []BlockLimit
	rdr := bytes.NewReader(r.Payload)
	dec := xdr.NewDecoder(rdr)

	// find out how many params were returned
	paramCount, _, err := dec.DecodeInt()
	if err != nil {
		return nil, err
	}

	// now decode each of the returned TypedParams. To do this we read the field
	// name and type, then use the type information to decode the value.
	for param := int32(0); param < paramCount; param++ {
		// Get the field name
		name, _, err := dec.DecodeString()
		if err != nil {
			return nil, err
		}
		// ...and the type
		ptype, _, err := dec.DecodeInt()
		if err != nil {
			return nil, err
		}

		// Now we can read the actual value.
		switch ptype {
		case TypeParamULLong:
			var val uint64
			_, err = dec.Decode(&val)
			if err != nil {
				return nil, err
			}
			lim := BlockLimit{name, val}
			limits = append(limits, lim)
		case TypeParamString:
			var val string
			_, err = dec.Decode(&val)
			if err != nil {
				return nil, err
			}
			// This routine doesn't currently return strings. As of libvirt 3+,
			// there's one string here, `group_name`.
		}
	}

	return limits, nil
}

// lookup returns a domain as seen by libvirt.
func (l *Libvirt) lookup(name string) (*Domain, error) {
	payload := struct {
		Name string
	}{name}

	buf, err := encode(&payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcDomainLookupByName, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))

	var d Domain
	_, err = dec.Decode(&d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// getQEMUError checks the provided response for QEMU process errors.
// If an error is found, it is extracted an returned, otherwise nil.
func getQEMUError(r response) error {
	pl := bytes.NewReader(r.Payload)
	dec := xdr.NewDecoder(pl)

	s, _, err := dec.DecodeString()
	if err != nil {
		return err
	}

	var e qemuError
	if err = json.Unmarshal([]byte(s), &e); err != nil {
		return err
	}

	if e.Error.Description != "" {
		return errors.New(e.Error.Description)
	}

	return nil
}

// New configures a new Libvirt RPC connection.
func New(conn net.Conn) *Libvirt {
	l := &Libvirt{
		conn:      conn,
		s:         0,
		r:         bufio.NewReader(conn),
		w:         bufio.NewWriter(conn),
		callbacks: make(map[uint32]chan response),
		events:    make(map[uint32]chan *DomainEvent),
	}

	go l.listen()

	return l
}
