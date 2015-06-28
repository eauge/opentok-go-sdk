package helpers

import "fmt"

var archiveResponseBody = "{\"createdAt\" : 1384221730555,\n \"duration\" : 60,\n \"hasAudio\" : %t,\n \"hasVideo\" : %t,\n \"id\" : \"%s\",\n \"name\" : \"%s\",\n \"partnerId\" : %d,\n \"reason\" : \"\",\n \"sessionId\" : \"%s\",\n \"size\" : 0,\n \"status\" : \"%s\",\n \"url\" : null}"

var archiveListResponseBody = "{ \"count\" : %d, \"items\" : [ %s ] }"

// ArchiveParams can be used to set up the desired
// archive when formatting against archiveResponseBody
type ArchiveParams struct {
	HasAudio  bool
	HasVideo  bool
	ID        string
	Name      string
	APIKey    int
	SessionID string
	Status    string
}

var archiveHelper *ArchiveHelper

func init() {
	archiveHelper = &ArchiveHelper{}
}

// Archive gives access to an ArchiveHelper instance
func Archive() *ArchiveHelper {
	return archiveHelper
}

// ArchiveHelper is an object helper to generate responses
// for the session resource
type ArchiveHelper struct {
}

// DefaultParams returns a default ArchiveParams struct
func (a *ArchiveHelper) DefaultParams() *ArchiveParams {
	return &ArchiveParams{
		HasAudio:  true,
		HasVideo:  true,
		ID:        "archiveId",
		Name:      "archive",
		APIKey:    123456,
		SessionID: "sessionId",
		Status:    "started",
	}
}

// RequestStart generates a request for ArchiveStart
func (a *ArchiveHelper) RequestStart(apiKey int, sessionID string, props map[string]interface{}) *Request {
	if props == nil {
		props = make(map[string]interface{})
	}

	var expectedProps = map[string]interface{}{
		"hasAudio":   true,
		"hasVideo":   true,
		"name":       "",
		"outputMode": "composed",
	}
	// We set the default values for the properties that have
	// not been set by the client
	for key, value := range expectedProps {
		if _, ok := props[key]; !ok {
			props[key] = value
		}
	}

	props["sessionId"] = sessionID
	url := fmt.Sprintf("%s/v2/partner/%d/archive", baseURL, apiKey)
	return NewRequestWithBodyJSON("POST", url, props)
}

// RequestStop generates a request for ArchiveStop
func (a *ArchiveHelper) RequestStop(apiKey int, archiveID string) *Request {
	url := fmt.Sprintf("%s/v2/partner/%d/archive/%s/stop",
		baseURL, apiKey, archiveID)
	return NewRequest("POST", url)
}

// RequestGet generates a request for ArchiveGet
func (a *ArchiveHelper) RequestGet(apiKey int, archiveID string) *Request {
	url := fmt.Sprintf("%s/v2/partner/%d/archive/%s",
		baseURL, apiKey, archiveID)
	return NewRequest("GET", url)
}

// RequestList generates a request for ArchiveList
func (a *ArchiveHelper) RequestList(apiKey, count, offset int) *Request {
	url := fmt.Sprintf("%s/v2/partner/%d/archive?offset=%d&count=%d",
		baseURL, apiKey, offset, count)
	return NewRequest("GET", url)
}

// RequestDelete generates a request for ArchiveDelete
func (a *ArchiveHelper) RequestDelete(apiKey int, archiveID string) *Request {
	url := fmt.Sprintf("%s/v2/partner/%d/archive/%s",
		baseURL, apiKey, archiveID)
	return NewRequest("DELETE", url)
}

// ValidResponseWithArchive generates a response that
// will contain a JSON archive
func (a *ArchiveHelper) ValidResponseWithArchive(params *ArchiveParams) *Response {
	body := fmt.Sprintf(archiveResponseBody, params.HasAudio,
		params.HasVideo, params.ID, params.Name, params.APIKey,
		params.SessionID, params.Status)
	return NewResponseWithBody(200, body)
}

// ValidResponseWithArchiveList generates a response that
// will contain a JSON archive list
func (a *ArchiveHelper) ValidResponseWithArchiveList(count int) *Response {
	if count <= 0 {
		panic("Count cannot be less than 1")
	}

	params := a.DefaultParams()
	archive := fmt.Sprintf(archiveResponseBody, params.HasAudio,
		params.HasVideo, params.ID, params.Name, params.APIKey,
		params.SessionID, params.Status)

	archiveList := fmt.Sprintf("%s", archive)
	for i := 1; i < count; i++ {
		archiveList = fmt.Sprintf("%s,%s", archiveList, archive)
	}
	archiveList = fmt.Sprintf(archiveListResponseBody, count, archiveList)
	return NewResponseWithBody(200, archiveList)
}

// ValidResponseEmpty generates a 200 empty response
func (a *ArchiveHelper) ValidResponseEmpty() *Response {
	return NewResponse(200)
}

// InvalidResponseAuth generates a failed 403 response
func (a *ArchiveHelper) InvalidResponseAuth() *Response {
	return NewResponse(403)
}
