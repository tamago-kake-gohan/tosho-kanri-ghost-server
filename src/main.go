package main

import (
	"log"
	"net/http"

	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
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
		log.Fatalf("データベースの疎通確認に失敗しました: %v", err)
	}
	session, _ := session.NewManager("memory", "gosessionid", 3600)
	mux := router.NewRouter(db, session)
	log.Printf("started server on 0.0.0.0%v", port)
	http.ListenAndServe(port, mux)
	go session.GC()
}
