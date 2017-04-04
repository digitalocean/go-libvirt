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
	"bytes"

	"github.com/davecgh/go-xdr/xdr2"
	"github.com/digitalocean/go-libvirt/internal/constants"
)

// StorageVolume represents a volume as seen by libvirt.
type StorageVolume struct {
	Pool string
	Name string
	Key  string
}

// StorageVolumeCreateFlags specifies options when performing a volume creation.
type StorageVolumeCreateFlags uint32

const (
	_ StorageVolumeCreateFlags = iota
	// StorageVolumeCreateFlagPreallocMetadata preallocates metadata
	StorageVolumeCreateFlagPreallocMetadata

	// StorageVolumeCreateFlagReflink use btrfs light copy
	StorageVolumeCreateFlagReflink
)

// StorageVolumeDeleteFlags specifies options when performing a volume deletion.
type StorageVolumeDeleteFlags uint32

const (
	// StorageVolumeDeleteFlagNormal delete metadata only (fast)
	StorageVolumeDeleteFlagNormal StorageVolumeDeleteFlags = iota

	// StorageVolumeDeleteFlagZeroes clear all data to zeros (slow)
	StorageVolumeDeleteFlagZeroes

	// StorageVolumeDeleteFlagWithSnapshots force removal of volume, even if in use
	StorageVolumeDeleteFlagWithSnapshots
)

// StorageVolumeCreateXML creates a volume.
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

// StorageVolumeDelete deletes a volume.
func (l *Libvirt) StorageVolumeDelete(v *StorageVolume, flags StorageVolumeDeleteFlags) error {
	payload := struct {
		Vol   StorageVolume
		Flags StorageVolumeDeleteFlags
	}{
		Vol:   *v,
		Flags: flags,
	}

	buf, err := encode(&payload)
	if err != nil {
		return err
	}

	resp, err := l.request(constants.ProcStorageVolumeDelete, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
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
