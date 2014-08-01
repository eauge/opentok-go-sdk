package opentok

import (
	"fmt"
	"time"
	"strconv"
	"strings"
	"testing"
	"os"
	"errors"
)

var (
	apiKey    = readIntVariable("API_KEY", true)
	apiSecret = readStringVariable("API_SECRET", true)
	apiUrl 	= readStringVariable("API_URL", false)
)

func TestOpenTok(t *testing.T) {
	ot := OpenTok{ApiKey: apiKey, ApiSecret: apiSecret}
	if err := validate(ot); err != nil {
		t.Error(fmt.Sprintf("Opentok was not validated : %s", err))
		return
	}
}

func TestInvalidOpenTok(t *testing.T) {
	ot := OpenTok{}
	if err := validate(ot); err == nil {
		t.Error(fmt.Sprintf("Invalid Opentok has been validated as correct"))
		return
	}
}

func TestCreateSession(t *testing.T) {
	options := SessionOptions{}
	session, err := createSessionHelper(options)

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}

	if err = validateSession(session); err != nil {
		t.Error(fmt.Sprintf("Session id is not valid : %s", err))
		return
	}
}

func TestCreateRelayedSession(t *testing.T) {
	options := SessionOptions{MediaMode: Relayed}
	session, err := createSessionHelper(options)

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}

	if err = validateSession(session); err != nil {
		t.Error(fmt.Sprintf("Session id is not valid : %s", err))
		return
	}
}

func TestCreateRoutedSession(t *testing.T) {
	options := SessionOptions{MediaMode: Routed}
	session, err := createSessionHelper(options)

	if  err != nil{
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}
}

func TestCreateSessionWithLocation(t *testing.T) {
	options := SessionOptions{Location: "127.0.0.1"}
	session, err := createSessionHelper(options)

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}
}

func TestCreateRoutedSessionWithLocation(t *testing.T) {
	options := SessionOptions{MediaMode: Routed, Location: "127.0.0.1"}
	session, err := createSessionHelper(options)

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create();err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}
}

func TestGenerateTokenWithoutSession(t *testing.T) {
	options := SessionOptions{}
	session, err := createSessionHelper(options)

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if _, err = session.GenerateToken(TokenProperties{}); err == nil {
		t.Error("Token should not be generated if session has not been created")
		return
	}
}

func TestGenerateToken(t *testing.T) {
	session, err := createSessionHelper(SessionOptions{})

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}

	var (
		token Token
		tokenProperties = TokenProperties{}
	)
	if token, err = session.GenerateToken(tokenProperties); err != nil {
		t.Error(fmt.Sprintf("Token could not be generated : %s", err))
		return
	}

	var decodedToken map[string]string
	decodedToken, err = decodeToken(token)
	if len(decodedToken) == 0 {
		t.Error("Decoded token should have more than one field")
		return
	}

	var decodedPartnerId int
	decodedPartnerId, err = strconv.Atoi(decodedToken["partner_id"])
	if err != nil {
		t.Error("Partner id could not be decoded from token")
		return
	}

	if decodedToken["session_id"] != session.id || decodedToken["role"] != string(Publisher) ||
		decodedPartnerId != apiKey {
		t.Error("Parameters in token inconsistent with properties provided")
		return
	}
}

func TestGenerateSubscriberToken(t *testing.T) {
	session, err := createSessionHelper(SessionOptions{})

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}

	var (
		token Token
		tokenProperties = TokenProperties {Role: Subscriber}
	)
	if token, err = session.GenerateToken(tokenProperties); err != nil {
		t.Error(fmt.Sprintf("Token could not be generated : %s", err))
		return
	}

	var decodedToken map[string]string
	decodedToken, err = decodeToken(token)
	if len(decodedToken) == 0 {
		t.Error("Decoded token should have more than one field")
		return
	}

	var decodedPartnerId int
	decodedPartnerId, err = strconv.Atoi(decodedToken["partner_id"])
	if err != nil {
		t.Error("Partner id could not be decoded from token")
		return
	}

	if decodedToken["session_id"] != session.id || decodedToken["role"] != string(Subscriber) ||
		decodedPartnerId != apiKey {
		t.Error("Parameters in token inconsistent with properties provided")
		return
	}
}

