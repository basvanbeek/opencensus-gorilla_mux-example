package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"

	ocmux "github.com/basvanbeek/opencensus-gorilla_mux-example"
)

const (
	host = "localhost:8000"
)

type client struct {
	c *http.Client
	r *mux.Router
}

func main() {
	// Initialize OpenCensus using Zipkin tracing backend
	defer ocmux.InitOpenCensusWithZipkin("client", "localhost:0").Close()

	// Create our Gorilla Mux template powered HTTP client
	client := newClient(host)

	// let's call someRoute
	params := []string{
		"object_id", "42",
	}
	client.call("someRoute", params...)

	// let's call otherRoute
	params = []string{
		"object_id", "42",
		"instance_id", "4937",
		"filter", "some filter text",
	}
	client.call("otherRoute", params...)
}

func newClient(host string) *client {
	router := mux.NewRouter()

	router.Schemes("HTTP").Methods("GET").Host(host).
		Path("/object/{object_id}").
		Name("someRoute")

	router.Schemes("HTTP").Methods("POST").Host(host).
		Path("/object/{object_id}/instance/{instance_id}").
		Queries("filter", "{filter}").
		Name("otherRoute")

	transport := &ochttp.Transport{
		FormatSpanName: ocmux.NameFromGorillaMux(router),
	}
	httpClient := &http.Client{Transport: transport}

	return &client{c: httpClient, r: router}
}

func (c *client) call(routeName string, pairs ...string) {
	var (
		methods []string
		err     error
		req     *http.Request
		res     *http.Response
	)

	fmt.Printf("REQUEST : %s\nPARAMS  : %#v\n", routeName, pairs)

	route := c.r.GetRoute(routeName)
	if route == nil {
		fmt.Printf("ERROR: unable to find route %s\n", routeName)
		return
	}

	if methods, err = route.GetMethods(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	if req, err = http.NewRequest(methods[0], "", nil); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	if req.URL, err = route.URL(pairs...); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	if res, err = c.c.Do(req); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("RESPONSE: %s\n\n", body)
}
