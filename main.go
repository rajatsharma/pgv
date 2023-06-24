package main

import (
	"database/sql"
	"errors"
	"log"

	"github.com/ardanlabs/conf/v3"
	_ "github.com/lib/pq"
)

var cfg struct {
	DBUrl string `conf:"env:POSTGRES_URL,help:Postgres Url,required,flag:db-url,short:u"`
}

func main() {
	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			log.Println(help)
			return
		}
		log.Fatalf("parsing config: %v", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	rows, err := db.Query("SELECT VERSION()")
	if err != nil {
		log.Fatalf("cannot query database: %v", err)
	}
	var res string
	for rows.Next() {
		if err := rows.Scan(&res); err != nil {
			return
		}
		println(res)
	}

	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Fatalf("unable to close rows: %v", err)
		}
	}(rows)

	if err = db.Close(); err != nil {
		log.Fatalf("unable to close db connection: %v", err)
	}
}
