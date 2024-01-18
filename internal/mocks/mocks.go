package mocks

// Properties is used for to store key value attributes of headers and request bodies
// type Properties map[string]string
type Properties map[string]interface{}

// Request describes the data we keep about a mock request
type Request struct {
	Verb    string     `json:"verb"`
	Headers Properties `json:"headers"`
	Body    Properties `json:"body"`
}

// Response describes the data we store about a mock response
type Response struct {
	StatusCode int        `json:"status"`
	Headers    Properties `json:"headers"`
	Body       Properties `json:"body"`
}

// Mock represents a single mock, it's endpoint, the request, and the response
type Mock struct {
	EndPoint string   `json:"endPoint"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

// Mocker is simple interface to describe the values which can load a suite of mocks
// New packages can be created which implement this interface to preload mocks to the cache
// such that they do not need to be individually loaded to the server via POST request
type Mocker interface {
	Mocks() []Mock
	Name() string
}
