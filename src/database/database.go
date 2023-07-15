package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("ローケルの取得に失敗しました: %v", err)
	}
	cfg := mysql.Config{
		DBName:               dbname,
		User:                 user,
		Passwd:               pass,
		Addr:                 host,
		Net:                  "tcp",
		Loc:                  jst,
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if nil != err {
		log.Fatalf("データベースへの接続に失敗しました: %v", err)
	}
	return db
}
