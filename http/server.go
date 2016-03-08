package gonzo

import (
	"fmt"
	"log"
	"net/http"	
	"github.com/gorilla/mux"
)

type Handler func(http.ResponseWriter, *http.Request);

type PrincipalHandler func(http.ResponseWriter, *http.Request, string);

type MicroService struct {
	Health func(w http.ResponseWriter, r *http.Request)	
	muxx *mux.Router
}

/** Default health check. This method can be overridden before the Start method
    is called. */
func ok(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

/** Instantiate a new microservice */
func NewMicroService() *MicroService {
    return &MicroService{
		muxx : mux.NewRouter(),
		Health: ok,
	}	
}

/** Register request */
func (m *MicroService) Handle(method string, path string, handler Handler) {
	fmt.Printf("Adding resource [%s] %s\n", method, path)
	m.muxx.Handle(path, AccessLogger{handler}).Methods(method)
}

/** Register principal request */
func (m *MicroService) Principal(method string, path string, handler PrincipalHandler) {
	fmt.Printf("Adding principal resource [%s] %s\n", method, path)
	m.muxx.Handle(path, AccessLogger{Principal{handler}.ServeHTTP}).Methods(method)
}
	
/** Not Allowed Requests */
func (ms *MicroService) NotAllowed(method string, path string) {
	fmt.Printf("NotAllowed resource [%s] %s\n", method, path)
	MethodNotAllowed := func (w http.ResponseWriter, r *http.Request) {	
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	ms.muxx.Handle(path, AccessLogger{MethodNotAllowed}).Methods(method)
}


/** Start a microservice with default health page. It uses port 8080 by
    convention. */
func (ms *MicroService) Start() {
	
	// add health
	ms.Handle("GET", "/health", ms.Health)
	 
	// start the web server
	if err := http.ListenAndServe(":8080", ms.muxx); err != nil {
		fmt.Println("error")
		log.Fatal("ListenAndServe:", err)
	} else {
		fmt.Println("running")
	}
}