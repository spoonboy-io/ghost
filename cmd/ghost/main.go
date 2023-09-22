package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Verb    string `json:"verb"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

type Response struct {
	StatusCode int    `json:"status"`
	Headers    string `json:"headers"`
	Body       string `json:"body"`
}

type Mock struct {
	EndPoint string   `json:"endPoint"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type MockLoaderResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
}

var mocks map[string]Mock

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}

func mockLoader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	res := MockLoaderResponse{}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		res.StatusCode = http.StatusMethodNotAllowed
		res.Status = "Method not allowed"
	} else {
		// parse the mock config from request body
		onErr := false
		if r.Body == nil {
			onErr = true
		} else {
			body, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			if err != nil {
				onErr = true
				log.Println("could not read request body")
			}

			mock := Mock{}
			if err := json.Unmarshal(body, &mock); err != nil {
				onErr = true
				log.Println("could not unmarshal request body")
			}

			// add/update the mocks list
			// key is endpoint, verb response statuscode
			mockKey := fmt.Sprintf("%s-%s-%d", mock.EndPoint, mock.Request.Verb, mock.Response.StatusCode)
			mocks[mockKey] = mock
		}

		if onErr {
			// bad request
			w.WriteHeader(http.StatusBadRequest)
			res.StatusCode = http.StatusBadRequest
			res.Status = "Bad request"
		}

		w.WriteHeader(http.StatusCreated)
		res.StatusCode = http.StatusCreated
		res.Status = "Created"
	}

	out, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}

	_, _ = w.Write(out)
}

func main() {
	mocks = map[string]Mock{}

	// port for server should be read from command line
	port := "9999"
	port = fmt.Sprintf(":%s", port)
	http.HandleFunc("/", handler)
	http.HandleFunc("/load/mock", mockLoader)

	log.Println("Starting Ghost server on port", port)
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatalln("failed to start server")
	}
}
