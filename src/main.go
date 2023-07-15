package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/database"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/router"
)

func main() {
	loadEnv()
	StartServer()
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("読み込み出来ませんでした: %v", err)
	}
}

func StartServer() {
	const (
		port = ":8080"
	)
	db := database.ConnectDB()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	mux := router.NewRouter(db)
	log.Printf("started server on 0.0.0.0%v", port)
	http.ListenAndServe(port, mux)
}
