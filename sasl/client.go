package sasl

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	saslMD5NoneSize = 32
	saslDataMax     = 65536

	sealingClientServer = "Digest H(A1) to client-to-server sealing key magic constant"
	sealingServerClient = "Digest H(A1) to server-to-client sealing key magic constant"
	signingClientServer = "Digest session key to client-to-server signing key magic constant"
	signingServerClient = "Digest session key to server-to-client signing key magic constant"

	mechDigestMD5      = "DIGEST-MD5"
	methodAuthenticate = "AUTHENTICATE"
	charsetUTF8        = "utf-8"
	algorithmMD5Sess   = "md5-sess"

	qopAuth     = "auth"
	qopAuthInt  = "auth-int"
	qopAuthConf = "auth-conf"

	fieldAlgorithm = "algorithm"
	fieldDigestURI = "digest-uri"
	fieldQOP       = "qop"
	fieldCipher    = "cipher"
	fieldCharset   = "charset"
	fieldRealm     = "realm"
	fieldNonce     = "nonce"
	fieldCNonce    = "cnonce"
	fieldMaxBuf    = "maxbuf"
	fieldNC        = "nc"
	fieldRSPAuth   = "rspauth"
)

type saslInfoCallback func([]CallbackRequestItem) error

type context struct {
	challenge
	KIReceive []byte
	KISend    []byte
	EncKey    []byte
	DecKey    []byte
	HA1       []byte
}

type client struct {
	LocalAddress  net.Addr
	RemoteAddress net.Addr
	Callback      saslInfoCallback
	context       *context
}

type challenge struct {
	Service    string
	Hostname   string
	DigestURI  string
	Username   string
	AuthzID    string
	NonceHex   string
	CNonceHex  string
	NC         uint32
	Algorithm  string
	QOP        string
	Realm      string
	Charset    string
	MaxBuf     uint32
	Response   string
	Method     string
	SessionKey string
	Ciphers    []string
	Cipher     string
	CipherType cipherType
}

// Client represents the client portion of the SASL protocol.
type Client interface {
	FindMech(mechList []string) (mech string, err error)
	ApplyChallenge(challenge string)
	MakeResponse() string
	SetCNonce(string)
	GetContext() SecurityContext
}

// NewClient creates a new instance of the client.
func NewClient(service, hostname string, auth AuthInfo) Client {
	nonce := make([]byte, saslMD5NoneSize)
	rand.Read(nonce)
	return &client{
		Callback: auth.Callback,
		context: &context{
			challenge: challenge{
				Service:   service,
				Hostname:  hostname,
				CNonceHex: base64.StdEncoding.EncodeToString(nonce),
				MaxBuf:    saslDataMax,
				Charset:   charsetUTF8,
				Method:    methodAuthenticate,
			},
		},
	}
}

func (client *client) GetContext() SecurityContext {
	switch client.context.QOP {
	case "auth-int":
		ctx := &securityIntegrityContext{
			HA1:    client.context.HA1,
			MaxBuf: client.context.MaxBuf,
		}
		ctx.generateKeys()
		return ctx
	case "auth-conf":
		ctx := &securityPrivacyContext{
			securityIntegrityContext: securityIntegrityContext{
				HA1:    client.context.HA1,
				MaxBuf: client.context.MaxBuf,
			},
			cipher: client.context.CipherType,
		}
		ctx.generateKeys()
		return ctx
	}
	return nil
}

func (client *client) FindMech(mechList []string) (mech string, err error) {
	for _, entry := range mechList {
		// eventually this would be the place to switch to different
		// mech variations, but for now we only support md5
		if entry == mechDigestMD5 {
			return entry, nil
		}
	}
	return "", fmt.Errorf("sasl: None of requested mech (%v) supported", mechList)
}

func md5Sum(a, b []byte) []byte {
	h := md5.New()
	_, _ = h.Write(a)
	_, _ = h.Write(b)
	return h.Sum(nil)
}

func (client *client) ApplyChallenge(challenge string) {
	parseChallenge(challenge, &client.context.challenge)
}

