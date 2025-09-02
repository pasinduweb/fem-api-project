package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/pasinduweb/fem-api-project/internal/app"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 3001, "Go server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/health", HealthCheck)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("App is running on port %d\n", port)

	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available\n")
}
