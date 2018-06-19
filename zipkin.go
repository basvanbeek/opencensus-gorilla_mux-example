package ocmux

import (
	"io"

	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	oczipkin "go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"
)

const (
	zipkinURL = "http://localhost:9411/api/v2/spans"
)

// InitOpenCensusWithZipkin initializes the OpenCensus Zipkin Exporter.
func InitOpenCensusWithZipkin(serviceName, hostPort string) io.Closer {
	// Always trace for this demo.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	rep := zipkinhttp.NewReporter(zipkinURL)

	localEndpoint, _ := zipkin.NewEndpoint(serviceName, hostPort)

	exporter := oczipkin.NewExporter(rep, localEndpoint)
	trace.RegisterExporter(exporter)

	return rep
}
