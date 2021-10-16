package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	sentrynegroni "github.com/getsentry/sentry-go/negroni"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"

	"github.com/benhall-1/appealscc/api/internal/db"
	"github.com/benhall-1/appealscc/api/routing"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	routing.SetupRequests(myRouter)

	n := negroni.Classic()
	n.UseHandler(myRouter)

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		TracesSampleRate: 0.2,
		TracesSampler: sentry.TracesSamplerFunc(func(ctx sentry.SamplingContext) sentry.Sampled {
			return sentry.SampledTrue
		}),
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	n.Use(sentrynegroni.New(sentrynegroni.Options{}))

	log.Fatal(http.ListenAndServe(":8080", n))
}

func main() {
	godotenv.Load()

	db.Open()
	db.Migrate()

	fmt.Println("AppealsCC API Server")
	handleRequests()
	sentry.CaptureMessage("It works!")
}
