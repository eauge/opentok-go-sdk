package opentok

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Archive struct that holds all the information
// retrieved from the server
type Archive struct {

	// Unix timestamp that specified when the
	// archive was created
	CreatedAt int64

	// Duration of the archive in seconds
	Duration int64

	// Id of the archive. This is a unique id
	// identifier for the archive. It's used to
	// stop, retrieve and delete archives
	Id string

	// Name for the archive. The user can choose
	// any name but it doesn't necessarily need
	// to be different between archives
	Name string

	// The apiKey to which the archive belongs
	PartnerId int

	// The session id to which the archive belongs
	SessionId string

	// Size of the archives in KB
	Size int

	// Url from where the archive can be retrieved. This is
	// only useful if the archive is in status available
	// in the OpenTok S3 Account
	Url string

	// The Status of the Archive. The possibilities are:
	// - Started: if the archive is being recorded
	// - Stopped: if the archive has been stopped and it hasn't
	//   been uploaded or available
	// - Deleted: if the archive has been deleted. Only available
	//   archives can be deleted
	// - Uploaded: if the archive has been uploaded to the
	//   partner storage account
	// - Available: if the archive has been uploaded to the
	//   OpenTok S3 account
	// - Expired: available archives are removed from the OpenTok
	//   S3 account after 3 days. Their status become expired.
	Status string
}

// Struct that is used to decode the Json file
// received from the server. It should not be used by
// the user because the user will receive a real slice.
type archiveList struct {
	Count int
	Items []Archive
}

// Specifies how the streams will be managed
// by the OpenTok platform
type mediaMode int

// There are two possibilities for mediaMode
const (

	// The streams are handled by our routing component, Mantis.
	// This mode is useful for sessions with more than 3 participants
	// or sessions that need archiving
	Routed mediaMode = iota

	// The streams are either send P2P, if possible. If that's not
	// possible, they are routed to a TURN server, but streams
	// are not manager more efficiently and archiving is not possible.
	// Useful for sessions with up to three people that do not need
	// archiving
	Relayed
)

// The struct containing the different options when creating a
// new session object
type SessionOptions struct {
	MediaMode mediaMode

	// An IP address that enables the user to specifiy a preferred
	// location to be considered when creating the session in
	// the platform
	Location string
}

// Type used to specify the type of a session in an archive
type SessionRole string

// Array of options for the roles that a participant in an OpenTok
// session can have. Only one of them can be chosen for a token
const (

	// Has the highest level of control in the session. Can subscribe and
	// pubslish to a session and can also force the disconnections of other
	// participants in the session
	Moderator SessionRole = "moderator"

	// The standard role. It allows the participants to subscribe to
	// the streams sent by other participants as well as pubslish
	// its own stream to the session
	Publisher SessionRole = "publisher"

	// The role with the lowest control. It's only allowed to subscribe
	// to the streams in the session published by publishers and
	// moderators and cannot publish its stream
	Subscriber SessionRole = "subscriber"
)

// This struct stores the properties that can be used to generate a
// new token in the GenerateToken() method
type TokenProperties struct {
	Role       SessionRole
	ExpireTime int64
	Data       string
}

// A session Object
type Session struct {

	// The unique ID identifier for the session. It's a 64 bit encoded
	// string
	Id string

	// The OpenTok API_KEY that is generated when a new project
	// is created in the OpenTok Dashboard
	ApiKey int

	// The API_SECRET provided when a project is created in the
	// OpenTok Dashboard
	apiSecret string

	// The options that were used to create the session
	options SessionOptions

	// The http client used to interact with the OpenTok platform
	client *httpClient
}

// Creates a new session object using the information of an OpenTok project and
// the options. This function creates an Object, but it does not initialize
// the object with the platform. The method Create has the purpose of creating
// the session in the platform
func NewSession(ot OpenTok, options SessionOptions) (s *Session, err error) {
	if err = validate(ot); err != nil {
		return nil, err
	}

	session := new(Session)
	session.ApiKey = ot.ApiKey
	session.apiSecret = ot.ApiSecret
	session.options = options
	session.client = newHttpClient(ot)
	return session, nil
}

// Creates a session in the platform and stores the session ID got from the
// server to the Session Object
func (s *Session) Create() (err error) {
	var id string
	id, err = s.allocateSession()

	if err != nil {
		return errors.New(fmt.Sprintf("Error when creating session: %s", err))
	}
	s.Id = id
	return nil
}

