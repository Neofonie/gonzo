package gonzo

import "net/http"

type Context struct {
	next ContextHandler
	r *http.Request
}

func (p Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// shadow the request to do actions on
	p.r = r	
	
	// execute next step in the chain
	p.next(w, r, &p)
}

func (p *Context) Principal() string {
	return p.r.Header.Get("X-Principal")
}