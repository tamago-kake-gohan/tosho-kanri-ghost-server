package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	loadEnv()
	fmt.Println("test")
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}
