package requests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type config struct {
	body interface{}
}

type Configurer func(*config)

func Body(body interface{}) Configurer {
	return func(config *config) {
		config.body = body
	}
}

type Response struct {
	responseBytes []byte
}

func Request(method, url string, configurers ...Configurer) (*Response, error) {
	var config config
	for _, configurer := range configurers {
		configurer(&config)
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(config.body)
	if err != nil {
		return nil, errors.Wrap(err, "encoding body failed")
	}

	req, err := http.NewRequest(method, url, &buf)
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

func (r *Response) JSON(response interface{}) error {
	if r.responseBytes != nil {
		err := json.NewDecoder(bytes.NewReader(r.responseBytes)).Decode(response)
		if err != nil {
			return errors.Wrap(err, "decoding response into JSON failed")
		}
	}
	return nil
}

func Get(url string, configurers ...Configurer) (*Response, error) {
	return Request(http.MethodGet, url, configurers...)
}

func Post(url string, configurers ...Configurer) (*Response, error) {
	return Request(http.MethodPost, url, configurers...)
}

func Put(url string, configurers ...Configurer) (*Response, error) {
	return Request(http.MethodPut, url, configurers...)
}

func Delete(url string, configurers ...Configurer) (*Response, error) {
	return Request(http.MethodDelete, url, configurers...)
}
