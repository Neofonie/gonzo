package gonzo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

/** Handle principal calls */
type AccessLogger struct {
	next func(http.ResponseWriter, *http.Request)
}

/** Wrap the server with the access-log */
func (p AccessLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer writeAccessLog(prepareAccessLog(r), start)
	p.next(w, r)
}

/** Write the access-log to stdout */
func writeAccessLog(accessLog map[string]string, start time.Time) {

	durationMillis := float64(time.Now().Sub(start).Nanoseconds()) / 1000000.0
	accessLog["durationMillis"] = fmt.Sprint(durationMillis)

	if result, err := json.Marshal(accessLog); err == nil {
		log.Println(string(result))
	} else {
		log.Printf("cannot convert access log to json: %v\n", err)
	}
}

/** Create a map with relevant access-log data */
func prepareAccessLog(req *http.Request) map[string]string {

	result := map[string]string{
		"requestUri":    req.RequestURI,
		"log-type":      "access",
		"remoteAddress": req.RemoteAddr,
		"requestMethod": req.Method,
		"service":       "proxy",
		"originAddress": originAddress(req),
	}

	for k, v := range req.Header {
		result[k] = strings.Join(v, ", ")
	}

	return result
}

/** Get last part of the Forwarded-For field */
func originAddress(req *http.Request) string {

	xff := req.Header.Get("X-Forwarded-For")

	if xff == "" {
		return req.RemoteAddr
	}

	i := strings.Index(xff, ",")
	if i == -1 {
		return xff
	}

	return xff[:i]
}
