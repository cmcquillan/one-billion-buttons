package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	log.SetOutput(os.Stdout)

	args := os.Args[1:]

	verb := "create"

	if len(args) > 0 {
		verb = strings.ToLower(args[0])
	}

	connStr := os.Getenv("PG_CONNECTION_STRING")

	if len(connStr) == 0 {
		log.Fatal("required environment: PG_CONNECTION_STRING")
	}

	dbc, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Printf("postgres connect error: %v", err)
	}

	log.Print("pinging postgres")
	errPing := dbc.Ping()

	if errPing != nil {
		log.Printf("postgres ping error: %v", errPing)
	}

	log.Print("pinging postgres success")

	failure := false

	{
		defer dbc.Close()

		switch verb {
		case "create":
			if errCreate := ExecDir(dbc, "./migrations"); errCreate != nil {
				log.Printf("failed to execute db creation: %v", errCreate)
				failure = true
			}
		case "reset":
			if errReset := ExecDir(dbc, "./reset"); errReset != nil {
				log.Printf("failed to reset db: %v", errReset)
				failure = true
			}
		default:
			log.Printf("%v does not match a valid command", verb)
			failure = true
		}
	}

	if failure {
		os.Exit(1)
	}
}
