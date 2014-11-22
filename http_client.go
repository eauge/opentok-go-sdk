package opentok

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type httpClient struct {
	apiKey    int
	apiSecret string
	apiUrl    string
	client    func() *http.Client
}

func newHttpClient(ot OpenTok) (c *httpClient) {
	c = new(httpClient)

	c.apiKey = ot.ApiKey
	c.apiSecret = ot.ApiSecret
	c.apiUrl = ot.apiUrl
	if len(ot.apiUrl) == 0 {
		c.apiUrl = "https://api.opentok.com"
	}

	if ot.AppEngine {
		c.client = clientWithTransport
	} else {
		c.client = clientWithoutTransport
	}
	return c
}

func (c *httpClient) Get(url string, headers map[string]string) (*http.Response, error) {
	return c.request("GET", c.apiUrl+"/"+url, headers, nil)
}

func (c *httpClient) Post(url string, headers, data map[string]string) (*http.Response, error) {
	return c.request("POST", c.apiUrl+"/"+url, headers, data)
}

func (c *httpClient) Delete(url string, headers map[string]string) error {
	var _, err = c.request("DELETE", c.apiUrl+"/"+url, headers, nil)
	return err
}

func (c *httpClient) request(method, url string, headers,
	data map[string]string) (res *http.Response, err error) {

	var (
		req    *http.Request
		client = c.client()
	)
	if req, err = c.createRequest(method, url, headers, data); err != nil {
		return nil, err
	}
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("Invalid response received from the server: %d", res.StatusCode))
	}

	return res, nil
}

func (c *httpClient) createRequest(method, url string, headers,
	data map[string]string) (r *http.Request, err error) {

	r, err = sendData(method, url, headers, data)
	if err != nil {
		return nil, err
	}

	// Adding headers common to all requests
	for key, value := range c.getCommonHeaders() {
		r.Header.Add(key, value)
	}

	// Adding headers specific to this request
	if headers != nil {
		for key, value := range headers {
			r.Header.Add(key, value)
		}
	}
	return r, nil
}

func sendData(method, url string, headers, data map[string]string) (r *http.Request, err error) {
	if data == nil {
		r, err = http.NewRequest(method, url, nil)
	} else {
		var dataString string
		dataString, err = processPostData(data, headers["Content-type"])
		if err != nil {
			return nil, err
		}
		var buffer = bytes.NewBufferString(dataString)
		r, err = http.NewRequest(method, url, bytes.NewReader(buffer.Bytes()))
	}
	return r, err
}

func processPostData(data map[string]string, contentType string) (string, error) {
	if contentType == "application/x-www-form-urlencoded" {
		return dataToQueryString(data), nil
	}
	return dataToJson(data)
}

func dataToJson(data map[string]string) (string, error) {
	var (
		dataBytes []byte
		err       error
	)

	if dataBytes, err = json.Marshal(data); err != nil {
		return "", err
	}
	return string(dataBytes), nil
}

func dataToQueryString(data map[string]string) string {
	var params = url.Values{}

	for key, value := range data {
		params.Add(key, value)
	}
	return params.Encode()
}

func (c *httpClient) getCommonHeaders() map[string]string {
	var partnerAuth = strconv.Itoa(c.apiKey) + ":" + c.apiSecret

	return map[string]string{
		"X-TB-PARTNER-AUTH": partnerAuth,
		"X-TB-VERSION":      "1",
	}
}

func clientWithoutTransport() *http.Client {
	return &http.Client{}
}

func clientWithTransport() *http.Client {
	transport := http.Transport{}
	return &http.Client{Transport: &transport}
}
