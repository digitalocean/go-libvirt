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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/digitalocean/go-libvirt/libvirttest"
)

func TestConnect(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Connect()
	if err != nil {
		t.Error(err)
	}
}

func TestDisconnect(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	err := l.Disconnect()
	if err != nil {
		t.Error(err)
	}
}

func TestDomains(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	domains, err := l.Domains()
	if err != nil {
		t.Error(err)
	}

	wantLen := 2
	gotLen := len(domains)
	if gotLen != wantLen {
		t.Errorf("expected %d domains to be returned, got %d", wantLen, gotLen)
	}

	for i, d := range domains {
		wantID := i + 1
		if d.ID != wantID {
			t.Errorf("expected domain ID %q, got %q", wantID, d.ID)
		}

		wantName := fmt.Sprintf("aaaaaaa-%d", i+1)
		if d.Name != wantName {
			t.Errorf("expected domain name %q, got %q", wantName, d.Name)
		}
	}
}

func TestEvents(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)
	done := make(chan struct{})

	stream, err := l.Events("test")
	if err != nil {
		t.Error(err)
	}

	go func() {
		var e DomainEvent
		select {
		case e = <-stream:
		case <-time.After(time.Second * 5):
			t.Error("expected event, received timeout")
		}

		result := struct {
			Device string `json:"device"`
			Len    int    `json:"len"`
			Offset int    `json:"offset"`
			Speed  int    `json:"speed"`
			Type   string `json:"type"`
		}{}

		if err := json.Unmarshal(e.Details, &result); err != nil {
			t.Error(err)
		}

		expected := "drive-ide0-0-0"
		if result.Device != expected {
			t.Errorf("expected device %q, got %q", expected, result.Device)
		}

		done <- struct{}{}
	}()

	// send an event to the listener goroutine
	conn.Test.Write(append(testEventHeader, testEvent...))

	// wait for completion
	<-done
}

func TestRun(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	res, err := l.Run("test", []byte(`{"query-version"}`))
	if err != nil {
		t.Error(err)
	}

	type version struct {
		Return struct {
			Package string `json:"package"`
			QEMU    struct {
				Major int `json:"major"`
				Micro int `json:"micro"`
				Minor int `json:"minor"`
			} `json:"qemu"`
		} `json:"return"`
	}

	var v version
	err = json.Unmarshal(res, &v)
	if err != nil {
		t.Error(err)
	}

	expected := 2
	if v.Return.QEMU.Major != expected {
		t.Errorf("expected qemu major version %d, got %d", expected, v.Return.QEMU.Major)
	}
}

func TestRunFail(t *testing.T) {
	conn := libvirttest.New()
	conn.Fail = true
	l := New(conn)

	_, err := l.Run("test", []byte(`{"drive-foo"}`))
	if err == nil {
		t.Error("expected qemu error")
	}
}

func TestVersion(t *testing.T) {
	conn := libvirttest.New()
	l := New(conn)

	version, err := l.Version()
	if err != nil {
		t.Error(err)
	}

	expected := "1.3.4"
	if version != expected {
		t.Errorf("expected version %q, got %q", expected, version)
	}
}
