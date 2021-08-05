package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/snowflake"
	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

// Base de donnée (définie globalement)
var database *runner.DB
var node *snowflake.Node

func init() {
	// We connect to the correct database for this challenge.
	var db *sql.DB
	var err error
	fmt.Println("Opening.")
	db, err = sql.Open("postgres", "dbname=ctf0 user=postgres password="+dbpswd+" host=db port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}
	fmt.Println("Opened.")

	runner.MustPing(db)

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)

	dat.EnableInterpolation = true

	// DB is PostgreSQL
	database = runner.NewDB(db, "postgres")
}

func main() {
	go hostService()

	// We load a node to ensure that everyone is able to do their SQL trick.
	node, _ = snowflake.NewNode(1)

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
