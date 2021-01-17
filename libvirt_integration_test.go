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

// +build integration

package libvirt

import (
	"bytes"
	"encoding/xml"
	"net"
	"sync"
	"testing"
	"time"
)

// In order for this test to work, libvirtd must be running and listening for
// tcp connections.
const testAddr = "127.0.0.1:16509"

func TestConnectIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Error(err)
	}
}

func TestConnectToURIIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.ConnectToURI(TestDefault); err != nil {
		t.Error(err)
	}
}

func TestDisconnectIntegration(t *testing.T) {
	l := New(testConn(t))
	if err := l.Disconnect(); err != nil {
		t.Error(err)
	}
}

func TestCapabilities(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Fatal(err)
	}

	resp, err := l.Capabilities()
	if err != nil {
		t.Fatal(err)
	}

	// verify UUID exists within returned XML
	var caps struct {
		Host struct {
			UUID string `xml:"uuid"`
		} `xml:"host"`
	}

	if err := xml.Unmarshal(resp, &caps); err != nil {
		t.Fatal(err)
	}

	if caps.Host.UUID == "" {
		t.Error("expected capabilities to contain a UUID")
	}
}

func TestSecretsIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Fatal(err)
	}

	secrets, err := l.Secrets()
	if err != nil {
		t.Fatal(err)
	}

	wantLen := 1
	gotLen := len(secrets)
	if gotLen != wantLen {
		t.Fatalf("expected %d secrets, got %d", wantLen, gotLen)
	}

	s := secrets[0]

	wantType := SecretUsageTypeVolume
	if s.UsageType != int32(wantType) {
		t.Errorf("expected usage type: %d, got %d", wantType, s.UsageType)
	}

	wantID := "/tmp"
	if s.UsageID != wantID {
		t.Errorf("expected usage id: %q, got %q", wantID, s.UsageID)
	}

	// 19fdc2f2-fa64-46f3-bacf-42a8aafca6dd
	wantUUID := [UUIDBuflen]byte{
		0x19, 0xfd, 0xc2, 0xf2, 0xfa, 0x64, 0x46, 0xf3,
		0xba, 0xcf, 0x42, 0xa8, 0xaa, 0xfc, 0xa6, 0xdd,
	}
	if s.UUID != wantUUID {
		t.Errorf("expected UUID %q, got %q", wantUUID, s.UUID)
	}
}

func TestStoragePoolIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Fatal(err)
	}

	wantName := "test"
	pool, err := l.StoragePool(wantName)
	if err != nil {
		t.Fatal(err)
	}

	gotName := pool.Name
	if gotName != wantName {
		t.Errorf("expected name %q, got %q", wantName, gotName)
	}
}

func TestStoragePoolInvalidIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Fatal(err)
	}

	_, err := l.StoragePool("test-does-not-exist")
	if err == nil {
		t.Errorf("expected non-existent storage pool return error")
	}
}

func TestStoragePoolsIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Fatal(err)
	}

	pools, err := l.StoragePools(ConnectListStoragePoolsActive)
	if err != nil {
		t.Error(err)
	}

	wantLen := 1
	gotLen := len(pools)
	if gotLen != wantLen {
		t.Fatalf("expected %d storage pool, got %d", wantLen, gotLen)
	}

	wantName := "test"
	gotName := pools[0].Name
	if gotName != wantName {
		t.Errorf("expected name %q, got %q", wantName, gotName)
	}
}

func TestStoragePoolsAutostartIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Fatal(err)
	}

	pools, err := l.StoragePools(ConnectListStoragePoolsAutostart)
	if err != nil {
		t.Error(err)
	}

	wantLen := 0
	gotLen := len(pools)
	if gotLen != wantLen {
		t.Errorf("expected %d storage pool, got %d", wantLen, gotLen)
	}
}

func TestStoragePoolRefreshIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Fatal(err)
	}

	pool, err := l.StoragePool("test")
	if err != nil {
		t.Error(err)
	}

	err = l.StoragePoolRefresh(pool, 0)
	if err != nil {
		t.Error(err)
	}
}

func TestStoragePoolRefreshInvalidIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
		t.Fatal(err)
	}

	pool, err := l.StoragePool("test-does-not-exist")
	if err == nil {
		t.Error(err)
	}

	err = l.StoragePoolRefresh(pool, 0)
	if err == nil {
		t.Error("expected non-existent storage pool to fail refresh")
	}
}

func TestXMLIntegration(t *testing.T) {
	l := New(testConn(t))

	if err := l.Connect(); err != nil {
		t.Error(err)
	}
	defer l.Disconnect()

	var flags DomainXMLFlags
	data, err := l.XML("test", flags)
	if err != nil {
		t.Fatal(err)
	}

	var v interface{}
	if err := xml.Unmarshal(data, &v); err != nil {
		t.Error(err)
	}
}

func TestVolumeUploadDownloadIntegration(t *testing.T) {
	testdata := []byte("Hello, world!")
	l := New(testConn(t))

	if err := l.Connect(); err != nil {
		t.Error(err)
	}
	defer l.Disconnect()

	pool, err := l.StoragePool("test")
	if err != nil {
		t.Fatal(err)
	}

	var volObj struct {
		XMLName  xml.Name `xml:"volume"`
		Name     string   `xml:"name"`
		Capacity struct {
			Value uint64 `xml:",chardata"`
		} `xml:"capacity"`
		Target struct {
			Format struct {
				Type string `xml:"type,attr"`
			} `xml:"format"`
		} `xml:"target"`
	}
	volObj.Name = "testvol"
	volObj.Capacity.Value = uint64(len(testdata))
	volObj.Target.Format.Type = "raw"
	xmlVol, err := xml.Marshal(volObj)
	if err != nil {
		t.Fatal(err)
	}
	vol, err := l.StorageVolCreateXML(pool, string(xmlVol), 0)
	if err != nil {
		t.Fatal(err)
	}
	defer l.StorageVolDelete(vol, 0)
	err = l.StorageVolUpload(vol, bytes.NewBuffer(testdata), 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	err = l.StorageVolDownload(vol, &buf, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(testdata, buf.Bytes()) != 0 {
		t.Fatal("download not what we uploaded")
	}
}

// verify we're able to concurrently communicate with libvirtd.
// see: https://github.com/digitalocean/go-libvirt/pull/52
func Test_concurrentWrite(t *testing.T) {
	l := New(testConn(t))

	if err := l.Connect(); err != nil {
		t.Error(err)
	}
	defer l.Disconnect()

	count := 10
	wg := sync.WaitGroup{}
	wg.Add(count)

	start := make(chan struct{})
	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			<-start

			_, err := l.Domains()
			if err != nil {
				t.Fatal(err)
			}
		}()
	}

	close(start)

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for execution to complete")
	}
}

func testConn(t *testing.T) net.Conn {
	conn, err := net.DialTimeout("tcp", testAddr, time.Second*2)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}
