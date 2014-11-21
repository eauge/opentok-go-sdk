package opentok

import (
	"bytes"
	"errors"
	"fmt"
)

// Struct that contains all the necessary information to
// create sessions and interact with the OpenTok platform
type OpenTok struct {

	// This is the ApiKey that you get after creating
	// a project at the OpenTok Dashboard
	ApiKey int

	// This is the ApiSecret that you get after
	// creating a project with the OpenTok Dashboard
	ApiSecret string

	// This is just used internally to test the sdk
	// for the different environments
	apiUrl string
}

// We create this new type token for the tokens created
// by a session.
type Token struct {
	prop TokenProperties
}

func newToken(prop TokenProperties) *Token {
	return &Token{prop: prop}
}

func (t Token) Value() string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("partner_id=%d", t.prop.apiKey))
	buffer.WriteString(fmt.Sprintf("&sig=%s:%s", t.signature(), t.key()))

	return "T1==" + encode64(buffer.String())
}

func (t Token) Expires() int64 {
	switch {
	case t.prop.createTime == 0:
		return 0
	case t.prop.ExpireTime == 0:
		// a day
		return t.prop.createTime + 60*60*24
	case t.prop.ExpireTime < 0:
		return 0
	}
	return t.prop.createTime + int64(t.prop.ExpireTime)
}

func (t Token) key() []byte {
	buffer := bytes.NewBuffer([]byte(""))
	buffer.WriteString(fmt.Sprintf("session_id=%s", t.prop.sessionId))
	buffer.WriteString(fmt.Sprintf("&create_time=%d", t.prop.createTime))
	buffer.WriteString(fmt.Sprintf("&nonce=%d", t.prop.nonce))
	buffer.WriteString(fmt.Sprintf("&role=%s", t.prop.Role.get()))
	buffer.WriteString(fmt.Sprintf("&expire_time=%d", t.Expires()))

	if len(t.prop.Data) > 0 && len(t.prop.Data) < 1000 {
		buffer.WriteString(fmt.Sprintf("&connection_data=%s", t.prop.Data))
	}
	return buffer.Bytes()
}

func (t Token) signature() string {
	return encodeHMAC(t.key(), []byte(t.prop.apiSecret))
}

// Validates that the necessary parameters in the OpenTok
// structure have been set
func validate(ot OpenTok) error {
	if ot.ApiKey == 0 {
		return errors.New("ApiKey is not set")
	}

	if len(ot.ApiSecret) == 0 {
		return errors.New("ApiSecret is not set")
	}

	return nil
}

// Retrieves an archive from the server. If the
// archive does not exist an error will be raised
func GetArchive(ot OpenTok, archiveId string) (archive Archive, err error) {
	if archiveId == "" {
		return Archive{}, errors.New("Archive id cannot be empty")
	}

	client := newHttpClient(ot)
	url := fmt.Sprintf("v2/partner/%d/archive/%s", ot.ApiKey, archiveId)
	response, err := client.get(url, nil)

	if err != nil {
		return Archive{}, err
	}

	return decodeArchive(response)
}

// Deletes an existing archive with status available. If
// the archive is in any other state the operation will
// fail and return an error
func DeleteArchive(ot OpenTok, archiveId string) error {
	if archiveId == "" {
		return errors.New("Archive id cannot be empty")
	}
	client := newHttpClient(ot)
	url := fmt.Sprintf("v2/partner/%d/archive/%s", ot.ApiKey, archiveId)
	return client.delete(url, nil)
}

// Returns a list of archives. If Count == 0, the limit of
// the number of archives returned by the server is limited
// by the server. Otherwise it will be count. Offset is
// useful for pagination
func ListArchives(ot OpenTok, count, offset int) ([]Archive, error) {
	if count < 0 {
		return nil, errors.New("count cannot be smaller than 1")
	}

	client := newHttpClient(ot)
	url := fmt.Sprintf("v2/partner/%d/archive?offset=%d", ot.ApiKey, offset)
	if count > 0 {
		url = fmt.Sprintf("%s&count=%d", url, count)
	}

	response, err := client.get(url, nil)
	if err != nil {
		return nil, err
	}

	return decodeArchiveList(response)
}
