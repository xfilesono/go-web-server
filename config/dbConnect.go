package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	db *sql.DB
)

func Connect() {

	passDB := envVariable("LOCALHOST_DB_PASSWD")
	userDB := envVariable("LOCALHOST_DB_USER")
	nameDB := envVariable("LOCALHOST_DB_NAME")

	conn := userDB + ":" + passDB + "@(127.0.0.1:3306)/" + nameDB

	var err error
	db, err = sql.Open("mysql", conn)
	if err = db.Ping(); err != nil {
		dbCheckErr(err)
	}
}

func GetDB() *sql.DB {
	fmt.Println("Database Connected")
	return db
}

func dbCheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func envVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	return os.Getenv(key)
}
