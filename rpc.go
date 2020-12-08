// Copyright 2018 The go-libvirt Authors.
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
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/digitalocean/go-libvirt/internal/constants"
	xdr "github.com/digitalocean/go-libvirt/internal/go-xdr/xdr2"
)

// ErrUnsupported is returned if a procedure is not supported by libvirt
var ErrUnsupported = errors.New("unsupported procedure requested")

// request and response types
const (
	// Call is used when making calls to the remote server.
	Call = iota

	// Reply indicates a server reply.
	Reply

	// Message is an asynchronous notification.
	Message

	// Stream represents a stream data packet.
	Stream

	// CallWithFDs is used by a client to indicate the request has
	// arguments with file descriptors.
	CallWithFDs

	// ReplyWithFDs is used by a server to indicate the request has
	// arguments with file descriptors.
	ReplyWithFDs
)

// request and response statuses
const (
	// StatusOK is always set for method calls or events.
	// For replies it indicates successful completion of the method.
	// For streams it indicates confirmation of the end of file on the stream.
	StatusOK = iota

	// StatusError for replies indicates that the method call failed
	// and error information is being returned. For streams this indicates
	// that not all data was sent and the stream has aborted.
	StatusError

	// StatusContinue is only used for streams.
	// This indicates that further data packets will be following.
	StatusContinue
)

// header is a libvirt rpc packet header
type header struct {
	// Program identifier
	Program uint32

	// Program version
	Version uint32

	// Remote procedure identifier
	Procedure uint32

	// Call type, e.g., Reply
	Type uint32

	// Call serial number
	Serial int32

	// Request status, e.g., StatusOK
	Status uint32
}

// packet represents a RPC request or response.
type packet struct {
	// Size of packet, in bytes, including length.
	// Len + Header + Payload
	Len    uint32
	Header header
}

// Global packet instance, for use with unsafe.Sizeof()
var _p packet

// internal rpc response
type response struct {
	Payload []byte
	Status  uint32
}

// typedParamDecoder is an empty struct with methods for handling libvirt's
// typed parameters, which are basically a union type.
type typedParamDecoder struct{}

// eventStream acts like a buffered channel with an unbounded buffer. The
// implementation consists of a pair of unbuffered channels and a goroutine to
// manage them. The eventStream can be cancelled by the client, in which case it will
// continue to receive events from libvirt until it can deregister itself. This
// prevents the rpc dispatcher from deadlocking by trying to send an event when
// the caller is not polling for one.
type eventStream struct {
	// Program specifies the source of the events - libvirt or QEMU.
	Program int32

	// CallbackID is the value returned by ConnectDomainCallbackRegisterAny
	CallbackID int32

	// Private members used to implement the unbounded channel behavior.
	fillCh  chan event
	drainCh chan event
	queue   []event
}

func newEventStream(ctx context.Context, program, callbackID int32) eventStream {
	ic := eventStream{
		Program:    program,
		CallbackID: callbackID,
		fillCh:     make(chan event),
		drainCh:    make(chan event),
		queue:      make([]event, 0),
	}

	// Start a goroutine to manage the queue
	go ic.process(ctx)

	return ic
}

// process is meant to be started in a separate goroutine. It will receive
// incoming events on one channel, and foward them to a second channel for a
// client to consume. Because of the event buffer, writes to the fill channel
// will not block.
func (ic *eventStream) process(ctx context.Context) {
	defer close(ic.drainCh)
	var sendEv event
Relay:
	for {
		if sendEv == nil {
			sendEv = ic.dequeue()
		}

		if sendEv == nil {
			select {
			case ev, ok := <-ic.fillCh:
				if !ok {
					// fillCh is closed; exit
					return
				}
				ic.enqueue(ev)
			case <-ctx.Done():
				break Relay
			}
		} else {
			select {
			case ev, ok := <-ic.fillCh:
				if !ok {
					// fillCh is closed; exit
					return
				}
				ic.enqueue(ev)
			case ic.drainCh <- sendEv:
				sendEv = nil
			case <-ctx.Done():
				break Relay
			}
		}
	}
	// Continue to drain & discard the incoming events until the fill channel
	// is closed.
	for ok := true; ok; _, ok = <-ic.fillCh {
	}
}

func (ic *eventStream) Close() {
	// Closing the fillCh will cause the goroutine started by Open to exit.
	close(ic.fillCh)
}

func (ic *eventStream) Send(ev event) {
	ic.fillCh <- ev
}

func (ic *eventStream) Recv() (event, bool) {
	ev, ok := <-ic.drainCh
	return ev, ok
}

func (ic *eventStream) enqueue(ev event) {
	ic.queue = append(ic.queue, ev)
}

