# Ghost

## A mock REST API server for developing client applications

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/spoonboy-io/ghost?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/spoonboy-io/ghost?style=flat-square)](https://goreportcard.com/report/github.com/spoonboy-io/ghost)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/spoonboy-io/ghost/build.yml?branch=main&style=flat-square)](https://github.com/spoonboy-io/ghost/actions/workflows/build.yml)
[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/spoonboy-io/ghost/unit_test.yml?branch=main&label=tests&style=flat-square)](https://github.com/spoonboy-io/ghost/actions/workflows/unit_test.yml)
[![GitHub Release Date](https://img.shields.io/github/release-date/spoonboy-io/ghost?style=flat-square)](https://github.com/spoonboy-io/ghost/releases)
[![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/spoonboy-io/ghost/latest?style=flat-square)](https://github.com/spoonboy-io/ghost/commits)
[![GitHub](https://img.shields.io/github/license/spoonboy-io/ghost?label=license&style=flat-square)](LICENSE)

## About

Ghost is a simple mock server for decoupled development of client API applications. 
Ghost is useful where an instance of the application hosting the REST API to be developed against is not available. 

Mocks can be loaded dynamically at runtime by making a POST request to the Ghost server, where they are cached until server shutdown.
Or, for mocks that are likely to be reused, these can be added as Go packages and built into the Ghost server application.

## Releases

You can find the [latest software here](https://github.com/spoonboy-io/ghost/releases/latest).

### Features

- POST mock request/response signature payloads to the running server
- Add packages of mocks for specific use cases which will likely be reused
- Call those mock API endpoints in your code/script and receive the registered responses
- Mocks are cached in memory

### Usage

#### Adding mocks at runtime

Load mocks to the server at runtime by making a POST request to the Ghost server on this endpoint

```shell
http://ghost/load/mock
```

The POST request body should be JSON which can be unmarshalled to a struct type of `mocks.Mock{}`

```go
// Mock represents a single mock, it's endpoint, the request, and the response
type Mock struct {
	EndPoint string   `json:"endPoint"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}
```
#### Creating mock packages to include at compile time

One package has already been created for Remedy and [can be found here](mocks/remedy/remedy.go). Use that as basis for creating
additional packages. The struct type of any Mocks package must satisfy the Mocker interface:

```go
// Mocker is simple interface to describe the values which can load a suite of mocks
// New packages can be created which implement this interface to preload mocks to the cache
// such that they do not need to be individually loaded to the server via POST request
type Mocker interface {
	Mocks() []Mock
	Name() string
}
```

Add packages you wish to include with your Ghost runtime in [main.go](cmd/ghost/main.go):

```go 
// as well as load mocks via the above server endpoint
// we have the ability to include packaged mocks for things we may reuse
packagedMocks := []mocks.Mocker{
	// add new packaged mocks here, which must satisfy the mocks.Mocker interface
	remedy.Remedy{},
}
```

### Installation

Clone the repository and use:

```
make build
```

Or grab the tar.gz or zip archive for your OS from the [releases page](https://github.com/spoonboy-io/ghost/releases/latest) and unpack it.

In either case, you can then start the server on the default port of `9999`:

```
./ghost
```

Or, specify the port:

```
./ghost -p 8008
```

To update the Ghost server, stop the server, replace the binary, then start the server.

### Development Opportunities

- Add mocks via configuration such as YAML

### License
Licensed under [Mozilla Public License 2.0](LICENSE)

