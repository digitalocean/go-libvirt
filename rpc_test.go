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
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/digitalocean/go-libvirt/internal/constants"
	"github.com/digitalocean/go-libvirt/internal/event"
	xdr "github.com/digitalocean/go-libvirt/internal/go-xdr/xdr2"
	"github.com/digitalocean/go-libvirt/libvirttest"
	"github.com/digitalocean/go-libvirt/socket"
	"github.com/stretchr/testify/assert"
)

var (
	// dc229f87d4de47198cfd2e21c6105b01
	testUUID = [UUIDBuflen]byte{
		0xdc, 0x22, 0x9f, 0x87, 0xd4, 0xde, 0x47, 0x19,
		0x8c, 0xfd, 0x2e, 0x21, 0xc6, 0x10, 0x5b, 0x01,
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

	testLifeCycle = []byte{
		0x00, 0x00, 0x00, 0x01, // callback id

		// domain name ("test")
		0x00, 0x00, 0x00, 0x04, 0x74, 0x65, 0x73, 0x74,

		// event data
		0x00, 0x00, 0x00, 0x50, 0xad, 0xf7, 0x3f, 0xbe, 0xca, 0x48, 0xac, 0x95,
		0x13, 0x8a, 0x31, 0xf4, 0xfe, 0x03, 0x2a, 0xff, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x01,
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
)

func TestDecodeEvent(t *testing.T) {
	var e DomainEvent
	err := eventDecoder(testEvent, &e)
	if err != nil {
		t.Error(err)
	}

	expCbID := int32(1)
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
	expectedCode := ErrOperationInvalid

	err := decodeError(testErrorMessage)
	e := err.(Error)
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

	err = fmt.Errorf("something went wrong: %w", err)
	ok = IsNotFound(err)
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
	l.callbacks = make(map[int32]chan response)
	id := int32(1)
	c := make(chan response)

	l.register(id, c)
	if _, ok := l.callbacks[id]; !ok {
		t.Error("expected callback to register")
	}
}

func TestDeregister(t *testing.T) {
	id := int32(1)

	l := &Libvirt{}
	l.callbacks = map[int32]chan response{
		id: make(chan response),
	}

	l.deregister(id)
	if _, ok := l.callbacks[id]; ok {
		t.Error("expected callback to deregister")
	}
}

func TestAddStream(t *testing.T) {
	id := int32(1)

	l := &Libvirt{}
	l.events = make(map[int32]*event.Stream)

	stream := event.NewStream(0, id)
	defer stream.Shutdown()

	l.addStream(stream)
	if _, ok := l.events[id]; !ok {
		t.Error("expected event stream to exist")
	}
}

func TestRemoveStream(t *testing.T) {
	id := int32(1)

	dialer := libvirttest.New()
	l := NewWithDialer(dialer)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	stream := event.NewStream(constants.QEMUProgram, id)
	defer stream.Shutdown()

	l.addStream(stream)

	fmt.Println("removing stream")
	err = l.removeStream(id)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := l.events[id]; ok {
		t.Error("expected event stream to be removed")
	}
}

func TestRemoveAllStreams(t *testing.T) {
	id1 := int32(1)
	id2 := int32(2)

	dialer := libvirttest.New()
	l := NewWithDialer(dialer)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	// verify it's a successful no-op when no streams have been added
	l.removeAllStreams()
	if len(l.events) != 0 {
		t.Fatal("expected no streams after remove all")
	}

	stream := event.NewStream(constants.QEMUProgram, id1)
	defer stream.Shutdown()

	l.addStream(stream)

	stream2 := event.NewStream(constants.QEMUProgram, id2)
	defer stream2.Shutdown()

	l.addStream(stream2)

	l.removeAllStreams()

	if len(l.events) != 0 {
		t.Error("expected all event streams to be removed")
	}
}

func TestStream(t *testing.T) {
	id := int32(1)
	stream := event.NewStream(constants.Program, 1)
	defer stream.Shutdown()

	l := &Libvirt{}
	l.events = map[int32]*event.Stream{
		id: stream,
	}

	var streamEvent DomainEvent
	err := eventDecoder(testEvent, &streamEvent)
	if err != nil { // event was malformed, drop.
		t.Error(err)
	}

	l.stream(streamEvent)
	e := <-stream.Recv()

	if e.(DomainEvent).Event != "BLOCK_JOB_COMPLETED" {
		t.Error("expected event")
	}
}

func TestSerial(t *testing.T) {
	count := int32(10)
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

	expected := count + int32(1)
	actual := l.serial()
	if expected != actual {
		t.Errorf("expected serial to be %d, got %d", expected, actual)
	}
}

func TestLookup(t *testing.T) {
	name := "test"

	dialer := libvirttest.New()
	l := NewWithDialer(dialer)

	err := l.Connect()
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer l.Disconnect()

	d, err := l.lookup(name)
	if err != nil {
		t.Error(err)
	}

	if d.Name != name {
		t.Errorf("expected domain %s, got %s", name, d.Name)
	}
}

func TestDeregisterAll(t *testing.T) {
	dialer := libvirttest.New()
	c1 := make(chan response)
	c2 := make(chan response)
	l := NewWithDialer(dialer)
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

// TestRouteDeadlock ensures that go-libvirt doesn't hang when trying to send
// both an event and a response (to a request) at the same time.
//
// Events are inherently asynchronous - the client may not be ready to receive
// an event when it arrives. We don't want that to prevent go-libvirt from
// continuing to receive responses to outstanding requests. This test checks for
// deadlocks where the client doesn't immediately consume incoming events.
func TestRouteDeadlock(t *testing.T) {
	id := int32(1)
	rch := make(chan response, 1)

	l := &Libvirt{
		callbacks: map[int32]chan response{
			id: rch,
		},
		events: make(map[int32]*event.Stream),
	}
	stream := event.NewStream(constants.Program, id)

	l.addStream(stream)

	respHeader := &socket.Header{
		Program: constants.Program,
		Serial:  id,
		Status:  socket.StatusOK,
	}
	eventHeader := &socket.Header{
		Program:   constants.Program,
		Procedure: constants.ProcDomainEventCallbackLifecycle,
		Status:    socket.StatusOK,
	}

	send := func(respCount, evCount int) {
		// Send the events first
		for i := 0; i < evCount; i++ {
			l.Route(eventHeader, testLifeCycle)
		}
		// Now send the requests.
		for i := 0; i < respCount; i++ {
			l.Route(respHeader, []byte{})
		}
	}

	cases := []struct{ rCount, eCount int }{
		{2, 0},
		{0, 2},
		{1, 1},
		{2, 2},
		{50, 50},
	}

	for _, tc := range cases {
		fmt.Printf("testing %d responses and %d events\n", tc.rCount, tc.eCount)
		go send(tc.rCount, tc.eCount)

		for i := 0; i < tc.rCount; i++ {
			r := <-rch
			assert.Equal(t, r.Status, uint32(socket.StatusOK))
		}
		for i := 0; i < tc.eCount; i++ {
			e := <-stream.Recv()
			fmt.Printf("event %v/%v received\n", i, len(cases))
			assert.Equal(t, "test", e.(*DomainEventCallbackLifecycleMsg).Msg.Dom.Name)
		}
	}

	// finally verify that canceling the context doesn't cause a deadlock.
	fmt.Println("checking for deadlock after context cancellation")
	send(0, 50)
}

// TestQEMUEventHandoffWindow reproduces the handoff window that used to
// exist during SubscribeQEMUEvents teardown: the stream's local reader
// (stream.Shutdown()) torn down before it was deregistered from routing
// (l.removeStream, via unsubscribeQEMUEvents). If an event for that
// callback ID arrived from libvirtd in that window, the shared
// socket-reader goroutine would call Route -> stream -> Push on a stream
// with nothing left to drain it, blocking forever. This simulates that
// exact ordering directly, since the real background goroutine's internal
// scheduling can't be raced deterministically from a test.
func TestQEMUEventHandoffWindow(t *testing.T) {
	id := int32(1)

	l := &Libvirt{events: make(map[int32]*event.Stream)}
	stream := event.NewStream(constants.QEMUProgram, id)
	l.addStream(stream)

	// Simulate the pre-fix teardown order: local reader gone, but the
	// stream is still registered for routing.
	stream.Shutdown()
	_, ok := <-stream.Recv()
	assert.False(t, ok)

	eventHeader := &socket.Header{
		Program:   constants.QEMUProgram,
		Procedure: constants.QEMUProcDomainMonitorEvent,
		Status:    socket.StatusOK,
	}

	done := make(chan struct{})
	go func() {
		l.Route(eventHeader, testEvent)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Route/Push blocked delivering an event to a torn-down stream")
	}
}

func TestGetResponseInterrupted(t *testing.T) {
	dialer := libvirttest.New()
	l := NewWithDialer(dialer)
	c := make(chan response)
	close(c)
	_, err := l.getResponse(c)
	assert.Equal(t, ErrInterrupted, err)
}

// Bidirectional console (proc 201) regression tests for the requestStream
// both-streams teardown deadlocks. See digitalocean/go-libvirt#260

// consoleTestServer is the server end of a net.Pipe used to script a fake
// libvirtd that speaks the streaming subset of the remote protocol
type consoleTestServer struct {
	conn net.Conn
}

func (s *consoleTestServer) readPacket() (socket.Header, []byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(s.conn, lenBuf[:]); err != nil {
		return socket.Header{}, nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf[:])
	buf := make([]byte, int(length)-4)
	if _, err := io.ReadFull(s.conn, buf); err != nil {
		return socket.Header{}, nil, err
	}
	h := socket.Header{
		Program:   binary.BigEndian.Uint32(buf[0:4]),
		Version:   binary.BigEndian.Uint32(buf[4:8]),
		Procedure: binary.BigEndian.Uint32(buf[8:12]),
		Type:      binary.BigEndian.Uint32(buf[12:16]),
		Serial:    int32(binary.BigEndian.Uint32(buf[16:20])),
		Status:    binary.BigEndian.Uint32(buf[20:24]),
	}
	if len(buf) <= 24 {
		return h, nil, nil
	}
	return h, buf[24:], nil
}

func (s *consoleTestServer) writePacket(serial int32, proc, typ, status uint32, payload []byte) error {
	buf := make([]byte, 28+len(payload))
	binary.BigEndian.PutUint32(buf[0:4], uint32(len(buf)))
	binary.BigEndian.PutUint32(buf[4:8], constants.Program)
	binary.BigEndian.PutUint32(buf[8:12], constants.ProtocolVersion)
	binary.BigEndian.PutUint32(buf[12:16], proc)
	binary.BigEndian.PutUint32(buf[16:20], typ)
	binary.BigEndian.PutUint32(buf[20:24], uint32(serial))
	binary.BigEndian.PutUint32(buf[24:28], status)
	copy(buf[28:], payload)
	_, err := s.conn.Write(buf)
	return err
}

// setupConsoleTest returns a Libvirt whose socket listener is running against
// a net.Pipe, plus the server end for scripting the fake libvirtd. It bypasses
// Connect() and the auth handshake to keep the tests focused on stream
// orchestration.
func setupConsoleTest(t *testing.T) (*Libvirt, *consoleTestServer) {
	t.Helper()
	client, server := net.Pipe()
	t.Cleanup(func() {
		client.Close()
		server.Close()
	})
	l := New(client)
	if err := l.socket.Connect(); err != nil {
		t.Fatalf("socket connect failed: %v", err)
	}
	return l, &consoleTestServer{conn: server}
}

// TestDomainOpenConsoleBidirectionalStdinEOF ensures that closing stdin
// cleanly terminates the whole stream
func TestDomainOpenConsoleBidirectionalStdinEOF(t *testing.T) {
	l, srv := setupConsoleTest(t)

	stdinData := []byte("root\r")
	consoleOutput := []byte("login: ")

	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)

		hdr, _, err := srv.readPacket()
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, uint32(socket.Call), hdr.Type)
		serial, proc := hdr.Serial, hdr.Procedure

		// Acknowledge the console-open request.
		if !assert.NoError(t, srv.writePacket(serial, proc, socket.Reply, socket.StatusOK, nil)) {
			return
		}

		// Drain the client->server stream until the client signals
		// end-of-stream (stdin EOF).
		var got bytes.Buffer
		for {
			hdr, payload, err := srv.readPacket()
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, uint32(socket.Stream), hdr.Type)
			if hdr.Status == socket.StatusOK {
				break
			}
			got.Write(payload)
		}
		assert.Equal(t, stdinData, got.Bytes())

		// Emit console output, then end the server->client stream.
		if !assert.NoError(t, srv.writePacket(serial, proc, socket.Stream, socket.StatusContinue, consoleOutput)) {
			return
		}
		assert.NoError(t, srv.writePacket(serial, proc, socket.Stream, socket.StatusOK, nil))
	}()

	var console bytes.Buffer
	done := make(chan error, 1)
	go func() {
		done <- l.DomainOpenConsoleBidirectional(Domain{}, nil, bytes.NewReader(stdinData), &console, 0)
	}()

	select {
	case err := <-done:
		assert.NoError(t, err)
		assert.Equal(t, consoleOutput, console.Bytes())
	case <-time.After(2 * time.Second):
		t.Fatal("DomainOpenConsoleBidirectional did not return after stdin EOF and stream end")
	}
	<-serverDone
}