func (ic *eventStream) dequeue() event {
	if len(ic.queue) == 0 {
		return nil
	}
	ev := ic.queue[0]
	ic.queue = ic.queue[1:]
	return ev
}

// libvirt error response
type libvirtError struct {
	Code     uint32
	DomainID uint32
	Padding  uint8
	Message  string
	Level    uint32
}

func (e libvirtError) Error() string {
	return e.Message
}

// checkError is used to check whether an error is a libvirtError, and if it is,
// whether its error code matches the one passed in. It will return false if
// these conditions are not met.
func checkError(err error, expectedError errorNumber) bool {
	e, ok := err.(libvirtError)
	if ok {
		return e.Code == uint32(expectedError)
	}
	return false
}

// IsNotFound detects libvirt's ERR_NO_DOMAIN.
func IsNotFound(err error) bool {
	return checkError(err, errNoDomain)
}

// listen processes incoming data and routes
// responses to their respective callback handler.
func (l *Libvirt) listen() {
	for {
		// response packet length
		length, err := pktlen(l.r)
		if err != nil {
			// When the underlying connection EOFs or is closed, stop
			// this goroutine
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			// invalid packet
			continue
		}

		// response header
		h, err := extractHeader(l.r)
		if err != nil {
			// invalid packet
			continue
		}

		// payload: packet length minus what was previously read
		size := int(length) - int(unsafe.Sizeof(_p))
		buf := make([]byte, size)
		_, err = io.ReadFull(l.r, buf)
		if err != nil {
			// invalid packet
			continue
		}

		// route response to caller
		l.route(h, buf)
	}
}

// callback sends rpc responses to their respective caller.
func (l *Libvirt) callback(id int32, res response) {
	l.cm.Lock()
	c, ok := l.callbacks[id]
	l.cm.Unlock()
	if ok {
		// we close the channel in deregister() so that we don't block here
		// forever without a receiver. If that happens, this write will panic.
		defer func() {
			recover()
		}()
		c <- res
	}
}

// TODO: This needs to be rewritten. The current code treats these lifecycle and
// qemu monitor events specially, but doens't require that the client have a
// goroutine running to collect them when they arrive. Without that, it's
// possible for the client to get stuck if it tries to read a request when this
// code is trying to send an event on the channel. Also, these are not the only
// event types supported by libvirt, just the ones someone has cared enough to
// add in here.

// route sends incoming packets to their listeners.
func (l *Libvirt) route(h *header, buf []byte) {
	// route events to their respective listener
	var streamEvent event
	switch {
	case h.Program == constants.QEMUProgram && h.Procedure == constants.QEMUProcDomainMonitorEvent:
		streamEvent = &DomainEvent{}
	case h.Program == constants.Program && h.Procedure == constants.ProcDomainEventCallbackLifecycle:
		streamEvent = &DomainEventCallbackLifecycleMsg{}
	}

	if streamEvent != nil {
		err := eventDecoder(buf, streamEvent)
		if err != nil { // event was malformed, drop.
			return
		}
		l.stream(streamEvent)
		return
	}

	// send responses to caller
	res := response{
		Payload: buf,
		Status:  h.Status,
	}
	l.callback(h.Serial, res)
}

// serial provides atomic access to the next sequential request serial number.
func (l *Libvirt) serial() int32 {
	return atomic.AddInt32(&l.s, 1)
}

// stream decodes domain events and sends them to the respective event listener.
func (l *Libvirt) stream(e event) {
	// Hold the lock while sending. The eventStream queue implementation
	// guarantees that this send will not block; holding the lock avoids the
	// possibility that the channel will be closed by a deregister call while
	// we're sending.
	l.em.Lock()
	c, ok := l.events[e.GetCallbackID()]
	defer l.em.Unlock()

	if ok {
		c.Send(e)
	}
}

// addStream configures the routing for an event stream.
func (l *Libvirt) addStream(s eventStream) {
	l.em.Lock()
	l.events[s.CallbackID] = s
	l.em.Unlock()
}

// removeStream removes the event stream from the internal list we use to
// dispatch events to the client. The caller should first unsubscribe from the
// event via a libvirt call.
func (l *Libvirt) removeStream(id int32) error {
	stream := l.events[id]

	l.em.Lock()
	delete(l.events, id)
	// We can now close the events stream. It is safe to do this because we are
	// holding the em lock. Since the route() call also must hold this lock when
	// looking up an eventStream and queueing an event to it, this Close() call
	// won't cause a panic in that code.
	stream.Close()
	l.em.Unlock()

	return nil
}

// register configures a method response callback
func (l *Libvirt) register(id int32, c chan response) {
	l.cm.Lock()
	l.callbacks[id] = c
	l.cm.Unlock()
}

