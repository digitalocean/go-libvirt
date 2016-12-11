package sasl

// CredentialType represents required data requested by the SASL implementation
type CredentialType int32

const (
	// CredentialUsername is used when the username is requested
	CredentialUsername CredentialType = iota
	// CredentialPassphrase is used when the password is requested
	CredentialPassphrase

//  So far unsupported but listed for completness
//  CredentialAuthname
//	CredentialLanguage
//	CredentialRealm
//	CredentialEchoPrompt
//	CredentialNoEchoPrompt
// 	CredentialCNonce
)

// AuthInfo represents the SASL callback instance for authentication data
type AuthInfo interface {
	Callback([]CallbackRequestItem) error
}

type defaultAuthInfo struct {
	username, password string
}

func (auth defaultAuthInfo) Callback(items []CallbackRequestItem) error {
	for _, item := range items {
		switch item.Type {
		case CredentialUsername:
			*item.Response = auth.username
		case CredentialPassphrase:
			*item.Response = auth.password
		default:
			item.Response = nil
		}
	}
	return nil
}

// MakeUserPassAuthInfo creates an AuthInfo implementation
// with embedded username and password that implements the
// required Callback method
func MakeUserPassAuthInfo(username, password string) AuthInfo {
	return &defaultAuthInfo{
		username: username,
		password: password,
	}
}

// CallbackRequestItem represents the requested information
// where Type says which type of information is requested
// and Reponse is a pointer to a string that receives the
// string through dereferencing and assinging to it
type CallbackRequestItem struct {
	Type     CredentialType
	Response *string
}
