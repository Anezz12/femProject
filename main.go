package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Anezz12/femProject/internal/app"
	"github.com/Anezz12/femProject/internal/routes"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	defer app.DB.Close()

	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Println(fmt.Sprintf("Application started at port %d", port))

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf("Error starting server: %s", err)
	}

}
