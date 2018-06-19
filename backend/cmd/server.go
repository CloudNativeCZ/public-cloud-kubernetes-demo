package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/api"
	"github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/pkg"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

type AppOptions struct {
	backingStoreOptions *redis.Options
	listenHost          int
	listenPort          int
	backingStoreHost    string
	backingStorePort    int
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
	flag.IntVar(&app.options.listenHost, "listenHost", 0, "Host address to listen on")
	flag.IntVar(&app.options.listenPort, "listenPort", 8080, "Host port to listen on")
	flag.StringVar(&app.options.backingStoreHost, "backingStoreHost", "localhost", "Redis host to connect to")
	flag.IntVar(&app.options.backingStorePort, "backingStorePort", 6379, "Redis port to connect to")
	flag.IntVar(&app.options.backingStoreOptions.DB, "db", 0, "Redis database to use")
}

func (app *App) parseFlags() {
	flag.Parse()

	app.options.backingStoreOptions.Addr = fmt.Sprintf("%s:%d", app.options.backingStoreHost, app.options.backingStorePort)
}

func (app *App) parseEnvVars() {
	password := os.Getenv("BACKING_STORE_PASSWORD")

	app.options.backingStoreOptions.Password = password
}

func (app *App) initiateBackingStore() {
	backingStore, err := pkg.NewClient(app.options.backingStoreOptions)
	if err != nil {
		panic("Could not connect to backing store")
	}

	app.backingStore = backingStore
}

func (app *App) serve() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%d:%d", app.options.listenHost, app.options.listenPort), nil))
}

func main() {
	app := newApp()
	app.addFlags()
	app.parseFlags()
	app.parseEnvVars()

	app.initiateBackingStore()

	questionsResource := api.NewQuestionsResource(app.backingStore)
	questionsResource.Register()

	app.serve()
}
