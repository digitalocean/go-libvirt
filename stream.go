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

// Stream represents a stream in ibvirt.
type Stream struct {
	l *Libvirt
}

// StreamNew creates new stream.
func (l *Libvirt) StreamNew() (*Stream, error) {
	return &Stream{l: l}, nil
}

// Abort forced closes stream.
func (s *Stream) Abort() error {
	return nil
}

// Finish shutdown stream.
func (s *Stream) Finish() error {
	return nil
}

// Read reads from stream
func (s *Stream) Read(p []byte) (int, error) {
	return 0, nil
}

// Write writes to stream
func (s *Stream) Write(p []byte) (int, error) {
	return 0, nil
}

// Close closes stream, to implement io.Closer interface
func (s *Stream) Close() error {
	return s.Finish()
}
