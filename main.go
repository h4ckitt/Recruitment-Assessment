package main

import (
	"assessment/config"
	"assessment/infra/db/sqlite"
	"assessment/interface/mux/controller"
	"assessment/interface/mux/router"
	"assessment/service"
	"log"
	"net/http"
	"os"
)

func main() {
	err := config.LoadEnv()

	if err != nil {
		log.Fatalf("An error occurred while trying to load config file: %v\n", err)
	}

	repo, err := sqlite.NewSqliteClient()

	if err != nil {
		log.Fatalf("An error occurred while bringing up the repository: %v\n", err)
	}

	svc := service.NewNumberService(service.NewValidator(), repo)

	numController := controller.NewNumberController(svc)

	r := router.InitRouter(numController)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Starting Server...")
	if err = http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalln(err)
	}
}
