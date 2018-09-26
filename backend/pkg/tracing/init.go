package tracing

import (
  "fmt"
  "time"

  log "github.com/golang/glog"

  "github.com/uber/jaeger-client-go"
  "github.com/uber/jaeger-client-go/config"
  "github.com/opentracing/opentracing-go"
)

func Init(serviceName string, tracingClientHostPort string) opentracing.Tracer {
  // Create a opentracing.Tracer that sends data to our tracing backend
  fmt.Println(tracingClientHostPort)
  cfg := config.Configuration{
    Sampler: &config.SamplerConfig{
      Type: "const",
      Param:  1,
    },
    Reporter: &config.ReporterConfig{
      LocalAgentHostPort: tracingClientHostPort,
      LogSpans: true,
      BufferFlushInterval:  1 * time.Second,  // How often to flush traces
    },
  }
  tracer, _, err := cfg.New(
    serviceName,
    config.Logger(jaeger.StdLogger),
  )
  if (err != nil){
    log.Errorf("Could not connect to tracing client: %v", err.Error())
    panic("Tracing host is wrong")
  }

  return tracer
}
