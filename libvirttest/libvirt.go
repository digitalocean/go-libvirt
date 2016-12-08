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

// Package libvirttest provides a mock libvirt server for RPC testing.
package libvirttest

import (
	"encoding/binary"
	"net"
	"sync/atomic"

	"github.com/digitalocean/go-libvirt/internal/constants"
)

var testDomainResponse = []byte{
	0x00, 0x00, 0x00, 0x38, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x17, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status

	// domain name ("test")
	0x00, 0x00, 0x00, 0x04, 0x74, 0x65, 0x73, 0x74,

	// uuid (dc229f87d4de47198cfd2e21c6105b01)
	0xdc, 0x22, 0x9f, 0x87, 0xd4, 0xde, 0x47, 0x19,
	0x8c, 0xfd, 0x2e, 0x21, 0xc6, 0x10, 0x5b, 0x01,

	// domain id (14)
	0x00, 0x00, 0x00, 0x0e,
}

var testRegisterEvent = []byte{
	0x00, 0x00, 0x00, 0x20, // length
	0x20, 0x00, 0x80, 0x87, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x04, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
	0x00, 0x00, 0x00, 0x01, // callback id
}

var testDeregisterEvent = []byte{
	0x00, 0x00, 0x00, 0x1c, // length
	0x20, 0x00, 0x80, 0x87, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x05, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
}

var testAuthReply = []byte{
	0x00, 0x00, 0x00, 0x1c, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x42, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
}

var testConnectReply = []byte{
	0x00, 0x00, 0x00, 0x1c, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x01, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
}

var testDisconnectReply = []byte{
	0x00, 0x00, 0x00, 0x1c, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x02, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
}

var testMigrateReply = []byte{
	0x00, 0x00, 0x00, 0x20, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x01, 0x31, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status

	// cookie out: 0
	0x00, 0x00, 0x00, 0x00,
}

