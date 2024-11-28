package main

import (
	"log"
	config "petitionsGO/configs"
	"petitionsGO/db"
	"petitionsGO/handlers"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg, err := config.LoadConfig(`configs/config.json`)
	if err != nil {
		log.Fatalf("Failed configuration: %v", err)
	}

	database := config.ConnectDB(cfg)
	defer database.Close()
	log.Println("База данных подключена")

	queries := db.New(database)

	handler := &handlers.Handler{
		Queries:      queries,
		JWTSecretKey: cfg.JWTSecret,
	}

	server := &Server{}
	if err := server.Run(cfg.Port, handler.InitRoutes()); err != nil {
		log.Fatalf("Error server: %v", err)
	}
}
