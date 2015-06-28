package opentok

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/eauge/opentok-go-sdk/helpers"
)

var apiKey int
var apiSecret string
var sessionID string
var archiveID string
var partnerAuth string

func TestMain(m *testing.M) {
	apiKey = 123456
	apiSecret = "API_SECRET"
	sessionID = "sessionId"
	archiveID = "archiveId"
	partnerAuth = fmt.Sprintf("%d:%s", apiKey, apiSecret)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	ot := New(apiKey, apiSecret)

	if ot.APIKey != apiKey {
		t.Fatalf("Different apikeys: expected: %d, in object: %d",
			apiKey, ot.APIKey)
	}
	if ot.APISecret != apiSecret {
		t.Fatalf("Different apisecrets: expected: %s, in object: %s",
			apiSecret, ot.APISecret)
	}
	if ot.client == nil {
		t.Fatalf("httpClient should not be nil")
	}
}

func TestNewPartnerAuth(t *testing.T) {
	ot := New(apiKey, apiSecret)
	partnerAuth := fmt.Sprintf("%d:%s", apiKey, apiSecret)
	if partnerAuth != ot.partnerAuth {
		t.Fatalf("Different partnerauths: expected: %s, in object: %s",
			partnerAuth, ot.partnerAuth)
	}
}

func TestNewWithAppEngine(t *testing.T) {
	ot := NewWithAppEngine(apiKey, apiSecret)

	if ot.APIKey != apiKey {
		t.Fatalf("Different apikeys: expected: %d, in object: %d",
			apiKey, ot.APIKey)
	}
	if ot.APISecret != apiSecret {
		t.Fatalf("Different apisecrets: expected: %s, in object: %s",
			apiSecret, ot.APISecret)
	}
	if ot.client == nil {
		t.Fatalf("httpClient should not be nil")
	}
}

func TestSession(t *testing.T) {
	req := helpers.Session().Request(make(map[string]string)).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Session().ValidResponse(sessionID, apiKey)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	session, err := ot.Session(&SessionProps{})

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if session == nil {
		t.Fatalf("Session should not be nil")
	}
	if session.ID != sessionID {
		t.Fatalf("Unexpected sessionId: expected: %s, received: %s",
			sessionID, session.ID)
	}
}

func TestSessionWithNil(t *testing.T) {
	req := helpers.Session().Request(make(map[string]string)).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Session().ValidResponse(sessionID, apiKey)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	session, err := ot.Session(nil)

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if session == nil {
		t.Fatalf("Session should not be nil")
	}
	if session.ID != sessionID {
		t.Fatalf("Unexpected sessionId: expected: %s, received: %s",
			sessionID, session.ID)
	}
}

func TestSessionWithLocation(t *testing.T) {
	req := helpers.Session().Request(map[string]string{
		"location": "127.0.0.1",
	}).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Session().ValidResponse(sessionID, apiKey)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	session, err := ot.Session(&SessionProps{Location: "127.0.0.1"})

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if session == nil {
		t.Fatalf("Session should not be nil")
	}
	if session.ID != sessionID {
		t.Fatalf("Unexpected sessionId: expected: %s, received: %s",
			sessionID, session.ID)
	}
}

func TestSessionWithMediaMode(t *testing.T) {
	req := helpers.Session().Request(map[string]string{
		"p2p.preference": "enabled",
	}).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Session().ValidResponse(sessionID, apiKey)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	session, err := ot.Session(&SessionProps{MediaMode: "enabled"})

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if session == nil {
		t.Fatalf("Session should not be nil")
	}
	if session.ID != sessionID {
		t.Fatalf("Unexpected sessionId: expected: %s, received: %s",
			sessionID, session.ID)
	}
}

func TestSessionWithArchiveMode(t *testing.T) {
	req := helpers.Session().Request(map[string]string{
		"archiveMode": "always",
	}).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Session().ValidResponse(sessionID, apiKey)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	session, err := ot.Session(&SessionProps{ArchiveMode: "always"})

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if session == nil {
		t.Fatalf("Session should not be nil")
	}
	if session.ID != sessionID {
		t.Fatalf("Unexpected sessionId: expected: %s, received: %s",
			sessionID, session.ID)
	}
}

