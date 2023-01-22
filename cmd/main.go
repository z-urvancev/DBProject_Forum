package main

import (
	"DBProject/internal/router"
	"DBProject/pkg/postgres"
	"log"
	"time"

	"DBProject/internal/handlers"
	"DBProject/internal/repositories"
	"DBProject/internal/usecases"
)

func main() {
	time.Sleep(5 * time.Second)
	log.Println("connecting to database...")
	db, err := postgres.NewConnect()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("connection successful")
	defer db.Close()

	repositories := repositories.NewRepositories(db)

	useCases := usecases.NewUseCases(repositories)

	handlers := handlers.NewHandlers(useCases)

	engine := router.NewEngine(handlers)
	log.Println("server is started")
	err = engine.Start(":5000")
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("server is down")
}
