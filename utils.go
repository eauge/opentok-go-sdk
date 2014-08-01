package opentok

import (
	"encoding/base64"
	"encoding/xml"
	"crypto/sha1"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
)

type XmlSessions struct {
	XMLName xml.Name `xml:"sessions"`
	Sessions []XmlSession `xml:"Session"`
}

type XmlSession struct {
	XMLName xml.Name `xml:"Session"`
	PartnerId int		`xml:"partner_id"`
	SessionId string	`xml:"session_id"`
	CreateDate string	`xml:"create_dt"`
}

func decodeArchive(body []byte) (archive Archive, err error) {
	err = json.Unmarshal(body, &archive)

	if err != nil {
		return Archive{}, err
	}
	return archive, nil
}

func decodeArchiveList(body []byte) (as []Archive, err error) {
	var archiveList archiveList

	err = json.Unmarshal(body, &archiveList)
	if err != nil {
		return nil, err
	}
	as = archiveList.Items
	return as, nil
}

func decodeSessionId(body []byte) (s string, err error) {
	var xmlSessions XmlSessions
	err = xml.Unmarshal(body, &xmlSessions)

	if err != nil {
		return "", err
	}

	return xmlSessions.Sessions[0].SessionId, nil
}

func encode64(data string) string {
	var encoder = base64.StdEncoding
	return encoder.EncodeToString([]byte(data))
}

func encodeHMAC(data string, key string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(data))
	mac := hash.Sum(nil)

	return hex.EncodeToString(mac)
}


// Used to decode a session id for testing purposes
func decode64(data string) (output string, err error) {
	var (
		decoder = base64.StdEncoding
		decoded []byte
	)

	for i := 0; i < 4; i ++ {
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
