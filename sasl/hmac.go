package sasl

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rc4"
	"crypto/subtle"
	"encoding/binary"
	"fmt"
	"io"
)

var (
	errEmptyBuffer                = fmt.Errorf("ERROR: Empty buffer")
	errBufferTooSmall             = fmt.Errorf("ERROR: Given buffer is too small")
	errPeerSequenceNumberMismatch = fmt.Errorf("ERROR: Peer sequence number mismatch")
	errInvalidPeerHMAC            = fmt.Errorf("ERROR: Invalid peer HMAC")
)

type cipherType int

const (
	cipherNone cipherType = iota
	cipherRC4_40
	cipherRC4_56
	cipherRC4
	cipherDES
	cipher3DES
	_cipherCount
)

// SecurityContext wraps and unwraps data in case of added
// security layers negotiated during SASL authentication
type SecurityContext interface {
	Unwrap([]byte) ([]byte, error)
	Wrap([]byte) ([]byte, error)
	MaxBufSize() uint32
}

type securityNullContext struct {
}

type securityContextReadWriter struct {
	Context SecurityContext
	r       *bufio.Reader
	buffer  *bytes.Reader
	w       *bufio.Writer
}

// WrapReadWriteSecurityContext wraps an io.Reader and an io.Writer to apply the requested security
// mechanism by the server on the communication layer. The resulting io.Reader and io.Writer must be used
// for the communication from that time on.
func WrapReadWriteSecurityContext(r io.Reader, w io.Writer, context SecurityContext) (io.Reader, io.Writer) {
	if context == nil {
		return r, w
	}

	wrapper := &securityContextReadWriter{
		Context: context,
		r:       bufio.NewReader(r),
		w:       bufio.NewWriter(w),
	}
	return wrapper, wrapper
}

func (secrw *securityContextReadWriter) Read(out []byte) (int, error) {
	if secrw.buffer == nil || secrw.buffer.Len() == 0 {
		b := make([]byte, secrw.Context.MaxBufSize())
		cnt, err := secrw.r.Read(b)
		if err != nil {
			return 0, err
		}
		b, err = secrw.Context.Unwrap(b[0:cnt])
		if err != nil {
			return 0, err
		}
		secrw.buffer = bytes.NewReader(b)
	}
	cnt, err := secrw.buffer.Read(out)
	if err != nil {
		return 0, err
	}
	if secrw.buffer.Len() == 0 {
		secrw.buffer = nil
	}
	return cnt, nil
}

func (secrw *securityContextReadWriter) Write(in []byte) (int, error) {
	buf, err := secrw.Context.Wrap(in)
	if err != nil {
		return 0, err
	}
	towrite := len(buf)
	for towrite > 0 {
		written, err := secrw.w.Write(buf)
		if err != nil {
			return 0, err
		}
		towrite -= written
	}
	err = secrw.w.Flush()
	if err != nil {
		return 0, err
	}
	return len(in), nil
}

func (sec *securityNullContext) Unwrap(data []byte) ([]byte, error) {
	return data, nil
}

func (sec *securityNullContext) Wrap(data []byte) ([]byte, error) {
	return data, nil
}

func (sec *securityNullContext) MaxBufSize() uint32 {
	return 0
}

type securityIntegrityContext struct {
	securityNullContext
	SeqNum     uint32
	PeerSeqNum uint32
	Ki         []byte
	PeerKi     []byte
	EncKey     []byte
	DecKey     []byte
	HA1        []byte
	MaxBuf     uint32
}

type securityPrivacyContext struct {
	securityIntegrityContext
	cipher    cipherType
	blockSize int
	EncCipher func(dst, src []byte)
	DecCipher func(dst, src []byte)
}

func (sec *securityIntegrityContext) MaxBufSize() uint32 {
	return sec.MaxBuf
}