// Generates a new token for the session so participants to the
// session can join
func (s *Session) GenerateToken(properties TokenProperties) (token Token, err error) {
	if len(s.Id) == 0 {
		return "", errors.New(fmt.Sprintf("Session has not been initialized yet. Please use Session.Create()"))
	}

	createTime := time.Now().Unix()
	nonce := rand.Int31() % 1000000

	role := properties.Role
	if len(properties.Role) == 0 {
		role = Publisher
	}

	var dataString string
	dataString, err = s.buildDataString(role, properties.ExpireTime,
		properties.Data, createTime, nonce)
	if err != nil {
		return "", err
	}
	return s.buildTokenString(dataString), nil
}

// Starts a new archive for the session. The archive id is generated by
// the OpenTok platform and the archive status becomes started
func (s *Session) StartArchive(name string) (archive Archive, err error) {
	if len(s.Id) == 0 {
		return Archive{}, errors.New(fmt.Sprintf("Session has not been initialized yet. Please use Session.Create()"))
	}

	var response []byte
	url := fmt.Sprintf("v2/partner/%d/archive", s.ApiKey)
	headers := map[string]string{"Content-type": "application/json"}
	data := map[string]string{"sessionId": s.Id, "name": name}
	response, err = s.client.post(url, headers, data)

	if err != nil {
		return Archive{}, err
	}

	archive, err = decodeArchive(response)

	if err != nil {
		return Archive{}, errors.New(fmt.Sprintf("Archive could not be decoded successfully"))
	}
	return archive, nil
}

// Stops an archive being recorded. If the archive is not in
// status started an error will be returned. The status of
// the archive becomes stopped
func (s *Session) StopArchive(archiveId string) (a Archive, err error) {
	if len(s.Id) == 0 {
		return Archive{}, errors.New(fmt.Sprintf("Session has not been initialized yet. Please use Session.Create()"))
	}

	if archiveId == "" {
		return Archive{}, errors.New("Archive id cannot be empty")
	}

	var response []byte
	url := fmt.Sprintf("v2/partner/%d/archive/%s/stop", s.ApiKey, archiveId)
	headers := map[string]string{"Content-type": "application/json"}
	response, err = s.client.post(url, headers, nil)

	if err != nil {
		return Archive{}, err
	}
	a, err = decodeArchive(response)

	if err != nil {
		return Archive{}, errors.New(fmt.Sprintf("Archive could not be decoded successfully"))
	}
	return a, nil
}

// Builds a token string using the data got from buildData string.
// It generates the final token that will be sent to the OpenTok
// client
func (s *Session) buildTokenString(dataString string) Token {
	signature := encodeHMAC(dataString, s.apiSecret)
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("partner_id=%d", s.ApiKey))
	buffer.WriteString(fmt.Sprintf("&sig=%s:%s", signature, dataString))

	return Token(fmt.Sprintf("T1==%s", encode64(buffer.String())))
}

// It builds the token string using the parameters provided by the user
func (s *Session) buildDataString(role SessionRole, expireTime int64, data string,
	createTime int64, nonce int32) (string, error) {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("session_id=%s", s.Id))
	buffer.WriteString(fmt.Sprintf("&create_time=%d", createTime))
	buffer.WriteString(fmt.Sprintf("&nonce=%d", nonce))
	buffer.WriteString(fmt.Sprintf("&role=%s", role))

	if expireTime > 0 {
		if isCorrectExpireTime(expireTime, createTime) {
			buffer.WriteString(fmt.Sprintf("&expire_time=%d", expireTime))
		} else {
			return "", errors.New("Invalid expire time set")
		}
	}

	if len(data) > 0 && len(data) < 1000 {
		buffer.WriteString(fmt.Sprintf("&connection_data=%s", data))
	}

	return buffer.String(), nil
}

// Creates a new session in the server using the OpenTok REST API
func (s *Session) createSession() (response []byte, err error) {
	var (
		p2p                       = "disabled"
		headers map[string]string = map[string]string{
			"Content-type": "application/x-www-form-urlencoded",
		}
	)

	if s.options.MediaMode == Relayed {
		p2p = "enabled"
	}
	var data map[string]string = map[string]string{
		"location":       s.options.Location,
		"p2p.preference": p2p,
	}

	response, err = s.client.post("session/create", headers, data)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Allocates a session by calling create new session and reads the
// response to return the id of the new session
func (s *Session) allocateSession() (id string, err error) {
	var response []byte
	response, err = s.createSession()
	if err != nil {
		return "", err
	}

	id, err = readAllocateResponse(response)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Reads an allocation response and returns the id of the
// allocated session
func readAllocateResponse(response []byte) (id string, err error) {
	id, err = decodeSessionId(response)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Checks if the expiration time set is less than 30 days but
// bigger than the creation time for the token
func isCorrectExpireTime(expireTime, createTime int64) bool {

	// Less than 30 days
	return expireTime > createTime && expireTime < createTime+2592000
}
