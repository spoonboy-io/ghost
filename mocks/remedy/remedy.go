/*
Package remedy provides mocks for Remedy server covering:
- authentication
- logout
- adding work order item
- getting status of work order item
*/
package remedy

import (
	"github.com/spoonboy-io/ghost/internal/mocks"
	"net/http"
)

// Remedy empty struct on which we implement the Mocker interface
type Remedy struct{}

// Mocks returns the mocks to be loaded as part of this package
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
				Body: mocks.Properties{
					"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ5WUJYMlBTdzlCeHRRcWZMZENpcTJuNmZUVDFkSmRHa2luZldFYW85SUNubEtxWjI4XC9FOXNxaW5FUkxVWDQ4Z3dlTnJKZit2VytZSEF6eUl0KzAwMkhPc2o4d3FYMTNOd2ZFMWxrY1hCZ0dvZ1E3Y2ZQSDJZdz09IiwibmJmIjoxNjk1MjE2ODExLCJpc3MiOiIxMjcuMC4wLjEiLCJleHAiOjE2OTUyMjA1MzEsIl9jYWNoZUlkIjozNTI4OSwiaWF0IjoxNjk1MjE2OTMxLCJqdGkiOiJJREdBNDYxN0pSVkpRQVNCMDNSSFNBMENDTFhQQ1UifQ.uWkZhPZ9eOjnaW9FMtDiXtWlcpmRODqv5MeRIdzGZPI": nil,
				},
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

		// create work order - good
		{
			EndPoint: "/api/arsys/v1/entry/SRM:RequestInterface_Create?fields=values(Request%20Number)",
			Request: mocks.Request{
				Verb: "POST",
				Headers: mocks.Properties{
					"Content-Type":  "application/json",
					"Authorization": "AR-JWT eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ5WUJYMlBTdzlCeHRRcWZMZENpcTJuNmZUVDFkSmRHa2luZldFYW85SUNubEtxWjI4XC9FOXNxaW5FUkxVWDQ4Z3dlTnJKZit2VytZSEF6eUl0KzAwMkhPc2o4d3FYMTNOd2ZFMWxrY1hCZ0dvZ1E3Y2ZQSDJZdz09IiwibmJmIjoxNjk1MjE2ODExLCJpc3MiOiIxMjcuMC4wLjEiLCJleHAiOjE2OTUyMjA1MzEsIl9jYWNoZUlkIjozNTI4OSwiaWF0IjoxNjk1MjE2OTMxLCJqdGkiOiJJREdBNDYxN0pSVkpRQVNCMDNSSFNBMENDTFhQQ1UifQ.uWkZhPZ9eOjnaW9FMtDiXtWlcpmRODqv5MeRIdzGZPI",
				},
				Body: mocks.Properties{
					"values": mocks.Properties{
						"TitleInstanceID":  "<%=instance.name%>",
						"Source Keyword":   "Morpheus CMP",
						"Login ID":         "<%=customOptions.vmowner%>",
						"OfferingTitle":    "Morpheus Provisioned Catalog Item",
						"SR Type Field 10": "Medium",
						"SR Type Field 14": "N/A",
						"SR Type Field 15": "Morpheus Catalog Item",
						"SR Type Field 17": "N/A",
						"SR Type Field 19": "<%=instance.instanceContext%>",
						"SR Type Field 20": "N/A",
						"SR Type Field 21": "<%=instance.hostname%>",
						"SR Type Field 22": "<%=instance.instanceTypeName%>",
						"SR Type Field 23": "<%=instance.description%>",
						"SR Type Field 28": "<%=instance.plan%>",
						"SR Type Field 29": "<%=customOptions.department%>",
						"SR Type Field 30": "<%=instance.backup.enabled%>",
						"SR Type Field 31": "Morpheus Instance Type",
						"SR Type Field 32": "<%=instance.instanceTypeName%>",
						"SR Type Field 33": "Morpheus Managed",
						"z1D Action":       "CREATE",
					},
				},
			},
			Response: mocks.Response{
				StatusCode: http.StatusCreated,
				Headers: mocks.Properties{
					"Content-Type": "application/json",
				},
				Body: mocks.Properties{
					"values": mocks.Properties{
						"Request Number": "1234",
					},
				},
			},
		},

		// get status of work order (waiting approval)
		{
			EndPoint: "/api/arsys/v1/entry/SRM:RequestApDetailSignature?fields=values(Request%20Number,Approval%20Status,Approvers)&q='Request%20Number'=%221234%22",
			Request: mocks.Request{
				Verb: "GET",
				Headers: mocks.Properties{
					"Content-Type":  "application/json",
					"Authorization": "AR-JWT eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ5WUJYMlBTdzlCeHRRcWZMZENpcTJuNmZUVDFkSmRHa2luZldFYW85SUNubEtxWjI4XC9FOXNxaW5FUkxVWDQ4Z3dlTnJKZit2VytZSEF6eUl0KzAwMkhPc2o4d3FYMTNOd2ZFMWxrY1hCZ0dvZ1E3Y2ZQSDJZdz09IiwibmJmIjoxNjk1MjE2ODExLCJpc3MiOiIxMjcuMC4wLjEiLCJleHAiOjE2OTUyMjA1MzEsIl9jYWNoZUlkIjozNTI4OSwiaWF0IjoxNjk1MjE2OTMxLCJqdGkiOiJJREdBNDYxN0pSVkpRQVNCMDNSSFNBMENDTFhQQ1UifQ.uWkZhPZ9eOjnaW9FMtDiXtWlcpmRODqv5MeRIdzGZPI",
				},
			},
			Response: mocks.Response{
				StatusCode: http.StatusOK,
				Headers: mocks.Properties{
					"Content-Type": "application/json",
				},
				Body: mocks.Properties{
					"entities": []mocks.Properties{{
						"values": mocks.Properties{
							"Request Number": "1234",
							//"Approval Status": "Pending",
							"Approval Status": "Approved",
						},
					},
					},
				},
			},
		},
	}
}

// Name returns the package name, which is displayed at start up
// when the packaged  mocks are loaded
func (Remedy) Name() string {
	return "Remedy"
}