func (sec *securityPrivacyContext) generateKeys() error {
	ha1len := 16 // full hash length is the default
	if sec.cipher == cipherRC4_40 {
		ha1len = 5 // only the first 5 bytes are used
	} else if sec.cipher == cipherRC4_56 {
		ha1len = 7 // only the first 7 bytes are used
	}

	sec.EncKey = md5Sum(sec.HA1[0:ha1len], []byte(sealingClientServer))
	sec.DecKey = md5Sum(sec.HA1[0:ha1len], []byte(sealingServerClient))

	sec.Ki = md5Sum(sec.HA1, []byte(signingClientServer))
	sec.PeerKi = md5Sum(sec.HA1, []byte(signingServerClient))

	switch sec.cipher {
	case cipherRC4_40, cipherRC4_56, cipherRC4:
		sec.blockSize = 1
		cipher, err := rc4.NewCipher(sec.DecKey)
		if err != nil {
			return err
		}
		sec.DecCipher = cipher.XORKeyStream
		cipher, err = rc4.NewCipher(sec.EncKey)
		if err != nil {
			return err
		}
		sec.EncCipher = cipher.XORKeyStream
	case cipher3DES:
		block, err := des.NewTripleDESCipher(sec.DecKey)
		if err != nil {
			return err
		}
		sec.blockSize = block.BlockSize()
		sec.DecCipher = block.Decrypt
		block, err = des.NewTripleDESCipher(sec.EncKey)
		if err != nil {
			return err
		}
		sec.EncCipher = block.Encrypt
	case cipherDES:
		block, err := des.NewCipher(sec.DecKey)
		if err != nil {
			return err
		}
		sec.blockSize = block.BlockSize()
		sec.DecCipher = cipher.NewCTR(block, sec.DecKey[len(sec.DecKey)-8:]).XORKeyStream
		block, err = des.NewCipher(sec.EncKey)
		if err != nil {
			return err
		}
		sec.EncCipher = cipher.NewCTR(block, sec.EncKey[len(sec.EncKey)-8:]).XORKeyStream
	default:
		// this shouldn't happen
		return fmt.Errorf("no such cipher type %d", sec.cipher)
	}
	return nil
}

func (sec *securityIntegrityContext) generateKeys() error {
	sec.EncKey = md5Sum(sec.HA1, []byte(sealingClientServer))
	sec.DecKey = md5Sum(sec.HA1, []byte(sealingServerClient))
	sec.Ki = md5Sum(sec.HA1, []byte(signingClientServer))
	sec.PeerKi = md5Sum(sec.HA1, []byte(signingServerClient))
	return nil
}

func (sec *securityIntegrityContext) getHMAC(msg []byte) []byte {
	x := hmac.New(md5.New, sec.Ki)
	_ = binary.Write(x, binary.BigEndian, sec.SeqNum)
	_, _ = x.Write(msg)
	return x.Sum(nil)[0:10]
}

func (sec *securityIntegrityContext) getPeerHMAC(msg []byte) []byte {
	x := hmac.New(md5.New, sec.PeerKi)
	_ = binary.Write(x, binary.BigEndian, sec.PeerSeqNum)
	_, _ = x.Write(msg)
	return x.Sum(nil)[0:10]
}

func (sec *securityIntegrityContext) Unwrap(data []byte) (out []byte, err error) {
	if data == nil || len(data) <= 16 {
		return nil, errBufferTooSmall
	}
	sec.PeerSeqNum++
	b := bytes.NewBuffer(data[len(data)-16:])
	hmac := make([]byte, 10)
	if _, err := b.Read(hmac); err != nil {
		return nil, err
	}

	msg := make([]byte, 2)
	if _, err := b.Read(msg); err != nil {
		return nil, err
	}

	var seqNum uint32
	if err := binary.Read(b, binary.BigEndian, &seqNum); err != nil {
		return nil, err
	}

	if seqNum != sec.PeerSeqNum {
		return nil, errPeerSequenceNumberMismatch
	}

	data = data[0 : len(data)-16]
	expectedHMAC := sec.getPeerHMAC(data)
	if subtle.ConstantTimeCompare(expectedHMAC, hmac) != 1 {
		return nil, errInvalidPeerHMAC
	}

	return data, nil
}

