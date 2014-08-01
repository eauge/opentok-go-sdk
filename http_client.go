package opentok

import (
	"net/http"
	"io"
	"io/ioutil"
	"bytes"
	"errors"
	"strconv"
	"fmt"
	"net/url"
	"encoding/json"
)

type httpClient struct {
	apiKey int
	apiSecret string
	apiUrl string
}

func newHttpClient (ot OpenTok) (c *httpClient) {
	c = new(httpClient)

	c.apiKey = ot.ApiKey
	c.apiSecret = ot.ApiSecret
	c.apiUrl= ot.apiUrl
	if len(ot.apiUrl) == 0 {
		c.apiUrl= "https://api.opentok.com"
	}

	return c
}

func (self *httpClient) get(url string, headers map[string]string) ([]byte, error) {
	return self.doRequest("GET", self.apiUrl + "/" + url, headers, nil)
}

func (self *httpClient) post(url string, headers, data map[string]string) ([]byte, error) {
	return self.doRequest("POST", self.apiUrl + "/" + url, headers, data)
}

func (self *httpClient) delete(url string, headers map[string]string) error {
	var _, err = self.doRequest("DELETE", self.apiUrl + "/" + url, headers, nil)
	return err
}

func (self *httpClient) doRequest(method, url string, headers,
	data map[string]string) (b []byte, err error) {

	var request *http.Request
	request, err = self.createRequest(method, url, headers, data)

	if err != nil {
		return nil, err
	}

	var client = &http.Client{}
	var response *http.Response

	response, err = client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Invalid response received from the server: %d", response.StatusCode))
	}

	return self.processResponse(response.Body)
}

func (self *httpClient) processResponse(body io.Reader) (b []byte, err error) {
	b, err = ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (self *httpClient) createRequest(method, url string, headers,
	data map[string]string) (r *http.Request, err error) {

	r, err = sendData(method, url, headers, data)
	if err != nil {
		return nil, err
	}

	// Adding headers common to all requests
	for key, value := range self.getCommonHeaders() {
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
	} else {
		return dataToJson(data)
	}
}

func dataToJson(data map[string]string) (string, error) {
	var dataBytes, err = json.Marshal(data)

	if err != nil {
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

func (self *httpClient) getCommonHeaders() map[string]string {
	var partnerAuth = strconv.Itoa(self.apiKey) + ":" + self.apiSecret

	return map[string]string {
		"X-TB-PARTNER-AUTH" : partnerAuth ,
		"X-TB-VERSION" : "1" ,
	}
}
