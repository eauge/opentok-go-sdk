package opentok

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"io"
)

type xmlSessions struct {
	XMLName  xml.Name     `xml:"sessions"`
	Sessions []xmlSession `xml:"Session"`
}

type xmlSession struct {
	XMLName    xml.Name `xml:"Session"`
	PartnerId  int      `xml:"partner_id"`
	SessionId  string   `xml:"session_id"`
	CreateDate string   `xml:"create_dt"`
}

func decodeArchive(body io.Reader) (archive *Archive, err error) {
	archive = new(Archive)
	if err = json.NewDecoder(body).Decode(archive); err != nil {
		return nil, err
	}
	return archive, nil
}

func decodeArchiveList(body io.Reader) ([]Archive, error) {
	var archiveList archiveList
	if err := json.NewDecoder(body).Decode(&archiveList); err != nil {
		return nil, err
	}
	return archiveList.Items, nil
}

func decodeSessionId(body io.Reader) (s string, err error) {
	var xmlSessions xmlSessions
	if err = xml.NewDecoder(body).Decode(&xmlSessions); err != nil {
		return "", err
	}
	return xmlSessions.Sessions[0].SessionId, nil
}

func encode64(data string) string {
	var encoder = base64.StdEncoding
	return encoder.EncodeToString([]byte(data))
}

func encodeHMAC(data []byte, key []byte) string {
	hash := hmac.New(sha1.New, key)
	hash.Write(data)
	mac := hash.Sum(nil)

	return hex.EncodeToString(mac)
}

// Used to decode a session id for testing purposes
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
