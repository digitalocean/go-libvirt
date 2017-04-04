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

	"github.com/davecgh/go-xdr/xdr2"
	"github.com/digitalocean/go-libvirt/internal/constants"
)

// Network represents a network as seen by libvirt.
type Network struct {
	Name string
	UUID [constants.UUIDSize]byte
	l    *Libvirt
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

	result.Net.l = l
	return &result.Net, nil
}

// SetAutostart set autostart for network.
func (n *Network) SetAutostart(autostart bool) error {
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

	resp, err := n.l.request(constants.ProcNetworkSetAutostart, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}
