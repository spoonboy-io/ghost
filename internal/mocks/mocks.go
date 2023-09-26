package mocks

type Properties map[string]string

type Request struct {
	Verb    string     `json:"verb"`
	Headers Properties `json:"headers"`
	Body    Properties `json:"body"`
}

type Response struct {
	StatusCode int        `json:"status"`
	Headers    Properties `json:"headers"`
	Body       string     `json:"body"`
}

type Mock struct {
	EndPoint string   `json:"endPoint"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}
