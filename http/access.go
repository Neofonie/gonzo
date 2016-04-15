package gonzo

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

/** Handle principal calls */
type AccessLogger struct {
	next Handler
}

// Wrap the request with access logging
func (p AccessLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	accessLog := prepareAccessLog(r)
	defer writeAccessLog(accessLog, start)
	myRW := &responseWriter{ResponseWriter: w}
	p.next(myRW, r)
	accessLog["status"] = myRW.statusCode
}

// We need this to remember the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Write the access-log to stdout
func writeAccessLog(accessLog map[string]interface{}, start time.Time) {

	durationMillis := float64(time.Now().Sub(start).Nanoseconds()) / 1000000.0
	accessLog["durationMillis"] = durationMillis

	if result, err := json.Marshal(accessLog); err == nil {
		log.Println(string(result))
	} else {
		log.Printf("cannot convert access log to json: %v\n", err)
	}
}

// Create a map with relevant access-log data
func prepareAccessLog(req *http.Request) map[string]interface{} {

	result := make(map[string]interface{})

	for k, v := range req.Header {
		result[k] = strings.Join(v, ", ")
	}

	result["requestUri"] = req.RequestURI
	result["log-type"] = "access"
	result["remoteAddress"] = req.RemoteAddr
	result["requestMethod"] = req.Method
	result["originAddress"] = originAddress(req)

	return result
}

// Get last part of the Forwarded-For field
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
