package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var sessionResponseBody = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><sessions><Session><session_id>%s</session_id><partner_id>%d</partner_id><create_dt>Sun Jun 28 02:55:27 PDT 2015</create_dt><media_server_url></media_server_url></Session></sessions>"

var sessionResponseBodyNoAuth = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><errorPayload><code>-1</code><message>No suitable authentication found</message></errorPayload>"

var sessionHelper *SessionHelper
var tokenHelper *TokenHelper

func init() {
	sessionHelper = &SessionHelper{}
	tokenHelper = &TokenHelper{}
}

// Session gives access to a SessionHelper instance
func Session() *SessionHelper {
	return sessionHelper
}

// Token gives access to a TokenHelper instance
func Token() *TokenHelper {
	return tokenHelper
}

// SessionHelper is an object helper to generate responses
// for the session resource
type SessionHelper struct {
}

// Request creates a new request to the session resource
// to create a session
func (s *SessionHelper) Request(props map[string]string) *Request {
	if props == nil {
		props = make(map[string]string)
	}

	var expectedProps = map[string]string{
		"location":       "",
		"p2p.preference": "disabled",
		"archiveMode":    "manual",
	}
	// We set the default values for the properties that have
	// not been set by the client
	for key, value := range expectedProps {
		if _, ok := props[key]; !ok {
			props[key] = value
		}
	}

	body := formURLEncode(props)
	url := fmt.Sprintf("%s/session/create", baseURL)
	return NewRequestWithBody("POST", url, body)
}

// ValidResponse generates a valid response for
// requests to the session resource to create a session
func (s *SessionHelper) ValidResponse(sessionID string, apiKey int) *Response {
	return NewResponseWithBody(200,
		fmt.Sprintf(sessionResponseBody, sessionID, apiKey))
}

// InvalidResponseNoAuth generates a response that is
// generated when the user tries to authenticate without
// X-TB-PARTNER-AUTH in the header
func (s *SessionHelper) InvalidResponseNoAuth() *Response {
	return NewResponseWithBody(403, sessionResponseBodyNoAuth)
}

// TokenHelper is a helper to decode a give token
type TokenHelper struct{}

// Decode decodes a token into a map so we can explore from
// the tests the parameters with which the token was created
func (t *TokenHelper) Decode(token string) (map[string]string, error) {
	if len(token) < 5 {
		return nil, errors.New("Token length is less than 5")
	}

	undecoded, err := decode64(token[4:])
	if err != nil {
		return nil, fmt.Errorf("Error decoding token: %s", err)
	}

	var apiKeyAndSignature = strings.Split(undecoded, ":")[0]
	var tokenParameters = strings.Split(undecoded, ":")[1]

	var apiKey = strings.Split(apiKeyAndSignature, "&")[0]

	var tokenParamsArray = strings.Split(tokenParameters, "&")
	var parameters map[string]string = map[string]string{
		strings.Split(apiKey, "=")[0]: strings.Split(apiKey, "=")[1],
	}

	for i := 0; i < len(tokenParamsArray); i++ {
		var param = tokenParamsArray[i]
		var key = strings.Split(param, "=")[0]
		var value = strings.Split(param, "=")[1]
		parameters[key] = value
	}

	return parameters, nil
}

func decode64(data string) (output string, err error) {
	var (
		decoder = base64.StdEncoding
		decoded []byte
	)

	for i := 0; i < 4; i++ {
		// We have to retry with different values
		decoded, err = decoder.DecodeString(data)
		if err == nil {
			break
		}
		data += "="
	}
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func formURLEncode(data map[string]string) string {
	var params = url.Values{}

	for key, value := range data {
		params.Add(key, value)
	}
	return params.Encode()
}