// deregister destroys a method response callback
func (l *Libvirt) deregister(id int32) {
	l.cm.Lock()
	if _, ok := l.callbacks[id]; ok {
		close(l.callbacks[id])
		delete(l.callbacks, id)
	}
	l.cm.Unlock()
}

// deregisterAll closes all the waiting callback channels. This is used to clean
// up if the connection to libvirt is lost. Callers waiting for responses will
// return an error when the response channel is closed, rather than just
// hanging.
func (l *Libvirt) deregisterAll() {
	l.cm.Lock()
	for id := range l.callbacks {
		// can't call deregister() here because we're already holding the lock.
		close(l.callbacks[id])
		delete(l.callbacks, id)
	}
	l.cm.Unlock()
}

// request performs a libvirt RPC request.
// returns response returned by server.
// if response is not OK, decodes error from it and returns it.
func (l *Libvirt) request(proc uint32, program uint32, payload []byte) (response, error) {
	return l.requestStream(proc, program, payload, nil, nil)
}

// requestStream performs a libvirt RPC request. The outStream and inStream
// parameters are optional, and should be nil for RPC endpoints that don't
// return a stream.
func (l *Libvirt) requestStream(proc uint32, program uint32, payload []byte,
	outStream io.Reader, inStream io.Writer) (response, error) {
	serial := l.serial()
	c := make(chan response)

	l.register(serial, c)
	defer l.deregister(serial)

	err := l.sendPacket(serial, proc, program, payload, Call, StatusOK)
	if err != nil {
		return response{}, err
	}

	resp, err := l.getResponse(c)
	if err != nil {
		return resp, err
	}

	if outStream != nil {
		abortOutStream := make(chan bool)
		outStreamErr := make(chan error)
		go func() {
			outStreamErr <- l.sendStream(serial, proc, program, outStream, abortOutStream)
		}()

		// Even without incoming stream server sends confirmation once all data is received
		resp, err = l.processIncomingStream(c, inStream)
		if err != nil {
			abortOutStream <- true
			return resp, err
		}

		err = <-outStreamErr
		if err != nil {
			return response{}, err
		}
	}

	if inStream != nil {
		return l.processIncomingStream(c, inStream)
	}

	return resp, nil
}

// processIncomingStream is called once we've successfully sent a request to
// libvirt. It writes the responses back to the stream passed by the caller
// until libvirt sends a packet with statusOK or an error.
func (l *Libvirt) processIncomingStream(c chan response, inStream io.Writer) (response, error) {
	for {
		resp, err := l.getResponse(c)
		if err != nil {
			return resp, err
		}
		// StatusOK here means end of stream
		if resp.Status == StatusOK {
			return resp, nil
		}
		// StatusError is handled in getResponse, so this is StatusContinue
		// StatusContinue is valid here only for stream packets
		// libvirtd breaks protocol and returns StatusContinue with empty Payload when stream finishes
		if len(resp.Payload) == 0 {
			return resp, nil
		}
		if inStream != nil {
			_, err = inStream.Write(resp.Payload)
			if err != nil {
				return response{}, err
			}
		}
	}
}

func (l *Libvirt) sendStream(serial int32, proc uint32, program uint32, stream io.Reader, abort chan bool) error {
	// Keep total packet length under 4 MiB to follow possible limitation in libvirt server code
	buf := make([]byte, 4*MiB-unsafe.Sizeof(_p))
	for {
		select {
		case <-abort:
			return l.sendPacket(serial, proc, program, nil, Stream, StatusError)
		default:
		}
		n, err := stream.Read(buf)
		if n > 0 {
			err2 := l.sendPacket(serial, proc, program, buf[:n], Stream, StatusContinue)
			if err2 != nil {
				return err2
			}
		}
		if err != nil {
			if err == io.EOF {
				return l.sendPacket(serial, proc, program, nil, Stream, StatusOK)
			}
			// keep original error
			err2 := l.sendPacket(serial, proc, program, nil, Stream, StatusError)
			if err2 != nil {
				return err2
			}
			return err
		}
	}
}

func (l *Libvirt) sendPacket(serial int32, proc uint32, program uint32, payload []byte, typ uint32, status uint32) error {

	p := packet{
		Header: header{
			Program:   program,
			Version:   constants.ProtocolVersion,
			Procedure: proc,
			Type:      typ,
			Serial:    serial,
			Status:    status,
		},
	}

	size := int(unsafe.Sizeof(p.Len)) + int(unsafe.Sizeof(p.Header))
	if payload != nil {
		size += len(payload)
	}
	p.Len = uint32(size)

	// write header
	l.mu.Lock()
	defer l.mu.Unlock()
	err := binary.Write(l.w, binary.BigEndian, p)
	if err != nil {
		return err
	}

	// write payload
	if payload != nil {
		err = binary.Write(l.w, binary.BigEndian, payload)
		if err != nil {
			return err
		}
	}

	return l.w.Flush()
}

