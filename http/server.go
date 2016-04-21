package gonzo

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"strconv"
	"net/http"
)

type MicroService struct {
	Health func(w http.ResponseWriter, r *http.Request, c *Context)
	muxx   *mux.Router
}

type ContextHandler func(http.ResponseWriter, *http.Request, *Context)


// Default health check. This method can be overridden before the Start method
// is called.
func ok(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintln(w, "ok")
}

// Instantiate a new microservice
func NewMicroService() *MicroService {
	return &MicroService{
		muxx:   mux.NewRouter(),
		Health: ok,
	}
}

 
// Wrap a Handler with AccessLogger and Principal
func (m *MicroService) Handle(method string, path string, handler ContextHandler) {
	fmt.Printf("Adding resource [%s] %s\n", method, path)
	m.muxx.Handle(path, Context {
	    next: AccessLogger{handler}.ServeHTTP,
	}).Methods(method)
}

// Wrap a Handler with AccessLogger and Principal
func (m *MicroService) Principal(method string, path string, handler ContextHandler) {
	fmt.Printf("Adding principal resource [%s] %s\n", method, path)
	m.muxx.Handle(path, Context{
	    next: AccessLogger{Principal{handler}.ServeHTTP}.ServeHTTP,
	}).Methods(method)
}

// Handle: Not Allowed Requests
func (ms *MicroService) NotAllowed(method string, path string) {
	fmt.Printf("NotAllowed resource [%s] %s\n", method, path)
	MethodNotAllowed := func(w http.ResponseWriter, r *http.Request, c *Context) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
 
	ms.muxx.Handle(path, Context{
	    next: AccessLogger{MethodNotAllowed}.ServeHTTP,
	}).Methods(method)
}


// Start a microservice with default health page on the given port
func (ms *MicroService) StartOnPort(port int) {

	// add health
	ms.Handle("GET", "/health", ms.Health)

	// start the web server
	fmt.Printf("Listening on %d....\n", port)
	
	if err := http.ListenAndServe(":" + strconv.Itoa(port), ms.muxx); err != nil {
		fmt.Println("error")
		log.Fatal("ListenAndServe:", err)
	} else {
		fmt.Println("running")
	}
}


// Start a microservice with default health page. It uses port 8080 by
// convention.
func (ms *MicroService) Start() {
 	ms.StartOnPort(8080)
}
