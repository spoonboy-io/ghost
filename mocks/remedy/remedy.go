/*
Package remedy provides mocks for Remedy backup server covering:
- authentication
- logout
- TODO adding work order item
- TODO getting status of work order item
*/
package remedy

import (
	"github.com/spoonboy-io/ghost/internal/mocks"
	"net/http"
)

type Remedy struct{}

func (Remedy) Mocks() []mocks.Mock {
	return []mocks.Mock{
		// authentication request
		{
			EndPoint: "/api/jwt/login",
			Request: mocks.Request{
				Verb: "POST",
				Headers: mocks.Properties{
					"Content-Type": "application/x-www-form-urlencoded",
				},
				Body: mocks.Properties{
					"username":   "admin",
					"password":   "password",
					"authString": "",
				},
			},
			Response: mocks.Response{
				StatusCode: http.StatusOK,
				Headers: mocks.Properties{
					"Content-Type": "text/plain",
				},
				Body: "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ5WUJYMlBTdzlCeHRRcWZMZENpcTJuNmZUVDFkSmRHa2luZldFYW85SUNubEtxWjI4XC9FOXNxaW5FUkxVWDQ4Z3dlTnJKZit2VytZSEF6eUl0KzAwMkhPc2o4d3FYMTNOd2ZFMWxrY1hCZ0dvZ1E3Y2ZQSDJZdz09IiwibmJmIjoxNjk1MjE2ODExLCJpc3MiOiIxMjcuMC4wLjEiLCJleHAiOjE2OTUyMjA1MzEsIl9jYWNoZUlkIjozNTI4OSwiaWF0IjoxNjk1MjE2OTMxLCJqdGkiOiJJREdBNDYxN0pSVkpRQVNCMDNSSFNBMENDTFhQQ1UifQ.uWkZhPZ9eOjnaW9FMtDiXtWlcpmRODqv5MeRIdzGZPI",
			},
		},

		// logout
		{
			EndPoint: "/api/jwt/logout",
			Request: mocks.Request{
				Verb: "POST",
				Headers: mocks.Properties{
					"Authorization": "AR-JWT eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ5WUJYMlBTdzlCeHRRcWZMZENpcTJuNmZUVDFkSmRHa2luZldFYW85SUNubEtxWjI4XC9FOXNxaW5FUkxVWDQ4Z3dlTnJKZit2VytZSEF6eUl0KzAwMkhPc2o4d3FYMTNOd2ZFMWxrY1hCZ0dvZ1E3Y2ZQSDJZdz09IiwibmJmIjoxNjk1MjE2ODExLCJpc3MiOiIxMjcuMC4wLjEiLCJleHAiOjE2OTUyMjA1MzEsIl9jYWNoZUlkIjozNTI4OSwiaWF0IjoxNjk1MjE2OTMxLCJqdGkiOiJJREdBNDYxN0pSVkpRQVNCMDNSSFNBMENDTFhQQ1UifQ.uWkZhPZ9eOjnaW9FMtDiXtWlcpmRODqv5MeRIdzGZPI",
				},
			},
			Response: mocks.Response{
				StatusCode: http.StatusNoContent,
				Headers:    mocks.Properties{},
			},
		},

		// create work order

		// get status of work order
	}
}

func (Remedy) Name() string {
	return "Remedy"
}