func minUint32(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

func (client *client) calcSecretHash(password string) []byte {
	h := md5.New()
	_, _ = io.WriteString(h, client.context.Username)
	_, _ = io.WriteString(h, ":")
	_, _ = io.WriteString(h, client.context.Realm)
	_, _ = io.WriteString(h, ":")
	_, _ = io.WriteString(h, password)
	x := h.Sum(nil)
	return x[:]
}

func (client *client) calcHA1(password string) []byte {
	ha1 := md5.New()

	tmp := client.calcSecretHash(password)
	if client.context.Algorithm == algorithmMD5Sess {
		ha1.Write(tmp)
		_, _ = io.WriteString(ha1, ":")
		_, _ = io.WriteString(ha1, client.context.NonceHex)
		_, _ = io.WriteString(ha1, ":")
		_, _ = io.WriteString(ha1, client.context.CNonceHex)
		if client.context.AuthzID != "" {
			_, _ = io.WriteString(ha1, ":")
			_, _ = io.WriteString(ha1, client.context.AuthzID)
		}
	}
	client.context.HA1 = ha1.Sum(nil)
	client.context.SessionKey = hex.EncodeToString(client.context.HA1)
	return client.context.HA1
}

func (client *client) calcHA2() string {
	digest := md5.New()
	_, _ = io.WriteString(digest, client.context.Method)
	_, _ = io.WriteString(digest, ":")
	_, _ = io.WriteString(digest, client.context.DigestURI)
	if client.context.QOP != qopAuth {
		_, _ = io.WriteString(digest, ":00000000000000000000000000000000")
	}
	return hex.EncodeToString(digest.Sum(nil))
}

func (client *client) genResponse(username, password string) string {
	client.context.Username = username
	respHash := md5.New()
	client.calcHA1(password)
	_, _ = io.WriteString(respHash, client.context.SessionKey)
	_, _ = io.WriteString(respHash, ":")
	_, _ = io.WriteString(respHash, client.context.NonceHex)
	_, _ = io.WriteString(respHash, ":")

	if client.context.QOP != "" {
		// we need the exact formatting, so use Sprintf
		_, _ = io.WriteString(respHash, fmt.Sprintf("%.08x", client.context.NC))
		_, _ = io.WriteString(respHash, ":")
		_, _ = io.WriteString(respHash, client.context.CNonceHex)
		_, _ = io.WriteString(respHash, ":")
		_, _ = io.WriteString(respHash, client.context.QOP)
		_, _ = io.WriteString(respHash, ":")
	}
	_, _ = io.WriteString(respHash, client.calcHA2())
	return hex.EncodeToString(respHash.Sum(nil))
}

func (client *client) requestInfo() string {
	password := ""
	items := []CallbackRequestItem{
		CallbackRequestItem{
			Type:     CredentialPassphrase,
			Response: &password,
		},
	}

	if client.context.Username == "" {
		items = append(items, CallbackRequestItem{
			Type:     CredentialUsername,
			Response: &client.context.Username,
		})
	}

	if client.Callback(items) == nil {
		return password
	}
	return ""
}

func (client *client) MakeResponse() string {
	// increase NC usage count
	client.context.NC++

	// calculate password response hash
	password := client.requestInfo()

	ciphers := ""
	if client.context.QOP == qopAuthConf {
		// in case of auth-conf we have to let the server side know which cipher we
		// would like to use
		ciphers = fmt.Sprintf("cipher=%s,", client.context.Cipher)
	}

	return fmt.Sprintf(
		`username=%q,realm=%q,nonce="%s",cnonce="%s",nc=%.08x,qop=%s,maxbuf=%d,digest-uri=%q,%sresponse=%s`,
		client.context.Username,
		client.context.Realm,
		client.context.NonceHex,
		client.context.CNonceHex,
		client.context.NC,
		client.context.QOP,
		client.context.MaxBuf,
		client.context.DigestURI,
		ciphers,
		client.genResponse(client.context.Username, password))
}

func (client *client) SetCNonce(s string) {
	client.context.CNonceHex = s
}

func parseChallenge(serverChallenge string, ch *challenge) *challenge {
	if ch == nil {
		ch = &challenge{}
	}

	fieldMapper := map[string]func(string){
		fieldDigestURI: func(s string) {
			ch.DigestURI = s
		},
		fieldAlgorithm: func(s string) {
			ch.Algorithm = s
		},
		fieldQOP: func(s string) {
			ch.QOP = s
		},
		fieldCipher: func(s string) {
			ch.Ciphers = strings.Split(s, ",")
			for _, cipher := range ch.Ciphers {
				switch cipher {
				case "rc4":
					ch.CipherType = cipherRC4
				case "des":
					ch.CipherType = cipherDES
				case "rc4-56":
					ch.CipherType = cipherRC4_56
				case "rc4-40":
					ch.CipherType = cipherRC4_40
				case "3des":
					ch.CipherType = cipher3DES
				default:
					continue
				}
				ch.Cipher = cipher
			}
		},
		fieldCharset: func(s string) {
			ch.Charset = s
		},
		fieldRealm: func(s string) {
			ch.Realm = s
		},
		fieldNonce: func(s string) {
			ch.NonceHex = s
		},
		fieldCNonce: func(s string) {
			ch.CNonceHex = s
		},
		fieldMaxBuf: func(s string) {
			if v, e := strconv.ParseUint(s, 10, 32); e == nil {
				ch.MaxBuf = minUint32(ch.MaxBuf, uint32(v))
			}
		},
		fieldNC: func(s string) {
			if v, e := strconv.ParseUint("0x"+s, 16, 32); e == nil {
				ch.NC = uint32(v)
			}
		},
		fieldRSPAuth: func(string) {
			ch.Method = ""
		},
	}

	for _, field := range strings.Split(serverChallenge, ",") {
		kv := strings.SplitN(field, "=", 2)
		if len(kv) == 2 {
			kv[1] = strings.Trim(kv[1], `"`)
			if mapper, ok := fieldMapper[kv[0]]; ok {
				mapper(kv[1])
			}
		}
	}
	if ch.DigestURI == "" {
		ch.DigestURI = ch.Service + "/localhost"
	}

	return ch
}
