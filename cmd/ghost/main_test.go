package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	Name           string
	Mock           Mock
	BreakBody      bool
	WantStatusCode int
	WantStatus     MockLoaderResponse
}

func TestMockLoader(t *testing.T) {
	mocks = map[string]Mock{}

	testCases := []TestCase{
		{
			Name: "Bad - sending GET request",
			Mock: Mock{
				EndPoint: "test/bad-get",
				Request: Request{
					Verb:    "GET",
					Headers: "{\"content-type\":\"application/json\"}",
				},
				Response: Response{
					StatusCode: 200,
					Headers:    "{\"content-type\":\"application/json\"}",
					Body:       "{\"hello\":\"world\"}",
				},
			},
			WantStatusCode: http.StatusMethodNotAllowed,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			Name: "Bad - sending POST request, bad body",
			Mock: Mock{
				EndPoint: "test/bad-post",
				Request: Request{
					Verb:    "POST",
					Headers: "{\"content-type\":\"application/json\"}",
				},
				Response: Response{
					StatusCode: 200,
					Headers:    "{\"content-type\":\"application/json\"}",
					Body:       "{\"hello\":\"world\"}",
				},
			},
			BreakBody:      true,
			WantStatusCode: http.StatusBadRequest,
			WantStatus: MockLoaderResponse{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			Name: "Good sending POST request with body",
			Mock: Mock{
				EndPoint: "test/good",
				Request: Request{
					Verb:    "POST",
					Headers: "{\"content-type\":\"application/json\"}",
					Body:    "{}",
				},
				Response: Response{
					StatusCode: 201,
					Headers:    "{\"content-type\":\"application/json\"}",
					Body:       "{\"hello\":\"world\"}",
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
			if tc.BreakBody {
				tReqBytes = tReqBytes[0:35]
			}
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(tc.Mock.Request.Verb, tc.Mock.EndPoint, bytes.NewReader(tReqBytes))
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
