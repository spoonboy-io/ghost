package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/spoonboy-io/ghost/internal/mocks"
	"github.com/spoonboy-io/koan"
	"io/ioutil"
	"net/http"
	"strings"
)

// MockLoaderResponse is used to provide response data when requests are made to
// the server to load mocks to the cache
type MockLoaderResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
}

// MockErrorResponse is used to respond when request cannot be matched against aa cached mock
// or its header and body do not match the data in the cached mock. It is not used on succcess,
// on success the mocks status response headers and response body are the response
type MockErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Detail     string `json:"detail"`
}

type App struct {
	Logger *koan.Logger
}

// MocksCache is a map cache of mocks, store on key of `uri-method`
var MocksCache = make(map[string]mocks.Mock)

// Handler is the handler for all requests it parses the request to match
// against cached mocks (using endpoint and request method), if a match is found the incoming
// request header and request body is checked agains the data specified in the mock, if the a match
// the mock response is emitted to the client, otherwise errors are returned which identify how the request
// was not a match or the data supplied was unacceptable
func (a *App) Handler(w http.ResponseWriter, r *http.Request) {
	var mock mocks.Mock
	var ok bool

	msg := fmt.Sprintf("request '%s'", r.URL)
	a.Logger.Info(msg)

	// strip end point, and verb
	key := fmt.Sprintf("%s-%s", r.URL, r.Method)
	if mock, ok = MocksCache[key]; !ok {
		res := MockErrorResponse{}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res.StatusCode = http.StatusBadRequest
		res.Status = "Bad Request"
		res.Detail = fmt.Sprintf("No mock for found for Url:%s and Method: %s", r.URL, r.Method)
		out, err := json.Marshal(res)
		if err != nil {
			a.Logger.Error("problem marshaling response", err)
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
		res.Detail = fmt.Sprintf("Request Headers do not meet expectations. Wanted: %vGot: %v", mock.Request.Headers, r.Header)
		out, err := json.Marshal(res)
		if err != nil {
			a.Logger.Error("problem marshaling response", err)
		}
		_, _ = w.Write(out)
		return
	}

	// request body
	allRequestBody := true
	reqBody := mocks.Properties{}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.Logger.Error("problem reading request body", err)
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
				a.Logger.Error("problem marshaling request body", err)
			}
		}
	}

	if len(reqBody) == 0 && len(mock.Request.Body) != 0 {
		allRequestBody = false
	} else {
		for mk, mv := range mock.Request.Body {
			var bv string
			var ok bool
			// type assertion
			switch reqBody[mk].(type) {
			case string:
				//fmt.Println("body", mk)
				if bv, ok = reqBody[mk].(string); !ok {
					fmt.Println("missing", mk)
					allRequestBody = false
				} else {
					if mv != bv {
						allRequestBody = false
					}
				}
			case nil:
				fmt.Println("body", mk)
				if _, ok := reqBody[mk]; !ok {
					allRequestBody = false
				}
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
			a.Logger.Error("problem marshaling response", err)
		}
		_, _ = w.Write(out)
		return
	}

	// if here we are good and we'll output the mock response
	w.WriteHeader(mock.Response.StatusCode)
	for k, v := range mock.Response.Headers {
		w.Header().Add(k, v.(string))
	}

	// handle text/plain
	if ct, ok := mock.Response.Headers["Content-Type"]; ok {
		if ct == "text/plain" {
			w.Header().Set("content-type", "text/plain")
			// convert to json
			out := ""
			for k, v := range mock.Response.Body {
				if v == nil {
					out += fmt.Sprintf("%s", k)
				}
			}

			if err != nil {
				a.Logger.Error("could not marshal response body", err)
			}

			msg := fmt.Sprintf("response '%s'", out)
			a.Logger.Info(msg)

			_, _ = w.Write([]byte(out))
			return
		}
	}

	body, err := json.Marshal(mock.Response.Body)
	if err != nil {
		a.Logger.Error("could not marshal response body", err)
	}

	msg = fmt.Sprintf("response '%s'", string(body))
	a.Logger.Info(msg)

	_, _ = w.Write(body)
}

// MockLoader allows mocks signatures to be loaded to the server cache on the fly
// via a POST request to ths listening endpoint
func (a *App) MockLoader(w http.ResponseWriter, r *http.Request) {
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
				a.Logger.Error("problem reading request body", err)
			}

			mock := mocks.Mock{}
			if err := json.Unmarshal(body, &mock); err != nil {
				onErr = true
				a.Logger.Error("problem unmarshaling request body", err)
			}

			// add/update the mocks list
			// key is endpoint, verb
			mockKey := fmt.Sprintf("%s-%s", mock.EndPoint, mock.Request.Verb)
			MocksCache[mockKey] = mock
			a.Logger.Info(fmt.Sprintf("added new mock '%s'\n", mockKey))
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
		a.Logger.Error("problem marshaling response", err)
	}
	_, _ = w.Write(out)
}
