package greq

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type config struct {
	body      interface{}
	bodyBytes []byte
}

// Configurer will modify a Requests configuration.
type Configurer func(*config)

// Body is a Configurer which will set a body on the request.
func Body(body interface{}) Configurer {
	return func(config *config) {
		config.body = body
	}
}

// BodyBytes is a Configurer which will set a body on the request. It has priority over Body.
func BodyBytes(bytes []byte) Configurer {
	return func(config *config) {
		config.bodyBytes = bytes
	}
}

// Response contains information obtained after the request has been run.
type Response struct {
	responseBytes []byte
}

// Request will make an HTTP request. It can be configured by passing in optional Configurers. It returns a Response which can be further processed if need be.
func Request(method, url string, configurers ...Configurer) (*Response, error) {
	var config config
	for _, configurer := range configurers {
		configurer(&config)
	}

	buf := &bytes.Buffer{}
	if config.bodyBytes != nil {
		buf = bytes.NewBuffer(config.bodyBytes)
	} else {
		err := json.NewEncoder(buf).Encode(config.body)
		if err != nil {
			return nil, errors.Wrap(err, "encoding body failed")
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, errors.Wrap(err, "creating request failed")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading body failed")
	}

	return &Response{
		responseBytes: bs,
	}, nil
}

// JSON will try to decode the response into a Go object.
func (r *Response) JSON(response interface{}) error {
	if r.responseBytes != nil {
		err := json.NewDecoder(bytes.NewReader(r.responseBytes)).Decode(response)
		if err != nil {
			return errors.Wrap(err, "decoding response into JSON failed")
		}
	}
	return nil
}

// Get makes a GET Request
func Get(url string, configurers ...Configurer) (*Response, error) {
	return Request(http.MethodGet, url, configurers...)
}

// Post makes a POST Request
func Post(url string, configurers ...Configurer) (*Response, error) {
	return Request(http.MethodPost, url, configurers...)
}

// Put makes a PUT Request
func Put(url string, configurers ...Configurer) (*Response, error) {
	return Request(http.MethodPut, url, configurers...)
}

// Delete makes a DELETE Request
func Delete(url string, configurers ...Configurer) (*Response, error) {
	return Request(http.MethodDelete, url, configurers...)
}
