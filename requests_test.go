package greq

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	person = "person"
	book   = "book"
)

type Command struct {
	ObjectToFetch string
}

type Person struct {
	Name      string
	HairColor string
}

type Book struct {
	Title      string
	CopiesSold int
}

var mockServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Auth") == "secret-key" {
		fmt.Fprintln(w, `{"name": "Admin", "hairColor": "Gold"}`)
		return
	}
	var command Command
	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		fmt.Fprint(w, "error")
		return
	}
	switch command.ObjectToFetch {
	case person:
		fmt.Fprintln(w, `{"name": "James", "hairColor": "Brown"}`)
		return
	case book:
		fmt.Fprintln(w, `{"title": "Hitchhiker's Guide to the Galaxy", "copiesSold": 42}`)
		return
	}
	fmt.Fprintln(w, `{"name": "James", "hairColor": "Brown"}`)
})

func TestRequest(t *testing.T) {
	testServer := httptest.NewServer(mockServerHandler)
	defer testServer.Close()

	tests := []struct {
		name             string
		method           func(string, ...Configurer) (*Response, error)
		url              string
		body             interface{}
		bodyBytes        []byte
		headers          map[string]string
		response         interface{}
		expectedResponse interface{}
	}{
		{
			name:     "regular get",
			method:   Get,
			url:      testServer.URL,
			response: &Person{},
			expectedResponse: &Person{
				Name:      "James",
				HairColor: "Brown",
			},
		},
		{
			name:   "get with headers",
			method: Get,
			url:    testServer.URL,
			headers: map[string]string{
				"Auth": "secret-key",
			},
			response: &Person{},
			expectedResponse: &Person{
				Name:      "Admin",
				HairColor: "Gold",
			},
		},
		{
			name:   "post with body",
			method: Post,
			url:    testServer.URL,
			body: Command{
				ObjectToFetch: book,
			},
			response: &Book{},
			expectedResponse: &Book{
				Title:      "Hitchhiker's Guide to the Galaxy",
				CopiesSold: 42,
			},
		},
		{
			name:      "post with body bytes",
			method:    Post,
			url:       testServer.URL,
			bodyBytes: []byte(`{"objectToFetch": "book"}`),
			response:  &Book{},
			expectedResponse: &Book{
				Title:      "Hitchhiker's Guide to the Galaxy",
				CopiesSold: 42,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.method(
				test.url,
				Body(test.body),
				BodyBytes(test.bodyBytes),
				Headers(test.headers),
			)
			if err != nil {
				t.Fatal(err)
			}
			err = res.JSON(test.response)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.response, test.expectedResponse) {
				t.Fatalf("expected %#v to equal %#v", test.response, test.expectedResponse)
			}
		})
	}

}

func ExampleRequest() {
	res, err := Post("people-and-books.com", Body(Command{ObjectToFetch: book}))
	if err != nil {
		log.Fatalln(err)
	}
	var book Book
	res.JSON(&book)
}
