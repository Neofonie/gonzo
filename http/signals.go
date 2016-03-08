package gonzo

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

// Start signal handling
func init() {
	log.Println("Registering signal handler")
	go handleSignals()
}

func handleSignals() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for sig := range c {
		switch sig {

		default:
			log.Panicf("unexpected signal: %v", sig)

		case syscall.SIGINT, syscall.SIGTERM:
			log.Printf("received signal %#v: graceful shutdown...\n", sig)
			time.Sleep(2 * time.Second) // wait for Consul cluster to propagate the leave event
			shutdown()

		case syscall.SIGQUIT:
			fmt.Printf("received signal %#v: printing stacktraces...\n", sig)
			stacktrace()

		}
	}
}

func shutdown() {
	r := recover()
	if r != nil {
		fmt.Printf("panic: %v\n", r)
		stacktrace()
	}
	flushLogs()
	fmt.Println("Goodbye!")
	os.Exit(0)
}

func stacktrace() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	pprof.Lookup("heap").WriteTo(os.Stdout, 1)
	pprof.Lookup("threadcreate").WriteTo(os.Stdout, 1)
	pprof.Lookup("block").WriteTo(os.Stdout, 1)
}