// TestDomainOpenConsoleBidirectionalServerEndsFirst ensures that when libvirtd
// ends the console stream while local stdin is still open and idle the call
// returns instead of waiting for stdin forever.
func TestDomainOpenConsoleBidirectionalServerEndsFirst(t *testing.T) {
	l, srv := setupConsoleTest(t)

	consoleOutput := []byte("kernel panic\r\n")

	// Emulate interactive stdin: open and idle.
	stdinR, stdinW := io.Pipe()
	t.Cleanup(func() { stdinW.Close() })

	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)

		hdr, _, err := srv.readPacket()
		if !assert.NoError(t, err) {
			return
		}
		serial, proc := hdr.Serial, hdr.Procedure

		if !assert.NoError(t, srv.writePacket(serial, proc, socket.Reply, socket.StatusOK, nil)) {
			return
		}

		// End the stream without the client having sent anything.
		if !assert.NoError(t, srv.writePacket(serial, proc, socket.Stream, socket.StatusContinue, consoleOutput)) {
			return
		}
		assert.NoError(t, srv.writePacket(serial, proc, socket.Stream, socket.StatusOK, nil))
	}()

	var console bytes.Buffer
	done := make(chan error, 1)
	go func() {
		done <- l.DomainOpenConsoleBidirectional(Domain{}, nil, stdinR, &console, 0)
	}()

	select {
	case err := <-done:
		assert.NoError(t, err)
		assert.Equal(t, consoleOutput, console.Bytes())
	case <-time.After(2 * time.Second):
		t.Fatal("DomainOpenConsoleBidirectional hung waiting for stdin after libvirtd ended the stream")
	}
	<-serverDone
}

