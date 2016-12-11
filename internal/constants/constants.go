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

// Package constants provides shared data for the libvirt package.
package constants

// protocol procedure numbers
const (
	ProgramVersion   = 1
	ProgramRemote    = 0x20008086
	ProgramQEMU      = 0x20008087
	ProgramKeepAlive = 0x6b656570
)

// libvirt procedure identifiers
// These are libvirt procedure numbers which correspond to each respective
// API call between remote_internal driver and libvirtd. Although stable.
// Each call is identified by a unique number which *may change at any time*.
//
// Examples:
//	REMOTE_PROC_CONNECT_OPEN = 1
//	REMOTE_PROC_DOMAIN_DEFINE_XML = 11
//	REMOTE_PROC_DOMAIN_MIGRATE_SET_MAX_SPEED = 207,
//
// See:
// https://libvirt.org/git/?p=libvirt.git;a=blob_plain;f=src/remote/remote_protocol.x;hb=HEAD
const (
	ProcConnectOpen              = 1
	ProcConnectClose             = 2
	ProcConnectGetCapabilties    = 7
	ProcDomainGetXMLDesc         = 14
	ProcDomainLookupByName       = 23
	ProcAuthList                 = 66
	ProcAuthSASLInit             = 67
	ProcAuthSASLStart            = 68
	ProcAuthSASLStep             = 69
	ProcAuthPolKit               = 70
	ProcConnectGetLibVersion     = 157
	ProcDomainMigrateSetMaxSpeed = 207
	ProcDomainGetState           = 212
	ProcDomainUndefineFlags      = 231
	ProcDomainDestroyFlags       = 234
	ProcConnectListAllDomains    = 273
	ProcMigratePerformParams     = 305
	ProcDomainDefineXMLFlags     = 350
)

// RemoteAuthType specifies the type of authentication used by libvirt.
type RemoteAuthType uint32

// Currently supported remote authentication types as defined by libvirt.
// See:
// https://libvirt.org/git/?p=libvirt.git;a=blob_plain;f=src/remote/remote_protocol.x;hb=HEAD
const (
	// RemoteAuthTypeNone means that no authentication is required.
	RemoteAuthTypeNone RemoteAuthType = 0
	// RemoteAuthTypeSASL means that the SASL authentication mechanism is used.
	RemoteAuthTypeSASL RemoteAuthType = 1
	// RemoteAuthTypePolKit means that PolKit is used for authentication.
	RemoteAuthTypePolKit RemoteAuthType = 2
)

// qemu procedure identifiers
const (
	QEMUDomainMonitor                       = 1
	QEMUConnectDomainMonitorEventRegister   = 4
	QEMUConnectDomainMonitorEventDeregister = 5
	QEMUDomainMonitorEvent                  = 6
)

const (
	// PacketLengthSize is the packet length, in bytes.
	PacketLengthSize = 4

	// HeaderSize is the packet header size, in bytes.
	HeaderSize = 24

	// UUIDSize is the length of a UUID, in bytes.
	UUIDSize = 16

	// AuthSASLDataMax is the upper limit on SASL auth negotiation packet.
	AuthSASLDataMax = 65536

	// AuthTypeListMax is the maximum number of auth types.
	AuthTypeListMax = 20
)
