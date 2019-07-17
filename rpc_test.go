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

package libvirt

import (
	"bytes"
	"sync"
	"testing"

	"github.com/digitalocean/go-libvirt/internal/constants"
	xdr "github.com/digitalocean/go-libvirt/internal/go-xdr/xdr2"
	"github.com/digitalocean/go-libvirt/libvirttest"
)

var (
	// dc229f87d4de47198cfd2e21c6105b01
	testUUID = [constants.UUIDSize]byte{
		0xdc, 0x22, 0x9f, 0x87, 0xd4, 0xde, 0x47, 0x19,
		0x8c, 0xfd, 0x2e, 0x21, 0xc6, 0x10, 0x5b, 0x01,
	}

	testHeader = []byte{
		0x20, 0x00, 0x80, 0x86, // program
		0x00, 0x00, 0x00, 0x01, // version
		0x00, 0x00, 0x00, 0x01, // procedure
		0x00, 0x00, 0x00, 0x00, // type
		0x00, 0x00, 0x00, 0x00, // serial
		0x00, 0x00, 0x00, 0x00, // status
	}

	testEventHeader = []byte{
		0x00, 0x00, 0x00, 0xb0, // length
		0x20, 0x00, 0x80, 0x87, // program
		0x00, 0x00, 0x00, 0x01, // version
		0x00, 0x00, 0x00, 0x06, // procedure
		0x00, 0x00, 0x00, 0x01, // type
		0x00, 0x00, 0x00, 0x00, // serial
		0x00, 0x00, 0x00, 0x00, // status
	}

	testEvent = []byte{
		0x00, 0x00, 0x00, 0x01, // callback id

		// domain name ("test")
		0x00, 0x00, 0x00, 0x04, 0x74, 0x65, 0x73, 0x74,

		// uuid (dc229f87d4de47198cfd2e21c6105b01)
		0xdc, 0x22, 0x9f, 0x87, 0xd4, 0xde, 0x47, 0x19,
		0x8c, 0xfd, 0x2e, 0x21, 0xc6, 0x10, 0x5b, 0x01,

		// domain id (14)
		0x00, 0x00, 0x00, 0x0e,

		// event name (BLOCK_JOB_COMPLETED)
		0x00, 0x00, 0x00, 0x13, 0x42, 0x4c, 0x4f, 0x43,
		0x4b, 0x5f, 0x4a, 0x4f, 0x42, 0x5f, 0x43, 0x4f,
		0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x00,

		// seconds (1462211891)
		0x00, 0x00, 0x00, 0x00, 0x57, 0x27, 0x95, 0x33,

		// microseconds (931791)
		0x00, 0x0e, 0x37, 0xcf,

		// event json data
		// ({"device":"drive-ide0-0-0","len":0,"offset":0,"speed":0,"type":"commit"})
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x48,
		0x7b, 0x22, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
		0x22, 0x3a, 0x22, 0x64, 0x72, 0x69, 0x76, 0x65,
		0x2d, 0x69, 0x64, 0x65, 0x30, 0x2d, 0x30, 0x2d,
		0x30, 0x22, 0x2c, 0x22, 0x6c, 0x65, 0x6e, 0x22,
		0x3a, 0x30, 0x2c, 0x22, 0x6f, 0x66, 0x66, 0x73,
		0x65, 0x74, 0x22, 0x3a, 0x30, 0x2c, 0x22, 0x73,
		0x70, 0x65, 0x65, 0x64, 0x22, 0x3a, 0x30, 0x2c,
		0x22, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x22,
		0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x22, 0x7d,
	}

	testErrorMessage = []byte{
		0x00, 0x00, 0x00, 0x37, // code (55, errOperationInvalid)
		0x00, 0x00, 0x00, 0x0a, // domain id

		// message ("Requested operation is not valid: domain is not running")
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x37,
		0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65,
		0x64, 0x20, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
		0x69, 0x6f, 0x6e, 0x20, 0x69, 0x73, 0x20, 0x6e,
		0x6f, 0x74, 0x20, 0x76, 0x61, 0x6c, 0x69, 0x64,
		0x3a, 0x20, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
		0x20, 0x69, 0x73, 0x20, 0x6e, 0x6f, 0x74, 0x20,
		0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x00,

		// error level
		0x00, 0x00, 0x00, 0x02,
	}

	testErrorNotFoundMessage = []byte{
		0x00, 0x00, 0x00, 0x2a, // code (42 errDoDmain)
		0x00, 0x00, 0x00, 0x0a, // domain id

		// message
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x38,
		0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x20, 0x6e,
		0x6f, 0x74, 0x20, 0x66, 0x6f, 0x75, 0x6e, 0x64,
		0x3a, 0x20, 0x6e, 0x6f, 0x20, 0x64, 0x6f, 0x6d,
		0x61, 0x69, 0x6e, 0x20, 0x77, 0x69, 0x74, 0x68,
		0x20, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e,
		0x67, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x20, 0x27,
		0x74, 0x65, 0x73, 0x74, 0x2d, 0x2d, 0x2d, 0x27,
		0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,

		// error level
		0x00, 0x00, 0x00, 0x01,
	}

	testDomain = &Domain{
		Name: "test-domain",
		UUID: testUUID,
		ID:   1,
	}
)

