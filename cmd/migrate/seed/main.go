package main

import (
	"gopher_social/internal/db"
	"gopher_social/internal/env"
	"gopher_social/internal/store"
	"log"
	"time"
)

type seedDBConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
	maxLifeTime  time.Duration
}

func main() {
	var seedDBConfig = seedDBConfig{
		addr:         env.GetString("DATABASE_URL", "postgres://admin:adm1npass@localhost/social?sslmode=disable"),
		maxOpenConns: 3,
		maxIdleConns: 3,
		maxIdleTime:  time.Minute,
		maxLifeTime:  2 * time.Minute,
	}

	seedDB, err := db.NewConn(
		seedDBConfig.addr,
		seedDBConfig.maxOpenConns,
		seedDBConfig.maxIdleConns,
		seedDBConfig.maxIdleTime,
		seedDBConfig.maxLifeTime,
	)

	if err != nil {
		log.Fatal(err)
	}
	defer seedDB.Close()

	store := store.NewStorage(seedDB)

	db.Seed(store, seedDB)
}
