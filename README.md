# opencensus-gorilla_mux-example

Example demonstrating [OpenCensus] tracing with [Gorilla] mux routing.

## zipkin

This demo uses [Zipkin] as the tracing backend as it is very easy to set-up.
See the [Zipkin quickstart](https://zipkin.io/pages/quickstart)

## usage

A typical way to see the example in action:
```bash
# run Zipkin from official docker image
docker run -d -p 9411:9411 openzipkin/zipkin

# build the client and server applications
go generate

# start the server
build/server &

# run the client
build/client
```

If all is well you should be able to see some generated spans at http://localhost:9411


[gorilla]:(http://www.gorillatoolkit.org/)
[zipkin]:(https://zipkin.io)
[opencensus]:(https://opencensus.io)