func (sec *securityIntegrityContext) Wrap(data []byte) (out []byte, err error) {
	if data == nil || len(data) == 0 {
		return nil, errEmptyBuffer
	}
	sec.SeqNum++
	b := bytes.NewBuffer(data)
	_, err = b.Write(sec.getHMAC(data))
	if err != nil {
		return nil, err
	}

	_, err = b.Write(data[0:1])
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, sec.SeqNum)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (sec *securityPrivacyContext) Unwrap(data []byte) (out []byte, err error) {
	if data == nil || len(data) <= 4 {
		return nil, errBufferTooSmall
	}

	b := bytes.NewBuffer(data)

	///   n      0-7      10      2         4
	// DATA + PADDING + HMAC + VERSION + SEQNUM

	padsize := 0
	if sec.blockSize > 1 {
		// There's padding involved
		padsize = int(data[len(data)-17])
	}

	tmp := b.Bytes()

	sec.DecCipher(tmp[:len(tmp)-6], tmp[:len(tmp)-6])

	out = make([]byte, len(data)-(16+padsize))
	_, err = b.Read(out)
	if err != nil {
		return nil, err
	}

	expectedHMAC := sec.getPeerHMAC(out)

	hmac := make([]byte, 10)
	_, err = b.Read(hmac)
	if err != nil {
		return nil, err
	}

	ver := uint16(0)
	err = binary.Read(b, binary.BigEndian, &ver)
	if err != nil {
		return nil, err
	}
	if ver != 1 {
		return nil, fmt.Errorf("Invalid version number")
	}

	seqnum := uint32(0)
	err = binary.Read(b, binary.BigEndian, &seqnum)
	if err != nil {
		return nil, err
	}

	if sec.PeerSeqNum != seqnum {
		return nil, fmt.Errorf("Invalid sequence number")
	}

	if subtle.ConstantTimeCompare(hmac, expectedHMAC) != 1 {
		return nil, fmt.Errorf("Invalid signature: \nExpected: % x\nGot:      % x", expectedHMAC, hmac)
	}

	sec.PeerSeqNum++
	return out, nil
}

func (sec *securityPrivacyContext) Wrap(data []byte) (out []byte, err error) {
	if data == nil || len(data) == 0 {
		return nil, errEmptyBuffer
	}

	length := len(data)
	enclength := length + 10 // 10 bytes of HMAC
	padding := []byte{}
	if sec.blockSize > 1 {
		pad := byte(sec.blockSize - (enclength % sec.blockSize))
		padding = bytes.Repeat([]byte{pad}, int(pad))
		enclength += int(pad)
	}

	b := bytes.NewBuffer(nil)

	// write the length of the data that consist of:
	// enclength := Length(Actual Data) + Len(padding) + len(HMAC)=10
	// + 2 Bytes version + 4 bytes SeqNum
	err = binary.Write(b, binary.BigEndian, uint32(enclength+4+2))
	if err != nil {
		return nil, err
	}

	// actual data
	_, err = b.Write(data)
	if err != nil {
		return nil, err
	}

	// write the padding
	_, err = b.Write(padding)
	if err != nil {
		return nil, err
	}

	// calculate the hmac digest (10 bytes)
	hmac := sec.getHMAC(data)

	// write the hmac digest
	_, err = b.Write(hmac)
	if err != nil {
		return nil, err
	}

	// now we encrypt the data directly in the buffer
	// and only for the data + padding + hmac (skipping the size at the beginning)
	sec.EncCipher(b.Bytes()[4:], b.Bytes()[4:])

	// write the magic version constant which is always 1
	err = binary.Write(b, binary.BigEndian, uint16(1))
	if err != nil {
		return nil, err
	}

	// write the sequence number
	err = binary.Write(b, binary.BigEndian, sec.SeqNum)
	if err != nil {
		return nil, err
	}

	// increate the sequence number
	sec.SeqNum++
	return b.Bytes(), nil
}
