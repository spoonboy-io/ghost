package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	Name           string
	Method         string
	Mock           Mock
	BreakBody      bool
	WantStatusCode int
	WantStatus     MockLoaderResponse
}

func TestMockLoader(t *testing.T) {
	//mocks = map[string]Mock{}

	testCases := []TestCase{
		{
			Name:   "Bad - sending GET request",
			Method: "GET",
			Mock: Mock{
				EndPoint: "test/bad-get",
				Request: Request{
					Verb: "GET",
					Headers: Properties{
						"content-type": "application/json",
					},
				},
				Response: Response{
					StatusCode: 200,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusMethodNotAllowed,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			Name:   "Bad - sending POST request, bad body",
			Method: "POST",
			Mock: Mock{
				EndPoint: "test/bad-post",
				Request: Request{
					Verb: "POST",
					Headers: Properties{
						"content-type": "application/json",
					},
				},
				Response: Response{
					StatusCode: 200,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			BreakBody:      true,
			WantStatusCode: http.StatusBadRequest,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			Name:   "Good sending POST request with body",
			Method: "POST",
			Mock: Mock{
				EndPoint: "test/good",
				Request: Request{
					Verb: "POST",
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: Properties{},
				},
				Response: Response{
					StatusCode: 201,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusCreated,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusCreated,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tReqBytes, err := json.Marshal(tc.Mock)
			if err != nil {
				t.Fatal(err)
			}

			if tc.BreakBody {
				tReqBytes = tReqBytes[0:35]
			}

			req, err := http.NewRequest(tc.Method, "/load/mock", bytes.NewReader(tReqBytes))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(mockLoader)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.WantStatus.StatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.WantStatus.StatusCode)
			}

		})

		fmt.Println(mocks)
	}
}

func seedMockData() {
	// seed some dummy data for tests
	dummies := []Mock{
		{
			EndPoint: "good/data",
			Request: Request{
				Verb: "POST",
				Headers: Properties{
					"content-type": "application/json",
				},
				Body: Properties{
					"testKey": "test value",
				},
			},
			Response: Response{
				StatusCode: 201,
				Headers: Properties{
					"content-type": "application/json",
				},
				Body: "{\"hello\":\"world\"}",
			},
		},
	}

	// load the dummies to mocks map
	for _, mock := range dummies {
		key := fmt.Sprintf("%s-%s", mock.EndPoint, mock.Request.Verb)
		mocks[key] = mock
	}

}

func TestHandler(t *testing.T) {

	seedMockData()

	testCases := []TestCase{
		{
			Name:   "Good, request uri, method, headers and body match stored mock",
			Method: "POST",
			Mock: Mock{
				EndPoint: "good/data",
				Request: Request{
					Verb: "POST",
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: Properties{
						"testKey": "test value",
					},
				},
				Response: Response{
					StatusCode: 201,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusCreated,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusCreated,
			},
		},
		{
			Name:   "Bad, request method not match",
			Method: "POST",
			Mock: Mock{
				EndPoint: "good/data",
				Request: Request{
					Verb: "GET",
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: Properties{
						"testKey": "test value",
					},
				},
				Response: Response{
					StatusCode: 201,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusBadRequest,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			Name:   "Bad, request uri not match",
			Method: "POST",
			Mock: Mock{
				EndPoint: "good/data/",
				Request: Request{
					Verb: "GET",
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: Properties{
						"testKey": "test value",
					},
				},
				Response: Response{
					StatusCode: 201,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusBadRequest,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			Name:   "Bad request header not match",
			Method: "POST",
			Mock: Mock{
				EndPoint: "good/data",
				Request: Request{
					Verb: "POST",
					Headers: Properties{
						"content-type": "application/text",
					},
					Body: Properties{
						"testKey": "test value",
					},
				},
				Response: Response{
					StatusCode: 201,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusNotAcceptable,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusNotAcceptable,
			},
		},
		{
			Name:   "Bad, request header not present",
			Method: "POST",
			Mock: Mock{
				EndPoint: "good/data",
				Request: Request{
					Verb: "POST",
					Headers: Properties{
						"authorization": "bearer xxxxxx",
					},
					Body: Properties{
						"testKey": "test value",
					},
				},
				Response: Response{
					StatusCode: 201,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusNotAcceptable,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusNotAcceptable,
			},
		},
		{
			Name:   "Bad, request body property not match",
			Method: "POST",
			Mock: Mock{
				EndPoint: "good/data",
				Request: Request{
					Verb: "POST",
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: Properties{
						"testKey": "this is the wrong value",
					},
				},
				Response: Response{
					StatusCode: 201,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusNotAcceptable,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusNotAcceptable,
			},
		},
		{
			Name:   "Bad, request body propery not present",
			Method: "POST",
			Mock: Mock{
				EndPoint: "good/data",
				Request: Request{
					Verb: "POST",
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: Properties{},
				},
				Response: Response{
					StatusCode: 201,
					Headers: Properties{
						"content-type": "application/json",
					},
					Body: "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusNotAcceptable,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusNotAcceptable,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tReqBytes, err := json.Marshal(tc.Mock.Request.Body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(tc.Mock.Request.Verb, tc.Mock.EndPoint, bytes.NewReader(tReqBytes))
			if err != nil {
				t.Fatal(err)
			}

			// add the headers
			for k, v := range tc.Mock.Request.Headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.WantStatus.StatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.WantStatus.StatusCode)

				detailBytes, err := ioutil.ReadAll(rr.Body)
				if err != nil {
					t.Fatal(err)
				}
				t.Log(string(detailBytes))
			}

		})
	}
}
