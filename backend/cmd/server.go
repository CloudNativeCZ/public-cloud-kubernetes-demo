package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	log "github.com/golang/glog"

	"github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/api"
  "github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/pkg"
  "github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/pkg/tracing"
	"github.com/go-redis/redis"

	"github.com/emicklei/go-restful"

  // Tracing dependencies
  "github.com/opentracing/opentracing-go"
)

type AppOptions struct {
	backingStoreOptions *redis.Options
	listenHost                string
	listenPort                int
  tracingClientHostPort     string
}

type App struct {
	backingStore *redis.Client
	options      *AppOptions
}

func newApp() *App {
	return &App{
		options: &AppOptions{
			backingStoreOptions: &redis.Options{},
		},
	}
}

func (app *App) addFlags() {
	flag.StringVar(&app.options.listenHost, "listenHost", "0.0.0.0", "Host address to listen on")
	flag.IntVar(&app.options.listenPort, "listenPort", 8080, "Host port to listen on")
}

func (app *App) parseFlags() {
	flag.Parse()
}

func (app *App) parseEnvVars() {
	host := os.Getenv("BACKING_STORE_HOST")
	port := os.Getenv("BACKING_STORE_PORT")
	user := os.Getenv("BACKING_STORE_DB")
	db, err := strconv.Atoi(user)
	if err != nil {
		panic("Wrong DB name")
	}

	app.options.backingStoreOptions.Addr = fmt.Sprintf("%s:%s", host, port)
	app.options.backingStoreOptions.DB = db

  // Pass tracing client backend via environment variables
  tracingClientHost := os.Getenv("TRACING_CLIENT_HOST")
  tracingClientPort := os.Getenv("TRACING_CLIENT_PORT")

  app.options.tracingClientHostPort = fmt.Sprintf("%s:%s", tracingClientHost, tracingClientPort)
}

func (app *App) initiateBackingStore() {
	backingStore, err := pkg.NewClient(app.options.backingStoreOptions)
	if err != nil {
		log.Errorf("Could not connect to redis: %s", err)
	}

	app.backingStore = backingStore
}

func (app *App) serve(container *restful.Container) {
	addr := fmt.Sprintf("%s:%d", app.options.listenHost, app.options.listenPort)
	server := &http.Server{Addr: addr, Handler: container}

	log.Fatal(server.ListenAndServe())
}

func main() {
  app := newApp()
  app.addFlags()
  app.parseFlags()
  app.parseEnvVars()

  tracer := tracing.Init("server", app.options.tracingClientHostPort)

  parent := tracer.StartSpan("Boostrap")

  child := tracer.StartSpan("BackingStoreBoostrap", opentracing.ChildOf(parent.Context()))
  app.initiateBackingStore()
  child.LogEvent("Backing store initiated")
  child.Finish()

  child = tracer.StartSpan("WSContainerBoostrap", opentracing.ChildOf(parent.Context()))
  wsContainer := restful.NewContainer()
  child.LogEvent("WS Container created")

  // Add container filter to enable CORS
  cors := restful.CrossOriginResourceSharing{
    ExposeHeaders:  []string{"X-My-Header"},
    AllowedHeaders: []string{"Content-Type", "Accept"},
    AllowedMethods: []string{"GET", "PUT"},
    CookiesAllowed: false,
    Container:      wsContainer}
  wsContainer.Filter(cors.Filter)
  child.LogEvent("WS Container CORS Filter created")

  // Add container filter to respond to OPTIONS
  wsContainer.Filter(wsContainer.OPTIONSFilter)
  child.LogEvent("WS Container CORS OPTIONS Filter added")

  child.LogEvent("App boot finished")
  child.Finish()

  child = tracer.StartSpan("QuestionsApiBoostrap", opentracing.ChildOf(parent.Context()))
  questionsResource := api.NewQuestionsResource(app.backingStore, app.options.tracingClientHostPort)
  child.LogEvent(fmt.Sprintf("Question Resource started with %v", app.backingStore))
  questionsResource.Register(wsContainer)

  child.LogEvent(fmt.Sprintf("Question Resource registered with WS Container %v", wsContainer))
  child.LogEvent("Starting serving content")
  child.Finish()
  parent.Finish()

  app.serve(wsContainer)
}
