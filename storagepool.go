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
	"encoding/hex"
	"strings"

	"github.com/davecgh/go-xdr/xdr2"
	"github.com/digitalocean/go-libvirt/internal/constants"
)

// StoragePool represents a storage pool as seen by libvirt.
type StoragePool struct {
	Name string
	UUID [constants.UUIDSize]byte
	l    *Libvirt
}

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

	result.Pool.l = l
	return &result.Pool, nil
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

	result.Pool.l = l
	return &result.Pool, nil
}

// Refresh refreshes the storage pool.
func (p *StoragePool) Refresh(flags uint32) error {
	req := struct {
		Pool  StoragePool
		Flags uint32
	}{
		Pool:  *p,
		Flags: flags, // unused per libvirt source, callers should pass 0
	}

	buf, err := encode(&req)
	if err != nil {
		return err
	}

	resp, err := p.l.request(constants.ProcStoragePoolRefresh, constants.ProgramRemote, &buf)
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

	for _, p := range result.Pools {
		p.l = l
	}
	return result.Pools, nil
}

// SetAutostart set autostart for domain.
func (p *StoragePool) SetAutostart(autostart bool) error {
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

	resp, err := p.l.request(constants.ProcStoragePoolSetAutostart, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}

// StorageVolumeCreateXML creates a volume.
func (p *StoragePool) StorageVolumeCreateXML(x []byte, flags StorageVolumeCreateFlags) (*StorageVolume, error) {
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

	resp, err := p.l.request(constants.ProcStorageVolumeCreateXML, constants.ProgramRemote, &buf)
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
	result.Volume.l = p.l

	return &result.Volume, nil
}

// StorageVolumeLookupByName returns a volume as seen by libvirt.
func (p *StoragePool) StorageVolumeLookupByName(name string) (*StorageVolume, error) {
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

	resp, err := p.l.request(constants.ProcStorageVolumeLookupByName, constants.ProgramRemote, &buf)
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

	result.Volume.l = p.l
	return &result.Volume, nil
}
