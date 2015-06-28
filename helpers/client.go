package helpers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

var baseURL = "https://api.opentok.com"

type body struct {
	buf *bytes.Buffer
}

func (b *body) Read(p []byte) (int, error) {
	return b.buf.Read(p)
}

func (b *body) Close() error {
	return nil
}

func newBody(b string) *body {
	return &body{
		buf: bytes.NewBufferString(b),
	}
}

// NewRequest creates a new request object with
// an empty body
func NewRequest(method, url string) *Request {
	return NewRequestWithBody(method, url, "")
}

// NewRequestWithBody creates a new request object with
// the body provided
func NewRequestWithBody(method, url, b string) *Request {
	reader := newBody(b)
	req, _ := http.NewRequest(method, url, reader)
	req.ContentLength = int64(reader.buf.Len())

	return &Request{
		req: req,
	}
}

// NewRequestWithBodyXML creates a new request object
// and encodes the body interface{} to xml format
func NewRequestWithBodyXML(method, url string, body interface{}) *Request {
	buf := bytes.NewBufferString("")
	if err := xml.NewEncoder(buf).Encode(body); err != nil {
		panic(err)
	}
	return NewRequestWithBody(method, url, buf.String())
}

// NewRequestWithBodyJSON creates a new request object
// and encodes the body interface{} to json format
func NewRequestWithBodyJSON(method, url string, body interface{}) *Request {
	buf := bytes.NewBufferString("")
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		panic(err)
	}
	return NewRequestWithBody(method, url, buf.String())
}

// Request is a simple object to encapsulate the
// *http.Request interface
type Request struct {
	req *http.Request
}

// AddHeader adds a header to the httpRequest
func (r *Request) AddHeader(key, value string) *Request {
	r.req.Header.Add(key, value)
	return r
}

// NewResponse creates a new response object with
// an empty body
func NewResponse(statusCode int) *Response {
	return NewResponseWithBody(statusCode, "")
}

// NewResponseWithBody creates a new response object with
// the body provided
func NewResponseWithBody(statusCode int, body string) *Response {
	res := &http.Response{}
	res.StatusCode = statusCode
	res.Body = newBody(body)
	return &Response{
		res: res,
	}
}

// NewResponseWithBodyXML creates a new response object
// and encodes the body interface{} to xml format
func NewResponseWithBodyXML(statusCode int, body interface{}) *Response {
	buf := bytes.NewBufferString("")
	if err := xml.NewEncoder(buf).Encode(body); err != nil {
		panic(err)
	}
	return NewResponseWithBody(statusCode, buf.String())
}

// NewResponseWithBodyJSON creates a new request object
// and encodes the body interface{} to json format
func NewResponseWithBodyJSON(statusCode int, body interface{}) *Response {
	buf := bytes.NewBufferString("")
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		panic(err)
	}
	return NewResponseWithBody(statusCode, buf.String())
}

// Response is a simple object to encapsulate the
// *http.Response interface
type Response struct {
	res *http.Response
}

// Client is a simple object to mock up http responses
type Client struct {
	reqResMap       map[*Request]*Response
	defaultResponse *Response
}

// Add adds a new (request, response) pair that will be used
// when a request is performed to find the correct response
func (c *Client) Add(req *Request, res *Response) *Client {
	c.reqResMap[req] = res
	return c
}

// Do mocks an http request. It implements the httpClient interface
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res := c.findResponse(req)
	if res == nil {
		return nil, fmt.Errorf("Could not find request for req: %s",
			req.URL.String())
	}

	return res, nil
}

func (c *Client) findResponse(req *http.Request) *http.Response {
	for cReq, cRes := range c.reqResMap {
		if equalReq(cReq.req, req) {
			return cRes.res
		}
	}
	if c.defaultResponse != nil {
		return c.defaultResponse.res
	}
	return nil
}

func equalReq(cReq, req *http.Request) bool {
	// match urls
	if cReq.URL.String() != req.URL.String() {
		return false
	}
	// match headers
	for header := range cReq.Header {
		if req.Header.Get(header) != cReq.Header.Get(header) {
			return false
		}
	}

	// match body.
	if cReq.ContentLength != req.ContentLength {
		return false
	}

	if cReq.ContentLength == 0 {
		return true
	}

	bodyCReq := make([]byte, cReq.ContentLength)
	bodyReq := make([]byte, req.ContentLength)

	cReq.Body.Read(bodyCReq)
	req.Body.Read(bodyReq)

	cReq.Body.Close()
	req.Body.Close()

	// We need to restore the body because it will be
	// used in future comparisons
	cReq.Body = newBody(string(bodyCReq))
	if string(bodyCReq) != string(bodyReq) {
		return false
	}
	return true
}

// SetDefaultResponse sets an optional default response that is
// used when Client.Do does not find an appropriate response for
// the provided request
func (c *Client) SetDefaultResponse(res *Response) {
	c.defaultResponse = res
}

// NewClient creates a new client
func NewClient() *Client {
	return &Client{
		reqResMap: make(map[*Request]*Response),
	}
}