func TestGenerateModeratorToken(t *testing.T) {
	session, err := createSessionHelper(SessionOptions{})

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}

	var (
		token Token
		tokenProperties = TokenProperties {Role: Moderator}
	)
	if token, err = session.GenerateToken(tokenProperties); err != nil {
		t.Error(fmt.Sprintf("Token could not be generated : %s", err))
		return
	}

	var decodedToken map[string]string
	decodedToken, err = decodeToken(token)
	if len(decodedToken) == 0 {
		t.Error("Decoded token should have more than one field")
		return
	}

	var decodedPartnerId int
	decodedPartnerId, err = strconv.Atoi(decodedToken["partner_id"])
	if err != nil {
		t.Error("Partner id could not be decoded from token")
		return
	}

	if decodedToken["session_id"] != session.id || decodedToken["role"] != string(Moderator) ||
		decodedPartnerId != apiKey {
		t.Error("Parameters in token inconsistent with properties provided")
		return
	}
}

func TestGenerateTokenWithExpiration(t *testing.T) {
	session, err := createSessionHelper(SessionOptions{})

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}

	var (
		token Token
		expireTime = time.Now().Unix() + 1000
		tokenProperties = TokenProperties {ExpireTime: expireTime}
	)
	if token, err = session.GenerateToken(tokenProperties); err != nil {
		t.Error(fmt.Sprintf("Token could not be generated : %s", err))
		return
	}

	var decodedToken map[string]string
	decodedToken, err = decodeToken(token)
	if len(decodedToken) == 0 {
		t.Error("Decoded token should have more than one field")
		return
	}

	var decodedPartnerId int
	decodedPartnerId, err = strconv.Atoi(decodedToken["partner_id"])
	if err != nil {
		t.Error("Partner id could not be decoded from token")
		return
	}

	var expireTimeNew int64
	expireTimeNew, err = strconv.ParseInt(decodedToken["expire_time"], 10, 64)
	if err != nil {
		t.Error("ExpireTime could not be decoded from token")
		return
	}

	if decodedToken["session_id"] != session.id || decodedToken["role"] != string(Publisher) ||
		decodedPartnerId != apiKey || expireTimeNew != expireTime {
		t.Error("Parameters in token inconsistent with properties provided")
		return
	}
}


func TestGenerateTokenWithData(t *testing.T) {
	session, err := createSessionHelper(SessionOptions{})

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}

	var (
		token Token
		data = "This is data for the token"
		tokenProperties = TokenProperties {Data: data}
	)
	if token, err = session.GenerateToken(tokenProperties); err != nil {
		t.Error(fmt.Sprintf("Token could not be generated : %s", err))
		return
	}

	var decodedToken map[string]string
	decodedToken, err = decodeToken(token)
	if len(decodedToken) == 0 {
		t.Error("Decoded token should have more than one field")
		return
	}

	var decodedPartnerId int
	decodedPartnerId, err = strconv.Atoi(decodedToken["partner_id"])
	if err != nil {
		t.Error("Partner id could not be decoded from token")
		return
	}

	if decodedToken["session_id"] != session.id || decodedToken["role"] != string(Publisher) ||
		decodedPartnerId != apiKey || data != decodedToken["connection_data"] {
		t.Error("Parameters in token inconsistent with properties provided")
		return
	}
}

func TestStartArchive(t *testing.T) {
	session, err := createSessionHelper(SessionOptions{})

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created : %s", err))
		return
	}

	if _, err = session.StartArchive("name"); err == nil {
		t.Error(fmt.Sprintf("StartArchive should fail without connections"))
		return
	}
}

func TestStopArchive(t *testing.T) {
	session, err := createSessionHelper(SessionOptions{})

	if  err != nil {
		t.Error(fmt.Sprintf("Session object should have been initialized : %s", err))
		return
	}

	if err = session.Create(); err != nil {
		t.Error(fmt.Sprintf("Session should have been created"))
		return
	}

	if _, err = session.StopArchive("ArchiveId"); err == nil {
		t.Error(fmt.Sprintf("StopArchive should fail without connections"))
		return
	}
}