func TestExtractHeader(t *testing.T) {
	r := bytes.NewBuffer(testHeader)
	h, err := extractHeader(r)
	if err != nil {
		t.Error(err)
	}

	if h.Program != constants.Program {
		t.Errorf("expected Program %q, got %q", constants.Program, h.Program)
	}

	if h.Version != constants.ProtocolVersion {
		t.Errorf("expected version %q, got %q", constants.ProtocolVersion, h.Version)
	}

	if h.Procedure != constants.ProcConnectOpen {
		t.Errorf("expected procedure %q, got %q", constants.ProcConnectOpen, h.Procedure)
	}

	if h.Type != Call {
		t.Errorf("expected type %q, got %q", Call, h.Type)
	}

	if h.Status != StatusOK {
		t.Errorf("expected status %q, got %q", StatusOK, h.Status)
	}
}

func TestPktLen(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0xa} // uint32:10
	r := bytes.NewBuffer(data)

	expected := uint32(10)
	actual, err := pktlen(r)
	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("expected packet length %q, got %q", expected, actual)
	}
}

func TestDecodeEvent(t *testing.T) {
	var e DomainEvent
	err := eventDecoder(testEvent, &e)
	if err != nil {
		t.Error(err)
	}

	expCbID := uint32(1)
	if e.CallbackID != expCbID {
		t.Errorf("expected callback id %d, got %d", expCbID, e.CallbackID)
	}

	expName := "test"
	if e.Domain.Name != expName {
		t.Errorf("expected domain %s, got %s", expName, e.Domain.Name)
	}

	expUUID := testUUID
	if !bytes.Equal(e.Domain.UUID[:], expUUID[:]) {
		t.Errorf("expected uuid:\t%x, got\n\t\t\t%x", expUUID, e.Domain.UUID)
	}

	expID := int32(14)
	if e.Domain.ID != expID {
		t.Errorf("expected id %d, got %d", expID, e.Domain.ID)
	}

	expEvent := "BLOCK_JOB_COMPLETED"
	if e.Event != expEvent {
		t.Errorf("expected %s, got %s", expEvent, e.Event)
	}

	expSec := uint64(1462211891)
	if e.Seconds != expSec {
		t.Errorf("expected seconds to be %d, got %d", expSec, e.Seconds)
	}

	expMs := uint32(931791)
	if e.Microseconds != expMs {
		t.Errorf("expected microseconds to be %d, got %d", expMs, e.Microseconds)
	}

	expDetails := []byte(`{"device":"drive-ide0-0-0","len":0,"offset":0,"speed":0,"type":"commit"}`)
	if e.Domain.ID != expID {
		t.Errorf("expected data %s, got %s", expDetails, e.Details)
	}
}

func TestDecodeError(t *testing.T) {
	expectedMsg := "Requested operation is not valid: domain is not running"
	expectedCode := errOperationInvalid

	err := decodeError(testErrorMessage)
	e := err.(libvirtError)
	if e.Message != expectedMsg {
		t.Errorf("expected error message %s, got %s", expectedMsg, err.Error())
	}
	if e.Code != uint32(expectedCode) {
		t.Errorf("expected code %d, got %d", expectedCode, e.Code)
	}
}

