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
	"encoding/xml"
	"net"
	"testing"
	"time"
)

const testAddr = "127.0.0.1:16509"

func TestConnectIntegration(t *testing.T) {
	l := New(testConn(t))
	defer l.Disconnect()

	if err := l.Connect(); err != nil {
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

func testConn(t *testing.T) net.Conn {
	conn, err := net.DialTimeout("tcp", testAddr, time.Second*2)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}
