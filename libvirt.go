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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
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

// Network represents a network as seen by libvirt.
type Network struct {
	Name string
	UUID [constants.UUIDSize]byte
}

// StorageVolume represents a volume as seen by libvirt.
type StorageVolume struct {
	Pool string
	Name string
	Key  string
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

// StorageVolumeCreateFlags specifies options when performing a volume creation.
type StorageVolumeCreateFlags uint32

const (
	_ StorageVolumeCreateFlags = iota
	// StorageVolumeCreateFlagPreallocMetadata preallocates metadata
	StorageVolumeCreateFlagPreallocMetadata

	// StorageVolumeCreateFlagReflink use btrfs light copy
	StorageVolumeCreateFlagReflink
)

// DomainCreateFlags specifies options when performing a domain creation.
type DomainCreateFlags uint32

const (
	// DomainCreateFlagNone is the default behavior.
	DomainCreateFlagNone DomainCreateFlags = 0

	// DomainCreateFlagPaused creates paused domain.
	DomainCreateFlagPaused DomainCreateFlags = 1 << (iota - 1)

	// DomainCreateFlagAutoDestroy destoy domain after libvirt connection closed.
	DomainCreateFlagAutoDestroy

	// DomainCreateFlagBypassCache avoid file system cache pollution.
	DomainCreateFlagBypassCache

	// DomainCreateFlagStartForceBoot boot, discarding any managed save
	DomainCreateFlagStartForceBoot

	// DomainCreateFlagStartValidate validate the XML document against schema
	DomainCreateFlagStartValidate
)

// DomainRebootFlagValues specifies options when performing a reboot.
type DomainRebootFlags uint32

const (
	// DomainRebootFlagDefault use hypervisor choice.
	DomainRebootFlagDefault DomainRebootFlags = 0

	// DomainRebootFlagACPI send ACPI event.
	DomainRebootFlagACPI DomainRebootFlags = 1 << (iota - 1)

	// DomainRebootFlagGuestAgent use guest agent.
	DomainRebootFlagGuestAgent

	// DomainRebootFlagInitctl use initctl.
	DomainRebootFlagInitctl

	// DomainRebootFlagSignal send a signal.
	DomainRebootFlagSignal

	// DomainRebootFlagParavirt use paravirt guest control.
	DomainRebootFlagParavirt
)

// DomainShutdownFlagValues specifies options when performing a shutdown.
type DomainShutdownFlags uint32

const (
	// DomainShutdownFlagDefault use hypervisor choice.
	DomainShutdownFlagDefault DomainShutdownFlags = 0

	// DomainShutdownFlagACPI send ACPI event.
	DomainShutdownFlagACPI DomainShutdownFlags = 1 << (iota - 1)

	// DomainShutdownFlagGuestAgent use guest agent.
	DomainShutdownFlagGuestAgent

	// DomainShutdownFlagInitctl use initctl.
	DomainShutdownFlagInitctl

	// DomainShutdownFlagSignal send a signal.
	DomainShutdownFlagSignal

	// DomainShutdownFlagParavirt use paravirt guest control.
	DomainShutdownFlagParavirt
)

// MigrateFlags specifies options when performing a migration.
type DomainMigrateFlags uint32

const (
	// DomainMigrateFlagLive performs a zero-downtime live migration.
	DomainMigrateFlagLive DomainMigrateFlags = 1 << iota

	// DomainMigrateFlagPeerToPeer creates a direct source to destination control channel.
	DomainMigrateFlagPeerToPeer

	// DomainMigrateFlagTunneled tunnels migration data over the libvirtd connection.
	DomainMigrateFlagTunneled

	// DomainMigrateFlagPersistDestination will persist the VM on the destination host.
	DomainMigrateFlagPersistDestination

	// DomainMigrateFlagUndefineSource undefines the VM on the source host.
	DomainMigrateFlagUndefineSource

	// DomainMigrateFlagPaused will pause the remote side VM.
	DomainMigrateFlagPaused

	// DomainMigrateFlagNonSharedDisk migrate non-shared storage with full disk copy.
	DomainMigrateFlagNonSharedDisk

	// DomainMigrateFlagNonSharedIncremental migrate non-shared storage with incremental copy.
	DomainMigrateFlagNonSharedIncremental

	// DomainMigrateFlagChangeProtection prevents any changes to the domain configuration through the whole migration process.
	DomainMigrateFlagChangeProtection

	// DomainMigrateFlagUnsafe will force a migration even when it is considered unsafe.
	DomainMigrateFlagUnsafe

	// DomainMigrateFlagOffline is used to perform an offline migration.
	DomainMigrateFlagOffline

	// DomainMigrateFlagCompressed compresses data during migration.
	DomainMigrateFlagCompressed

	// DomainMigrateFlagAbortOnError will abort a migration on I/O errors encountered during migration.
	DomainMigrateFlagAbortOnError

	// DomainMigrateFlagAutoConverge forces convergence.
	DomainMigrateFlagAutoConverge

	// DomainMigrateFlagRDMAPinAll enables RDMA memory pinning.
	DomainMigrateFlagRDMAPinAll
)

// UndefineFlags specifies options available when undefining a domain.
type DomainUndefineFlags uint32

const (
	// DomainUndefineFlagManagedSave removes all domain managed save data.
	DomainUndefineFlagManagedSave DomainUndefineFlags = 1 << iota

	// DomainUndefineFlagSnapshotsMetadata removes all domain snapshot metadata.
	DomainUndefineFlagSnapshotsMetadata

	// DomainUndefineFlagNVRAM removes all domain NVRAM files.
	DomainUndefineFlagNVRAM
)

// DomainDefineXMLFlags specifies options available when defining a domain.
type DomainDefineXMLFlags uint32

const (
	// DefineValidate validates the XML document against schema
	DefineValidate DomainDefineXMLFlags = 1
)

// DomainDestroyFlags specifies options available when destroying a domain.
type DomainDestroyFlags uint32

const (
	// DestroyFlagDefault default behavior, forcefully terminate the domain.
	DestroyFlagDefault DomainDestroyFlags = 1 << iota

	// DestroyFlagGraceful only sends a SIGTERM no SIGKILL.
	DestroyFlagGraceful
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

// LookupDomainByName return Domain by its name.
func (l *Libvirt) LookupDomainByName(name string) (*Domain, error) {
	return l.lookupByName(name)
}

// LookupDomainByUUID return Domain by its uuid.
func (l *Libvirt) LookupDomainByUUID(uuid string) (*Domain, error) {
	return l.lookupByUUID(uuid)
}

// DomainState returns state of the domain managed by libvirt.
func (l *Libvirt) DomainState(d *Domain) (DomainState, error) {
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

// DomainEvents streams domain events.
// If a problem is encountered setting up the event monitor connection
// an error will be returned. Errors encountered during streaming will
// cause the returned event channel to be closed.
func (l *Libvirt) DomainEvents(d *Domain) (<-chan DomainEvent, error) {
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
		err := decodeError(res.Payload)
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

// DomainMigrate synchronously migrates the domain specified by dom, e.g.,
// 'prod-lb-01', to the destination hypervisor specified by dest, e.g.,
// 'qemu+tcp://example.com/system'. The flags argument determines the
// type of migration and how it will be performed. For more information
// on available migration flags and their meaning, see MigrateFlag*.
func (l *Libvirt) DomainMigrate(d *Domain, dest string, flags DomainMigrateFlags) error {
	_, err := url.Parse(dest)
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
		Flags            DomainMigrateFlags
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

// DomainMigrateSetMaxSpeed set the maximum migration bandwidth (in MiB/s) for a
// domain which is being migrated to another host. Specifying a negative value
// results in an essentially unlimited value being provided to the hypervisor.
func (l *Libvirt) DomainMigrateSetMaxSpeed(d *Domain, speed int64) error {
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
	d, err := l.lookupByName(dom)
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

// StoragePoolLookupByName returns the storage pool associated with the provided name.
// An error is returned if the requested storage pool is not found.
func (l *Libvirt) StoragePoolLookupByName(name string) (*StoragePool, error) {
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

// NetworkLookupByName returns the network associated with the provided name.
// An error is returned if the requested network is not found.
func (l *Libvirt) NetworkLookupByName(name string) (*Network, error) {
	req := struct {
		Name string
	}{
		Name: name,
	}

	buf, err := encode(&req)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcNetworkLookupByName, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	result := struct {
		Net Network
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Net, nil
}

// StoragePoolLookupByUUID returns the storage pool associated with the provided uuid.
// An error is returned if the requested storage pool is not found.
func (l *Libvirt) StoragePoolLookupByUUID(uuid string) (*StoragePool, error) {
	req := struct {
		UUID [constants.UUIDSize]byte
	}{}

	_, err := hex.Decode(req.UUID[:], []byte(strings.Replace(uuid, "-", "", -1)))
	if err != nil {
		return nil, err
	}

	buf, err := encode(&req)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcStoragePoolLookupByUUID, constants.ProgramRemote, &buf)
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
func (l *Libvirt) StoragePoolRefresh(pool *StoragePool, flags uint32) error {
	req := struct {
		Pool  StoragePool
		Flags uint32
	}{
		Pool:  *pool,
		Flags: flags, // unused per libvirt source, callers should pass 0
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

// DomainUndefine undefines the domain specified by dom, e.g., 'prod-lb-01'.
// The flags argument allows additional options to be specified such as
// cleaning up snapshot metadata. For more information on available
// flags, see DomainUndefineFlag*.
func (l *Libvirt) DomainUndefine(d *Domain, flags DomainUndefineFlags) error {
	payload := struct {
		Domain Domain
		Flags  DomainUndefineFlags
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

// DomainSuspend suspends the domain.
func (l *Libvirt) DomainSuspend(d *Domain) error {
	payload := struct {
		Domain Domain
	}{
		Domain: *d,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainSuspend, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// DomainResume resume domain.
func (l *Libvirt) DomainResume(d *Domain) error {
	payload := struct {
		Domain Domain
	}{
		Domain: *d,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainResume, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// DomainSetAutostart set autostart for domain.
func (l *Libvirt) DomainSetAutostart(d *Domain, autostart bool) error {
	payload := struct {
		Domain    Domain
		Autostart int32
	}{}

	payload.Domain = *d
	if autostart {
		payload.Autostart = 1
	} else {
		payload.Autostart = 0
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainSetAutostart, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// StoragePoolSetAutostart set autostart for domain.
func (l *Libvirt) StoragePoolSetAutostart(p *StoragePool, autostart bool) error {
	payload := struct {
		Pool      StoragePool
		Autostart int32
	}{}

	payload.Pool = *p
	if autostart {
		payload.Autostart = 1
	} else {
		payload.Autostart = 0
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcStoragePoolSetAutostart, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// NetworkSetAutostart set autostart for domain.
func (l *Libvirt) NetworkSetAutostart(n *Network, autostart bool) error {
	payload := struct {
		Network   Network
		Autostart int32
	}{}

	payload.Network = *n
	if autostart {
		payload.Autostart = 1
	} else {
		payload.Autostart = 0
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcNetworkSetAutostart, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// DomainDestroy destroys the domain.
// The flags argument allows additional options to be specified such as
// allowing a graceful shutdown with SIGTERM than SIGKILL.
// For more information on available flags, see DomainDestroyFlag*.
func (l *Libvirt) DomainDestroy(d *Domain, flags DomainDestroyFlags) error {
	payload := struct {
		Domain Domain
		Flags  DomainDestroyFlags
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

// DomainReboot reboot the domain.
// The flags argument allows additional options to be specified.
// For more information on available flags, see DomainRebootFlags*.
func (l *Libvirt) DomainReboot(d *Domain, flags DomainRebootFlags) error {
	payload := struct {
		Domain Domain
		Flags  DomainRebootFlags
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

// DomainShutdown reboot the domain.
// The flags argument allows additional options to be specified.
// For more information on available flags, see DomainShutdownFlags*.
func (l *Libvirt) DomainShutdown(d *Domain, flags DomainShutdownFlags) error {
	payload := struct {
		Domain Domain
		Flags  DomainShutdownFlags
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

// DomainXML returns a domain's raw XML definition, akin to `virsh dumpxml <domain>`.
// See DomainXMLFlag* for optional flags.
func (l *Libvirt) DomainXML(d *Domain, flags DomainXMLFlags) ([]byte, error) {
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
		XML   []byte
		Flags DomainDefineXMLFlags
	}{
		XML:   x,
		Flags: flags,
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

// StorageVolumeCreateXML defines a volume, but does not start it.
func (l *Libvirt) StorageVolumeCreateXML(p *StoragePool, x []byte, flags StorageVolumeCreateFlags) (*StorageVolume, error) {
	payload := struct {
		Pool  StoragePool
		XML   []byte
		Flags StorageVolumeCreateFlags
	}{
		Pool:  *p,
		XML:   x,
		Flags: flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcStorageVolumeCreateXML, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	result := struct {
		Volume StorageVolume
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Volume, nil

}

// DomainCreate start defined domain.
func (l *Libvirt) DomainCreate(d *Domain) error {
	payload := struct {
		Domain Domain
	}{
		Domain: *d,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainCreate, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// DomainCreateXML start domain based on xml.
func (l *Libvirt) DomainCreateXML(x []byte, flags DomainCreateFlags) error {
	payload := struct {
		XML   []byte
		Flags DomainCreateFlags
	}{
		XML:   x,
		Flags: flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcDomainCreateXML, constants.ProgramRemote, &buf)
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

// lookupByName returns a domain as seen by libvirt.
func (l *Libvirt) lookupByName(name string) (*Domain, error) {
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

// StorageVolumeLookupByName returns a volume as seen by libvirt.
func (l *Libvirt) StorageVolumeLookupByName(p *StoragePool, name string) (*StorageVolume, error) {
	payload := struct {
		Pool StoragePool
		Name string
	}{
		Pool: *p,
		Name: name,
	}

	buf, err := encode(&payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcStorageVolumeLookupByName, constants.ProgramRemote, &buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	if r.Status != StatusOK {
		return nil, decodeError(r.Payload)
	}

	result := struct {
		Volume StorageVolume
	}{}

	dec := xdr.NewDecoder(bytes.NewReader(r.Payload))
	_, err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Volume, nil
}

// lookupByUUID returns a domain as seen by libvirt.
func (l *Libvirt) lookupByUUID(uuid string) (*Domain, error) {
	payload := struct {
		UUID [constants.UUIDSize]byte
	}{}
	_, err := hex.Decode(payload.UUID[:], []byte(strings.Replace(uuid, "-", "", -1)))
	if err != nil {
		return nil, err
	}

	buf, err := encode(&payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.request(constants.ProcDomainLookupByUUID, constants.ProgramRemote, &buf)
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