func (l *Libvirt) getResponse(c chan response) (response, error) {
	resp := <-c
	if resp.Status == StatusError {
		return resp, decodeError(resp.Payload)
	}

	return resp, nil
}

// encode XDR encodes the provided data.
func encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	_, err := xdr.Marshal(&buf, data)

	return buf.Bytes(), err
}

// decodeError extracts an error message from the provider buffer.
func decodeError(buf []byte) error {
	var e libvirtError

	dec := xdr.NewDecoder(bytes.NewReader(buf))
	_, err := dec.Decode(&e)
	if err != nil {
		return err
	}

	if strings.Contains(e.Message, "unknown procedure") {
		return ErrUnsupported
	}
	// if libvirt returns ERR_OK, ignore the error
	if checkError(e, errOk) {
		return nil
	}

	return e
}

// eventDecoder decoder an event from a xdr buffer.
func eventDecoder(buf []byte, e interface{}) error {
	dec := xdr.NewDecoder(bytes.NewReader(buf))
	_, err := dec.Decode(e)
	return err
}

// pktlen determines the length of an incoming rpc response.
// If an error is encountered reading the provided Reader, the
// error is returned and response length will be 0.
func pktlen(r io.Reader) (uint32, error) {
	buf := make([]byte, unsafe.Sizeof(_p.Len))

	// read exactly constants.PacketLengthSize bytes
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(buf), nil
}

// extractHeader returns the decoded header from an incoming response.
func extractHeader(r io.Reader) (*header, error) {
	buf := make([]byte, unsafe.Sizeof(_p.Header))

	// read exactly constants.HeaderSize bytes
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	h := &header{
		Program:   binary.BigEndian.Uint32(buf[0:4]),
		Version:   binary.BigEndian.Uint32(buf[4:8]),
		Procedure: binary.BigEndian.Uint32(buf[8:12]),
		Type:      binary.BigEndian.Uint32(buf[12:16]),
		Serial:    int32(binary.BigEndian.Uint32(buf[16:20])),
		Status:    binary.BigEndian.Uint32(buf[20:24]),
	}

	return h, nil
}

// Decode decodes a TypedParam. These are part of the libvirt spec, and not xdr
// proper. TypedParams contain a name, which is called Field for some reason,
// and a Value, which itself has a "discriminant" - an integer enum encoding the
// actual type, and a value, the length of which varies based on the actual
// type.
func (tpd typedParamDecoder) Decode(d *xdr.Decoder, v reflect.Value) (int, error) {
	// Get the name of the typed param first
	name, n, err := d.DecodeString()
	if err != nil {
		return n, err
	}
	val, n2, err := tpd.decodeTypedParamValue(d)
	n += n2
	if err != nil {
		return n, err
	}
	tp := &TypedParam{Field: name, Value: *val}
	v.Set(reflect.ValueOf(*tp))

	return n, nil
}

// decodeTypedParamValue decodes the Value part of a TypedParam.
func (typedParamDecoder) decodeTypedParamValue(d *xdr.Decoder) (*TypedParamValue, int, error) {
	// All TypedParamValues begin with a uint32 discriminant that tells us what
	// type they are.
	discriminant, n, err := d.DecodeUint()
	if err != nil {
		return nil, n, err
	}
	var n2 int
	var tpv *TypedParamValue
	switch discriminant {
	case 1:
		var val int32
		n2, err = d.Decode(&val)
		tpv = &TypedParamValue{D: discriminant, I: val}
	case 2:
		var val uint32
		n2, err = d.Decode(&val)
		tpv = &TypedParamValue{D: discriminant, I: val}
	case 3:
		var val int64
		n2, err = d.Decode(&val)
		tpv = &TypedParamValue{D: discriminant, I: val}
	case 4:
		var val uint64
		n2, err = d.Decode(&val)
		tpv = &TypedParamValue{D: discriminant, I: val}
	case 5:
		var val float64
		n2, err = d.Decode(&val)
		tpv = &TypedParamValue{D: discriminant, I: val}
	case 6:
		var val int32
		n2, err = d.Decode(&val)
		tpv = &TypedParamValue{D: discriminant, I: val}
	case 7:
		var val string
		n2, err = d.Decode(&val)
		tpv = &TypedParamValue{D: discriminant, I: val}

	default:
		err = fmt.Errorf("invalid parameter type %v", discriminant)
	}
	n += n2

	return tpv, n, err
}
