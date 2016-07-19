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

// qemuError represents a QEMU process error.
type qemuError struct {
	Error struct {
		Class       string `json:"class"`
		Description string `json:"desc"`
	} `json:"error"`
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
