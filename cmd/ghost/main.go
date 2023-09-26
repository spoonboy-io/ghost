package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spoonboy-io/ghost/internal/mocks"
	"github.com/spoonboy-io/ghost/mocks/remedy"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MockLoaderResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
}

type MockErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Detail     string `json:"detail"`
}

var mocksCache = make(map[string]mocks.Mock)

func handler(w http.ResponseWriter, r *http.Request) {
	var mock mocks.Mock
	var ok bool

	// strip end point, and verb
	key := fmt.Sprintf("%s-%s", r.URL, r.Method)
	if mock, ok = mocksCache[key]; !ok {
		res := MockErrorResponse{}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res.StatusCode = http.StatusBadRequest
		res.Status = "Bad Request"
		res.Detail = fmt.Sprintf("No mock for found for Url:%s and Method: %s", r.URL, r.Method)
		out, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
		}
		_, _ = w.Write(out)
		return
	}

	// we have a mock we can respond with
	// check request meets expectations
	// request headers
	allHeaders := true

	for mk, mv := range mock.Request.Headers {
		hv := r.Header.Get(mk)
		if hv == "" {
			allHeaders = false
		}
		if hv != mv {
			allHeaders = false
		}
	}

	if !allHeaders {
		res := MockErrorResponse{}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)
		res.StatusCode = http.StatusNotAcceptable
		res.Status = "Not Acceptable"
		res.Detail = fmt.Sprintf("Request Headers do not meet expectations. Wanted: %v, Got: %v", mock.Request.Headers, r.Header)
		out, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
		}
		_, _ = w.Write(out)
		return
	}

	// request body
	allRequestBody := true
	reqBody := mocks.Properties{}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	if len(bytes) > 0 {
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			// make the map by parsing
			pairs := strings.Split(string(bytes), "&")
			for _, pair := range pairs {
				kv := strings.Split(pair, "=")
				reqBody[kv[0]] = kv[1]
			}
		} else {
			// json to unmarshal
			err = json.Unmarshal(bytes, &reqBody)
			if err != nil {
				log.Fatalln("could not unmarshal request body", err)
			}
		}
	}

	for mk, mv := range mock.Request.Body {
		var bv string
		var ok bool
		if bv, ok = reqBody[mk]; !ok {
			allRequestBody = false
		} else {
			if mv != bv {
				allRequestBody = false
			}
		}
	}

	if !allRequestBody {
		res := MockErrorResponse{}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)
		res.StatusCode = http.StatusNotAcceptable
		res.Status = "Not Acceptable"
		res.Detail = fmt.Sprintf("Request Body does not meet expecations. Wanted: %v, Got: %v", mock.Request.Body, reqBody)
		out, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
		}
		_, _ = w.Write(out)
		return
	}

	// if here we are good and we'll output the mock response
	w.WriteHeader(mock.Response.StatusCode)
	for k, v := range mock.Response.Headers {
		w.Header().Add(k, v)
	}

	// handle text/plain
	if ct, ok := mock.Response.Headers["Content-Type"]; ok {
		if ct == "text/plain" {
			w.Header().Set("content-type", "text/plain")
			_, _ = w.Write([]byte(mock.Response.Body))
			return
		}
	}

	body, err := json.Marshal(mock.Response.Body)
	if err != nil {
		log.Fatalln("could not marshal response body", err)
	}
	_, _ = w.Write(body)
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

			mock := mocks.Mock{}
			if err := json.Unmarshal(body, &mock); err != nil {
				onErr = true
				log.Println("could not unmarshal request body")
			}

			// add/update the mocks list
			// key is endpoint, verb
			mockKey := fmt.Sprintf("%s-%s", mock.EndPoint, mock.Request.Verb)
			mocksCache[mockKey] = mock
			log.Printf("added new mock '%s'\n", mockKey)
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
	// read port from cli -p flag or default to 9999
	var port int
	flag.IntVar(&port, "p", 9999, "Specify a port number (default is 9999")
	flag.Parse()
	portStr := fmt.Sprintf(":%d", port)

	// handlers
	// everything hits this endpoint
	http.HandleFunc("/", handler)
	// except this one, where we can load mock config in realtime
	http.HandleFunc("/load/mock", mockLoader)

	// as well as load mocks via the above server endpoint
	// we have the ability to include packaged mocks for things we may reuse
	packagedMocks := []mocks.Mocker{
		// add packaged mocks here, which must satisfy the mocks.Mocker interface
		remedy.Remedy{},
	}

	// add packaged mocks to mocksCache
	for _, pkg := range packagedMocks {
		pkgMocks := pkg.Mocks()
		log.Printf("loading mocks from '%s' package\n", pkg.Name())
		for _, mock := range pkgMocks {
			mockKey := fmt.Sprintf("%s-%s", mock.EndPoint, mock.Request.Verb)
			mocksCache[mockKey] = mock
		}
	}

	log.Println("starting Ghost server on port", portStr)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatalln("failed to start server")
	}
}
