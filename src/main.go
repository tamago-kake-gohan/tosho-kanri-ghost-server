package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/database"
	"tamago-kake-gohan.github.io/tosho-kanri-ghost/src/router"
)

func main() {
	loadEnv()
	fmt.Println("test")
	StartServer()
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}

func StartServer() {
	const (
		port = ":8080"
	)
	db := database.ConnectDB()
	mux := router.NewRouter(db)
	http.ListenAndServe(port, mux)
}
