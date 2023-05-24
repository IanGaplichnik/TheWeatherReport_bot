package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"main.go/clients/events/telegram"
	"main.go/clients/geocoding"
	"main.go/clients/tgClient"
	"main.go/consumer/eventConsumer"
	"main.go/storage/sqlite"
	"main.go/tokens"
)

const (
	tgBotHost         = "api.telegram.org"
	storagePathSqlite = "data/sqlite/storage.db"
	geoClientHost     = "api.openweathermap.org"
	batchSize         = 100
)

func main() {
	tokens := new(tokens.Tokens)
	tokens.MustToken()

	s, err := sqlite.New(storagePathSqlite)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)

	}
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)

	}

	geocodingClient := geocoding.New(geoClientHost, tokens.GeocodingToken())

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, tokens.TgToken()),
		s,
		geocodingClient)

	log.Print("service started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
