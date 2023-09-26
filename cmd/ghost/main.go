package main

import (
	"flag"
	"fmt"
	"github.com/spoonboy-io/ghost/internal/handlers"
	"github.com/spoonboy-io/ghost/internal/mocks"
	"github.com/spoonboy-io/ghost/mocks/remedy"
	"github.com/spoonboy-io/koan"
	"github.com/spoonboy-io/reprise"
	"log"
	"net/http"
)

var (
	version   = "Development build"
	goversion = "Unknown"
)

var logger *koan.Logger

func main() {
	// write a console banner
	reprise.WriteSimple(&reprise.Banner{
		Name:         "Ghost",
		Description:  "Mock Server for decoupled development against REST APIs",
		Version:      version,
		GoVersion:    goversion,
		WebsiteURL:   "https://spoonboy.io",
		VcsURL:       "https://github.com/spoonboy-io/ghost",
		VcsName:      "Github",
		EmailAddress: "hello@spoonboy.io",
	})

	// read port from cli -p flag or default to 9999
	var port int
	flag.IntVar(&port, "p", 9999, "Specify a port number (default is 9999")
	flag.Parse()
	portStr := fmt.Sprintf(":%d", port)

	// handlers
	// everything hits this endpoint
	http.HandleFunc("/", handlers.Handler)
	// except this one, where we can load mock config in realtime
	http.HandleFunc("/load/mock", handlers.MockLoader)

	// as well as load mocks via the above server endpoint
	// we have the ability to include packaged mocks for things we may reuse
	packagedMocks := []mocks.Mocker{
		// add new packaged mocks here, which must satisfy the mocks.Mocker interface
		remedy.Remedy{},
	}

	// add packaged mocks to mocksCache
	for _, pkg := range packagedMocks {
		pkgMocks := pkg.Mocks()
		log.Printf("loading mocks from '%s' package\n", pkg.Name())
		for _, mock := range pkgMocks {
			mockKey := fmt.Sprintf("%s-%s", mock.EndPoint, mock.Request.Verb)
			handlers.MocksCache[mockKey] = mock
		}
	}

	log.Println("starting Ghost server on port", portStr)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatalln("failed to start server")
	}
}