func TestGetArchive(t *testing.T) {
	ot := OpenTok{ApiKey: apiKey, ApiSecret: apiSecret, apiUrl: apiUrl}
	if _, err := GetArchive(ot, "ArchiveId"); err == nil {
		t.Error(fmt.Sprintf("GetArchive should fail for an archive that does not exist"))
		return
	}
}

func TestDeleteArchive(t *testing.T) {
	ot := OpenTok{ApiKey: apiKey, ApiSecret: apiSecret, apiUrl: apiUrl}

	if err := DeleteArchive(ot, "ArchiveId"); err == nil {
		t.Error(fmt.Sprintf("Should fail to delete an archive that does not exist"))
		return
	}
}

func TestListArchives(t *testing.T) {
	ot := OpenTok{ApiKey: apiKey, ApiSecret: apiSecret, apiUrl: apiUrl}
	if archives, err := ListArchives(ot, 0, 100); err != nil || len(archives) < 0 {
		t.Error(fmt.Sprintf("Should return a list of archives"))
		return
	}
}

func createSessionHelper(options SessionOptions) (s *Session, err error) {
	var (
		opentok = OpenTok{ApiKey: apiKey, ApiSecret: apiSecret, apiUrl: apiUrl}
	)

	if s, err = NewSession(opentok, options); err != nil {
		return nil, err
	}
	return s, nil
}

func decodeToken(token Token) (map[string]string, error) {
	var tokenString = string(token)

	if len(tokenString) == 0 {
		return nil, errors.New("Token is an empty string")
	}
	undecodedToken, err := decode64(tokenString[4:])

	if err != nil {
		return nil, errors.New("Error decoding token")
	}

	var apiKeyAndSignature = strings.Split(undecodedToken, ":")[0]
	var tokenParameters = strings.Split(undecodedToken, ":")[1]

	var apiKey = strings.Split(apiKeyAndSignature, "&")[0]
	// var encodedSignature = strings.Split(apiKeyAndSignature, "&")[1]

	var tokenParamsArray = strings.Split(tokenParameters, "&")
	var parameters map[string]string = map[string]string {
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

func validateSession(session *Session) (err error) {

	if session == nil {
		return errors.New(fmt.Sprintf("Error when decoding session: nil session provided"))
	}
	sessionId := session.id
	if len(sessionId) < 2{
		return errors.New(fmt.Sprintf("Error when decoding session: sessionId is not long enough: %s", sessionId))
	}

	// remove sentinal (e.g. '1_', '2_')
	decodedSessionId := sessionId[2:]

	// replace invalid base64 chars
	decodedSessionId = strings.Replace(decodedSessionId, "-", "+", -1)
	decodedSessionId = strings.Replace(decodedSessionId, "_", "/", -1)
	decodedSessionId, err = decode64(decodedSessionId)
	if err != nil {
		return errors.New(fmt.Sprintf("Error when decoding session: %s", err))
	}

	fields := strings.Split(decodedSessionId, "~")
	var sessionApiKey int
	sessionApiKey, err = strconv.Atoi(fields[1])

	if err != nil {
		return errors.New(fmt.Sprintf("Error when decoding session: %s", err))
	}
	if sessionApiKey != apiKey {
		return errors.New("Api keys do not match")
	}
	return nil
}


func readIntVariable(variable string, mandatory bool) int {
	value := os.Getenv(variable)
	if len(value) == 0 && mandatory {
		panic(fmt.Sprintf("Environment variable : %s expected", variable))
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Environment variable : %s was expected to be an integer : %s", variable, err))
	}
	return intValue
}

func readStringVariable(variable string, mandatory bool) string {
	value := os.Getenv(variable)
	if len(value) == 0 && mandatory {
		panic(fmt.Sprintf("Environment variable : %s expected", variable))
	}

	return value
}