var testRunReply = []byte{
	0x00, 0x00, 0x00, 0x74, // length
	0x20, 0x00, 0x80, 0x87, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x01, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status

	// {"return":{"qemu":{"micro":1,"minor":5,"major":2},"package":""},"id":"libvirt-53"}
	0x00, 0x00, 0x00, 0x52, 0x7b, 0x22, 0x72, 0x65,
	0x74, 0x75, 0x72, 0x6e, 0x22, 0x3a, 0x7b, 0x22,
	0x71, 0x65, 0x6d, 0x75, 0x22, 0x3a, 0x7b, 0x22,
	0x6d, 0x69, 0x63, 0x72, 0x6f, 0x22, 0x3a, 0x31,
	0x2c, 0x22, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x22,
	0x3a, 0x35, 0x2c, 0x22, 0x6d, 0x61, 0x6a, 0x6f,
	0x72, 0x22, 0x3a, 0x32, 0x7d, 0x2c, 0x22, 0x70,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x22, 0x3a,
	0x22, 0x22, 0x7d, 0x2c, 0x22, 0x69, 0x64, 0x22,
	0x3a, 0x22, 0x6c, 0x69, 0x62, 0x76, 0x69, 0x72,
	0x74, 0x2d, 0x35, 0x33, 0x22, 0x7d,

	// All trailing NULL characters should be removed
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

var testRunReplyFail = []byte{
	0x00, 0x00, 0x00, 0x8c, // length
	0x20, 0x00, 0x80, 0x87, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x01, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x0a, // serial
	0x00, 0x00, 0x00, 0x00, // status

	// {"id":"libvirt-68","error":{"class":"CommandNotFound","desc":"The command drive-foo has not been found"}}`
	0x00, 0x00, 0x00, 0x69, 0x7b, 0x22, 0x69, 0x64,
	0x22, 0x3a, 0x22, 0x6c, 0x69, 0x62, 0x76, 0x69,
	0x72, 0x74, 0x2d, 0x36, 0x38, 0x22, 0x2c, 0x22,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3a, 0x7b,
	0x22, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x22, 0x3a,
	0x22, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64,
	0x22, 0x2c, 0x22, 0x64, 0x65, 0x73, 0x63, 0x22,
	0x3a, 0x22, 0x54, 0x68, 0x65, 0x20, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x20, 0x64, 0x72,
	0x69, 0x76, 0x65, 0x2d, 0x66, 0x6f, 0x6f, 0x20,
	0x68, 0x61, 0x73, 0x20, 0x6e, 0x6f, 0x74, 0x20,
	0x62, 0x65, 0x65, 0x6e, 0x20, 0x66, 0x6f, 0x75,
	0x6e, 0x64, 0x22, 0x7d, 0x7d, 0x00, 0x00, 0x00,
}

var testSetSpeedReply = []byte{
	0x00, 0x00, 0x00, 0x1c, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0xcf, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
}

var testDomainsReply = []byte{
	0x00, 0x00, 0x00, 0x6c, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x01, 0x11, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status

	// struct of domains
	0x00, 0x00, 0x00, 0x02,

	// first domain
	// name - aaaaaaa-1
	0x00, 0x00, 0x00, 0x09, 0x61, 0x61, 0x61, 0x61,
	0x61, 0x61, 0x61, 0x2d, 0x31, 0x00, 0x00, 0x00,
	// uuid - dc:32:9f:87:d4:de:47:19:8c:fd:2e:21:c6:10:5b:01
	0xdc, 0x32, 0x9f, 0x87, 0xd4, 0xde, 0x47, 0x19,
	0x8c, 0xfd, 0x2e, 0x21, 0xc6, 0x10, 0x5b, 0x01,
	// id
	0x00, 0x00, 0x00, 0x01,

	// second domain
	// name - aaaaaaa-2
	0x00, 0x00, 0x00, 0x09, 0x61, 0x61, 0x61, 0x61,
	0x61, 0x61, 0x61, 0x2d, 0x32, 0x00, 0x00, 0x00,
	// uuid - dc:22:9f:87:d4:de:47:19:8c:fd:2e:21:c6:10:5b:01
	0xdc, 0x22, 0x9f, 0x87, 0xd4, 0xde, 0x47, 0x19, 0x8c,
	0xfd, 0x2e, 0x21, 0xc6, 0x10, 0x5b, 0x01, 0x00, 0x00,
	// id
	0x00, 0x02, 0x00,

	// count of domains returned
	0x00, 0x00, 0x02,
}

var testDomainStateReply = []byte{
	0x00, 0x00, 0x00, 0x24, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0xd4, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
	0x00, 0x00, 0x00, 0x01, // state
	0x00, 0x00, 0x00, 0x01, // reason
}

var testUndefineReply = []byte{
	0x00, 0x00, 0x00, 0x1c, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0xe7, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
}

var testDestroyReply = []byte{
	0x00, 0x00, 0x00, 0x1c, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0xea, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
}

var testVersionReply = []byte{
	0x00, 0x00, 0x00, 0x24, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x00, 0x9d, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
	0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x4d, 0xfc, // version (1003004)
}

var testDefineXML = []byte{
	0x00, 0x00, 0x00, 0x38, // length
	0x20, 0x00, 0x80, 0x86, // program
	0x00, 0x00, 0x00, 0x01, // version
	0x00, 0x00, 0x01, 0x5e, // procedure
	0x00, 0x00, 0x00, 0x01, // type
	0x00, 0x00, 0x00, 0x00, // serial
	0x00, 0x00, 0x00, 0x00, // status
	0x00, 0x00, 0x00, 0x04, // dom
	0x74, 0x65, 0x73, 0x74, // name
	// uuid
	0xaf, 0xc2, 0xef, 0x71, 0x66, 0xe0, 0x45, 0xa7,
	0xa5, 0xec, 0xd8, 0xba, 0x1e, 0xa8, 0x17, 0x7d,
	0xff, 0xff, 0xff, 0xff, // id
}

// MockLibvirt provides a mock libvirt server for testing.
type MockLibvirt struct {
	net.Conn
	Test   net.Conn
	Fail   bool
	serial uint32
}

// New creates a new mock Libvirt server.
func New() *MockLibvirt {
	serv, conn := net.Pipe()

	m := &MockLibvirt{
		Conn: conn,
		Test: serv,
	}

	go m.handle(serv)

	return m
}

func (m *MockLibvirt) handle(conn net.Conn) {
	for {
		// packetLengthSize + headerSize
		buf := make([]byte, 28)
		conn.Read(buf)

		// extract program
		prog := binary.BigEndian.Uint32(buf[4:8])

		// extract procedure
		proc := binary.BigEndian.Uint32(buf[12:16])

		switch prog {
		case constants.ProgramRemote:
			m.handleRemote(proc, conn)
		case constants.ProgramQEMU:
			m.handleQEMU(proc, conn)
		}
	}
}

func (m *MockLibvirt) handleRemote(procedure uint32, conn net.Conn) {
	switch procedure {
	case constants.ProcAuthList:
		conn.Write(m.reply(testAuthReply))
	case constants.ProcConnectOpen:
		conn.Write(m.reply(testConnectReply))
	case constants.ProcConnectClose:
		conn.Write(m.reply(testDisconnectReply))
	case constants.ProcConnectGetLibVersion:
		conn.Write(m.reply(testVersionReply))
	case constants.ProcDomainLookupByName:
		conn.Write(m.reply(testDomainResponse))
	case constants.ProcConnectListAllDomains:
		conn.Write(m.reply(testDomainsReply))
	case constants.ProcDomainGetState:
		conn.Write(m.reply(testDomainStateReply))
	case constants.ProcDomainMigrateSetMaxSpeed:
		conn.Write(m.reply(testSetSpeedReply))
	case constants.ProcMigratePerformParams:
		conn.Write(m.reply(testMigrateReply))
	case constants.ProcDomainUndefineFlags:
		conn.Write(m.reply(testUndefineReply))
	case constants.ProcDomainDestroyFlags:
		conn.Write(m.reply(testDestroyReply))
	case constants.ProcDomainDefineXMLFlags:
		conn.Write(m.reply(testDefineXML))
	}
}

func (m *MockLibvirt) handleQEMU(procedure uint32, conn net.Conn) {
	switch procedure {
	case constants.QEMUConnectDomainMonitorEventRegister:
		conn.Write(m.reply(testRegisterEvent))
	case constants.QEMUConnectDomainMonitorEventDeregister:
		conn.Write(m.reply(testDeregisterEvent))
	case constants.QEMUDomainMonitor:
		if m.Fail {
			conn.Write(m.reply(testRunReplyFail))
		} else {
			conn.Write(m.reply(testRunReply))
		}
	}
}

// reply automatically injects the correct serial
// number into the provided response buffer.
func (m *MockLibvirt) reply(buf []byte) []byte {
	atomic.AddUint32(&m.serial, 1)
	binary.BigEndian.PutUint32(buf[20:24], m.serial)

	return buf
}
