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
	ProcDomainLookupByName       = 23
	ProcAuthList                 = 66
	ProcConnectGetLibVersion     = 157
	ProcDomainMigrateSetMaxSpeed = 207
	ProcDomainGetState           = 212
	ProcDomainUndefineFlags      = 231
	ProcDomainDestroyFlags       = 234
	ProcConnectListAllDomains    = 273
	ProcMigratePerformParams     = 305
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
)

// Domain state
const (
	// No state
	DomainNoState = 0
	// The domain is running
	DomainRunning = 1
	// The domain is blocked on resource
	DomainBlocked = 2
	// The domain is paused by user
	DomainPaused = 3
	// The domain is being shut down
	DomainShutdown = 4
	// The domain is shut off
	DomainShutoff = 5
	// The domain is crashed
	DomainCrashed = 6
	// The domain is suspended by guest power management
	DomainPMSuspended = 7
	// This value will increase over time as new events are added to the libvirt
	// API. It reflects the last state supported by this version of the libvirt API.
	DomainLast = 8
)
