package runtime

import (
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
)

// Mux is a HTTP Server mux for go-kit based services.
type Mux struct {
	endpoints []endpoint
}

type endpoint struct {
	*kithttp.Server

	method       string
	pathSegments []string
}

// NewMux returns a new mux with a blank state.
func NewMux() *Mux {
	return &Mux{}
}

// AddEndpoint adds the specified endpoint to the Mux.
func (mux *Mux) AddEndpoint(method, pathSegments string, ep *kithttp.Server) {
	mux.endpoints = append(mux.endpoints, endpoint{
		pathSegments: strings.Split(pathSegments, "/"),
		method:       method,
		Server:       ep,
	})
}

func (mux *Mux) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	pathParts := strings.Split(req.URL.Path, "/")
	invalidMethod := false

	if req.Method == "OPTIONS" {
		wr.Header().Set("Access-Control-Allow-Headers", "Authorization")
		wr.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		wr.WriteHeader(200)

		return
	}

	for _, endp := range mux.endpoints {
		if len(pathParts) == len(endp.pathSegments) && matchPath(endp.pathSegments, pathParts) {
			if endp.method != req.Method {
				invalidMethod = true
				continue
			}

			endp.ServeHTTP(wr, req)

			return
		}
	}

	if invalidMethod {
		http.Error(wr, "Method not allowed", 405)
		return
	}

	http.NotFound(wr, req)
}

func matchPath(endp, pathParts []string) bool {
	for i := 0; i < len(pathParts); i++ {
		endpPath := endp[i]
		partsPath := pathParts[i]
		if len(endpPath) > 0 {
			if endpPath[0] == '{' && endpPath[len(endpPath)-1] == '}' {
				continue
			}
		}

		if endpPath != partsPath {
			return false
		}
	}

	return true
}
