package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// loads env details
func loadEnv() (string, int, string, string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("could not load environment file")
	}

	port, _ := strconv.Atoi(os.Getenv("port"))

	return os.Getenv("host"), port, os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname")
}

// sets up a database connection
func database() {
	host, port, user, password, dbname := loadEnv()

	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

	db, err = sql.Open("postgres", info)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Established a successful connection!")
}
