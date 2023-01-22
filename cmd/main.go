package main

import (
	"DBProject/internal/router"
	"DBProject/pkg/postgres"
	"fmt"
	"log"

	"DBProject/internal/handlers"
	"DBProject/internal/repositories"
	"DBProject/internal/usecases"
)

func main() {
	db, err := postgres.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repositories := repositories.NewRepositories(db)

	useCases := usecases.NewUseCases(repositories)

	handlers := handlers.NewHandlers(useCases)

	engine := router.NewEngine(handlers)

	err = engine.Start(":80")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