func TestSessionWithAllProps(t *testing.T) {
	req := helpers.Session().Request(map[string]string{
		"archiveMode":    "always",
		"location":       "127.0.0.1",
		"p2p.preference": "enabled",
	}).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Session().ValidResponse(sessionID, apiKey)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	session, err := ot.Session(&SessionProps{
		ArchiveMode: "always",
		Location:    "127.0.0.1",
		MediaMode:   "enabled",
	})

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if session == nil {
		t.Fatalf("Session should not be nil")
	}
	if session.ID != sessionID {
		t.Fatalf("Unexpected sessionId: expected: %s, received: %s",
			sessionID, session.ID)
	}
}

func TestSessionInvalidAuth(t *testing.T) {
	req := helpers.Session().Request(make(map[string]string))
	res := helpers.Session().InvalidResponseNoAuth()
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	session, err := ot.Session(&SessionProps{})
	if session != nil {
		t.Fatalf("Session should be nil")
	}
	if err == nil {
		t.Fatalf("Expected err not to be nil: %s", err)
	}
	if match, _ := regexp.Match("403", []byte(err.Error())); !match {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestToken(t *testing.T) {
	ot := New(apiKey, apiSecret)
	session := &Session{ID: sessionID}

	token, err := ot.Token(session, &TokenProps{})
	decodedMap, _ := helpers.Token().Decode(token.String())

	if err != nil {
		t.Fatalf("Err should not be nil: %s", err)
	}
	if decodedMap["partner_id"] != strconv.Itoa(apiKey) {
		t.Fatalf("Invalid apiKey in token: %s, expected %d",
			decodedMap["partner_id"], apiKey)
	}
	if decodedMap["role"] != string(Publisher) {
		t.Fatalf("Invalid role in token: %s, expected %s",
			decodedMap["role"], Publisher)
	}
	if decodedMap["session_id"] != sessionID {
		t.Fatalf("Invalid role in token: %s, expected %s",
			decodedMap["session_id"], sessionID)
	}
}

func TestTokenWithNil(t *testing.T) {
	ot := New(apiKey, apiSecret)
	session := &Session{ID: sessionID}

	token, err := ot.Token(session, nil)
	decodedMap, _ := helpers.Token().Decode(token.String())

	if err != nil {
		t.Fatalf("Err should not be nil: %s", err)
	}
	if decodedMap["partner_id"] != strconv.Itoa(apiKey) {
		t.Fatalf("Invalid apiKey in token: %s, expected %d",
			decodedMap["partner_id"], apiKey)
	}
	if decodedMap["role"] != string(Publisher) {
		t.Fatalf("Invalid role in token: %s, expected %s",
			decodedMap["role"], Publisher)
	}
	if decodedMap["session_id"] != sessionID {
		t.Fatalf("Invalid role in token: %s, expected %s",
			decodedMap["session_id"], sessionID)
	}
}

func TestTokenWithParams(t *testing.T) {
	ot := New(apiKey, apiSecret)
	session := &Session{ID: sessionID}
	expireTime := time.Now().Unix() + 60
	data := "Some data"
	role := Subscriber

	token, err := ot.Token(session, &TokenProps{
		ExpireTime: expireTime,
		Data:       data,
		Role:       role,
	})
	decodedMap, _ := helpers.Token().Decode(token.String())

	if err != nil {
		t.Fatalf("Err should not be nil: %s", err)
	}
	if decodedMap["role"] != string(role) {
		t.Fatalf("Invalid role in token: %s, expected %s",
			decodedMap["role"], role)
	}
	if decodedMap["expire_time"] != strconv.FormatInt(expireTime, 10) {
		t.Fatalf("Invalid expireTime in token: %s, expected %d",
			decodedMap["expire_time"], expireTime)
	}
	if decodedMap["connection_data"] != data {
		t.Fatalf("Invalid connectionData in token: %s, expected %s",
			decodedMap["connection_data"], data)
	}
}

func TestTokenFails(t *testing.T) {
	ot := New(apiKey, apiSecret)
	session := &Session{}

	if _, err := ot.Token(session, &TokenProps{}); err == nil {
		t.Fatalf("Err should not be nil")
	}
}

func TestArchiveStart(t *testing.T) {
	req := helpers.Archive().RequestStart(apiKey, sessionID, nil).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	params := &helpers.ArchiveParams{
		HasAudio:  true,
		HasVideo:  true,
		ID:        archiveID,
		Name:      "archive",
		APIKey:    apiKey,
		SessionID: sessionID,
		Status:    "started",
	}
	res := helpers.Archive().ValidResponseWithArchive(params)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	archive, err := ot.ArchiveStart(sessionID, &ArchiveProps{
		SessionID: sessionID,
		HasAudio:  true,
		HasVideo:  true,
	})

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if archive == nil {
		t.Fatalf("Session should not be nil")
	}
	if archive.ID != archiveID {
		t.Fatalf("Unexpected archiveId: expected: %s, received: %s",
			archiveID, archive.ID)
	}
	if !archive.HasAudio {
		t.Fatalf("archive should have audio")
	}
	if !archive.HasVideo {
		t.Fatalf("archive should have video")
	}
	if archive.Name != params.Name {
		t.Fatalf("err: archive name: %s, expected: %s",
			archive.Name, params.Name)
	}
	if archive.Status != "started" {
		t.Fatalf("err: unexpected archive status: %s, expected: %s",
			archive.Status, "started")
	}
}

func TestArchiveStartIndividual(t *testing.T) {
	req := helpers.Archive().RequestStart(apiKey, sessionID, map[string]interface{}{
		"hasVideo":   false,
		"outputMode": "individual",
	}).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	params := &helpers.ArchiveParams{
		HasAudio:  true,
		HasVideo:  false,
		ID:        archiveID,
		Name:      "archive",
		APIKey:    apiKey,
		SessionID: sessionID,
		Status:    "started",
	}
	res := helpers.Archive().ValidResponseWithArchive(params)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	archive, err := ot.ArchiveStart(sessionID, &ArchiveProps{
		HasAudio:   true,
		HasVideo:   false,
		OutputMode: Individual,
	})

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if archive.HasVideo {
		t.Fatalf("archive should not have video")
	}
}

func TestArchiveStartFails(t *testing.T) {
	ot := newOpenTokWithClient(apiKey, apiSecret, helpers.NewClient())
	if _, err := ot.ArchiveStart("", nil); err == nil {
		t.Fatalf("Error should not be nil")
	}
}

func TestArchiveStop(t *testing.T) {
	req := helpers.Archive().RequestStop(apiKey, archiveID).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Archive().ValidResponseEmpty()
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	if err := ot.ArchiveStop(archiveID); err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
}

func TestArchiveStopFails(t *testing.T) {
	req := helpers.Archive().RequestStop(apiKey, archiveID).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Archive().ValidResponseEmpty()
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	if err := ot.ArchiveStop(""); err == nil {
		t.Fatalf("Expected err not to be nil")
	}
}

func TestArchiveDelete(t *testing.T) {
	req := helpers.Archive().RequestDelete(apiKey, archiveID).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Archive().ValidResponseEmpty()
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	if err := ot.ArchiveDelete(archiveID); err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
}

func TestArchiveDeleteFails(t *testing.T) {
	req := helpers.Archive().RequestDelete(apiKey, archiveID).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Archive().ValidResponseEmpty()
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	if err := ot.ArchiveDelete(""); err == nil {
		t.Fatalf("Expected err not to be nil")
	}
}

func TestArchiveGet(t *testing.T) {
	req := helpers.Archive().RequestGet(apiKey, archiveID).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	params := helpers.Archive().DefaultParams()
	res := helpers.Archive().ValidResponseWithArchive(params)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	archive, err := ot.ArchiveGet(archiveID)

	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if archive == nil {
		t.Fatalf("Session should not be nil")
	}
	if archive.ID != archiveID {
		t.Fatalf("Unexpected archiveId: expected: %s, received: %s",
			archiveID, archive.ID)
	}
}

func TestArchiveGetFails(t *testing.T) {
	req := helpers.Archive().RequestGet(apiKey, archiveID).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	params := helpers.Archive().DefaultParams()
	res := helpers.Archive().ValidResponseWithArchive(params)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	if _, err := ot.ArchiveGet(""); err == nil {
		t.Fatalf("Expected err not to be nil")
	}
}

func TestArchiveList(t *testing.T) {
	count := 10
	req := helpers.Archive().RequestList(apiKey, count, 0).
		AddHeader("X-TB-PARTNER-AUTH", partnerAuth)
	res := helpers.Archive().ValidResponseWithArchiveList(count)
	client := helpers.NewClient().
		Add(req, res)
	ot := newOpenTokWithClient(apiKey, apiSecret, client)

	archiveList, err := ot.ArchiveList(count, 0)
	if err != nil {
		t.Fatalf("Expected err to be nil: %s", err)
	}
	if archiveList.Count != count {
		t.Fatalf("Expected archivelist.Count to equal count: %d, actual: %d",
			count, archiveList.Count)
	}
	if len(archiveList.Archives) != count {
		t.Logf("archiveList: %v", archiveList)
		t.Fatalf("Expected archivelist to have length: %d, actual: %d",
			count, len(archiveList.Archives))
	}
}
