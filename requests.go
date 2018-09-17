package requests

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type Requester struct {
	method   string
	url      string
	body     interface{}
	response interface{}
}

func Request(method, url string) *Requester {
	return &Requester{
		method: method,
		url:    url,
	}
}

func (r *Requester) Body(body interface{}) *Requester {
	r.body = body
	return r
}

func (r *Requester) JSON(response interface{}) *Requester {
	r.response = response
	return r
}

func (r *Requester) Run() error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r.body)
	if err != nil {
		return errors.Wrap(err, "encoding body failed")
	}

	req, err := http.NewRequest(r.method, r.url, &buf)
	if err != nil {
		return errors.Wrap(err, "creating request failed")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if r.response != nil {
		err = json.NewDecoder(res.Body).Decode(r.response)
		if err != nil {
			return errors.Wrap(err, "decoding response into JSON failed")
		}
	}

	return nil
}

func Get(url string) *Requester {
	return Request(http.MethodGet, url)
}

func Post(url string) *Requester {
	return Request(http.MethodPost, url)
}

func Put(url string) *Requester {
	return Request(http.MethodPut, url)
}

func Delete(url string) *Requester {
	return Request(http.MethodDelete, url)
}
