package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/benhall-1/appealscc/api/internal/db"
	"github.com/benhall-1/appealscc/api/routing"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	routing.SetupRequests(myRouter)

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		TracesSampleRate: 0.2,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	godotenv.Load()

	db.Open()
	db.Migrate()

	fmt.Println("AppealsCC API Server")
	handleRequests()
	sentry.CaptureMessage("It works!")
}
