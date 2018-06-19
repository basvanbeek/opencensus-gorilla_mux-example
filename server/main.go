package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"go.opencensus.io/plugin/ochttp"

	ocmux "github.com/basvanbeek/opencensus-gorilla_mux-example"
)

const (
	tplSome  = `Serving SomeRoute. Object %q`
	tplOther = `Serving OtherRoute. Object %q, Instance %q, Filter %q`
)

func main() {
	// Initialize OpenCensus using Zipkin tracing backend
	defer ocmux.InitOpenCensusWithZipkin("server", "localhost:8000").Close()

	router := mux.NewRouter()
	router.Use(ocmux.Middleware())
	router.Methods("GET").
		Path("/object/{object_id}").
		HandlerFunc(someRoute)

	router.Methods("POST").
		Path("/object/{object_id}/instance/{instance_id}").
		Queries("filter", "{filter}").
		HandlerFunc(otherRoute)

	handler := &ochttp.Handler{Handler: router}

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("unable to create listener: %v", err)
	}

	// run.Group manages our goroutine lifecycles
	// see: https://www.youtube.com/watch?v=LHe1Cb_Ud_M&t=15m45s
	var g run.Group
	// set-up our HTTP service handler
	{
		g.Add(func() error {
			return http.Serve(listener, handler)

		}, func(error) {
			listener.Close()
		})
	}
	// set-up our signal handler
	{
		var (
			cancelInterrupt = make(chan struct{})
			c               = make(chan os.Signal, 2)
		)
		defer close(c)

		g.Add(func() error {
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	// spawn our goroutines and wait for shutdown
	log.Println("exit", g.Run())

}

func someRoute(w http.ResponseWriter, r *http.Request) {
	var (
		vars     = mux.Vars(r)
		objectID = vars["object_id"]
	)
	fmt.Fprintf(w, tplSome, objectID)
}

func otherRoute(w http.ResponseWriter, r *http.Request) {
	var (
		vars       = mux.Vars(r)
		objectID   = vars["object_id"]
		instanceID = vars["instance_id"]
		filter     = vars["filter"]
	)
	fmt.Fprintf(w, tplOther, objectID, instanceID, filter)
}