// TestDomainOpenConsoleBidirectionalStreamError ensures a stream error is
// returned even when the sender goroutine has already exited
func TestDomainOpenConsoleBidirectionalStreamError(t *testing.T) {
	l, srv := setupConsoleTest(t)

	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)

		hdr, _, err := srv.readPacket()
		if !assert.NoError(t, err) {
			return
		}
		serial, proc := hdr.Serial, hdr.Procedure

		if !assert.NoError(t, srv.writePacket(serial, proc, socket.Reply, socket.StatusOK, nil)) {
			return
		}

		// Empty stdin: the client immediately sends its stream-finish packet
		// and the sender goroutine exits.
		shdr, _, err := srv.readPacket()
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, uint32(socket.Stream), shdr.Type)
		assert.Equal(t, uint32(socket.StatusOK), shdr.Status)

		// Then the server fails the stream mid-session.
		assert.NoError(t, srv.writePacket(serial, proc, socket.Stream, socket.StatusError, testErrorMessage))
	}()

	var console bytes.Buffer
	done := make(chan error, 1)
	go func() {
		done <- l.DomainOpenConsoleBidirectional(Domain{}, nil, bytes.NewReader(nil), &console, 0)
	}()

	select {
	case err := <-done:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "domain is not running")
	case <-time.After(2 * time.Second):
		t.Fatal("DomainOpenConsoleBidirectional did not return after a stream error")
	}
	<-serverDone
}
