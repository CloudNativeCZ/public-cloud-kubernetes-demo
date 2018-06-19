package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"log"
	"net/http"

	"github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/api"
	"github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/pkg"
	"github.com/go-redis/redis"
)

type AppOptions struct {
	backingStoreOptions *redis.Options
	listenHost          int
	listenPort          int
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
}

func (app *App) parseFlags() {
        flag.Parse()
}

func (app *App) parseEnvVars() {
        password := os.Getenv("BACKING_STORE_PASSWORD")
        host := os.Getenv("BACKING_STORE_HOST")
        port := os.Getenv("BACKING_STORE_PORT")
        user := os.Getenv("BACKING_STORE_DB")
        db, err := strconv.Atoi(user)
        if err != nil {
                panic("Wrong DB name")
        }

        app.options.backingStoreOptions.Password = password
        app.options.backingStoreOptions.Addr = fmt.Sprintf("%s:%d", host, port)
        app.options.backingStoreOptions.DB = db
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
