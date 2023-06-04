package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	db *sql.DB
)

func Connect() (*sql.DB, error) {

	passDB := envVariable("LOCALHOST_DB_PASSWD")
	userDB := envVariable("LOCALHOST_DB_USER") + ":"
	nameDB := envVariable("LOCALHOST_DB_NAME")
	server := "@(127.0.0.1:3306)/"

	conn := userDB + passDB + server + nameDB

	var err error
	db, err = sql.Open("mysql", conn)
	if err = db.Ping(); err != nil {
		dbCheckErr(err)
	} else {
		fmt.Println("Database Connected")
	}

	return db, err
}

func GetDB() (*sql.DB, error) {
	return Connect()
}

func dbCheckErr(err error) {
	if err != nil {
		log.Fatalln(http.StatusUnauthorized)
	}
}

func envVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	return os.Getenv(key)
}
