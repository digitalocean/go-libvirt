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

import "github.com/digitalocean/go-libvirt/internal/constants"

// StorageVolume represents a volume as seen by libvirt.
type StorageVolume struct {
	Pool string
	Name string
	Key  string
	l    *Libvirt
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

// Delete deletes a volume.
func (v *StorageVolume) Delete(flags StorageVolumeDeleteFlags) error {
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

	resp, err := v.l.request(constants.ProcStorageVolumeDelete, constants.ProgramRemote, &buf)
	if err != nil {
		return err
	}

	r := <-resp
	if r.Status != StatusOK {
		return decodeError(r.Payload)
	}

	return nil
}
