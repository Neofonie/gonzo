package gonzo

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

var flushLogs = func() {}

// Init logging
func init() {
	w := newLogWriter()
	log.SetFlags(0) // do not log with time
	log.SetOutput(w)

	flushLogs = func() {
		w.Flush()
	}
	go flushLogsEvery(time.Second)
}

// flushLogsEvery writes the logs every d time units to its sink.
// The function never returns and must be run in a seperate goroutine.
func flushLogsEvery(d time.Duration) {
	for {
		time.Sleep(d)
		flushLogs()
	}
}

type logWriter struct {
	bufio.Writer
	rw sync.RWMutex
}

func newLogWriter() *logWriter {
	return &logWriter{
		Writer: *bufio.NewWriter(os.Stdout), // buffered logging
	}
}

// Write is passed as a Writer to the log package. The log package has its own
// syncronization for concurrent access. But since we want to call Flush from outside
// the log package Flush's access to the buffered writer must be exclusive to avoid data races.
func (l *logWriter) Write(p []byte) (n int, err error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	return l.Writer.Write(p)
}

// Flush forces the buffered data to be written.
func (l *logWriter) Flush() {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Writer.Flush()
}
