package gonzo

import (
"net/http"
"github.com/gorilla/mux"
)

type Context struct {
	next ContextHandler
	r *http.Request
	paths map[string]string
}

func (p Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// shadow the request to do actions on
	p.r = r	
	p.paths = mux.Vars(r)
	
	// execute next step in the chain
	p.next(w, r, &p)
}

// The principal handler, when set
func (p *Context) Principal() string {
	return p.r.Header.Get("X-Principal")
}

// Get a named path in the Handler definition
func (p *Context) NamedPath(key string) string {
	return p.paths[key]
}