func TestErrNotFound(t *testing.T) {
	err := decodeError(testErrorNotFoundMessage)
	ok := IsNotFound(err)
	if !ok {
		t.Errorf("expected true, got %t", ok)
	}
}

func TestEncode(t *testing.T) {
	data := "test"
	buf, err := encode(data)
	if err != nil {
		t.Error(err)
	}

	dec := xdr.NewDecoder(bytes.NewReader(buf))
	res, _, err := dec.DecodeString()
	if err != nil {
		t.Error(err)
	}

	if res != data {
		t.Errorf("expected %s, got %s", data, res)
	}
}

func TestRegister(t *testing.T) {
	l := &Libvirt{}
	l.callbacks = make(map[uint32]chan response)
	id := uint32(1)
	c := make(chan response)

	l.register(id, c)
	if _, ok := l.callbacks[id]; !ok {
		t.Error("expected callback to register")
	}
}

func TestDeregister(t *testing.T) {
	id := uint32(1)

	l := &Libvirt{}
	l.callbacks = map[uint32]chan response{
		id: make(chan response),
	}

	l.deregister(id)
	if _, ok := l.callbacks[id]; ok {
		t.Error("expected callback to deregister")
	}
}

func TestAddStream(t *testing.T) {
	id := uint32(1)

	l := &Libvirt{}
	l.events = make(map[uint32]eventStream)

	l.addStream(id, newEventStream(0, 0))
	if _, ok := l.events[id]; !ok {
		t.Error("expected event stream to exist")
	}
}

func TestRemoveStream(t *testing.T) {
	id := uint32(1)

	conn := libvirttest.New()
	l := New(conn)
	l.events[id] = newEventStream(constants.QEMUConnectDomainMonitorEventDeregister, constants.ProgramQEMU)

	err := l.removeStream(id, false)
	if err != nil {
		t.Error(err)
	}

	if _, ok := l.events[id]; ok {
		t.Error("expected event stream to be removed")
	}
}

func TestStream(t *testing.T) {
	id := uint32(1)
	c := make(chan event, 1)

	l := &Libvirt{}
	l.events = map[uint32]eventStream{
		id: eventStream{Events: c},
	}

	var streamEvent DomainEvent
	err := eventDecoder(testEvent, &streamEvent)
	if err != nil { // event was malformed, drop.
		t.Error(err)
	}

	l.stream(streamEvent)
	e := <-c

	if e.(DomainEvent).Event != "BLOCK_JOB_COMPLETED" {
		t.Error("expected event")
	}
}

func TestSerial(t *testing.T) {
	count := uint32(10)
	l := &Libvirt{}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			l.serial()
			wg.Done()
		}()
	}

	wg.Wait()

	expected := count + uint32(1)
	actual := l.serial()
	if expected != actual {
		t.Errorf("expected serial to be %d, got %d", expected, actual)
	}
}

func TestLookup(t *testing.T) {
	id := uint32(1)
	c := make(chan response)
	name := "test"

	conn := libvirttest.New()
	l := New(conn)

	l.register(id, c)

	d, err := l.lookup(name)
	if err != nil {
		t.Error(err)
	}

	if d.Name != name {
		t.Errorf("expected domain %s, got %s", name, d.Name)
	}

	// The callback should now be deregistered.
	if _, ok := l.callbacks[id]; ok {
		t.Error("expected callback to deregister")
	}
}

func TestDeregisterAll(t *testing.T) {
	conn := libvirttest.New()
	c1 := make(chan response)
	c2 := make(chan response)
	l := New(conn)
	if len(l.callbacks) != 0 {
		t.Error("expected callback map to be empty at test start")
	}
	l.register(1, c1)
	l.register(2, c2)
	if len(l.callbacks) != 2 {
		t.Error("expected callback map to have 2 entries after inserts")
	}
	l.deregisterAll()
	if len(l.callbacks) != 0 {
		t.Error("expected callback map to be empty after deregisterAll")
	}
}
