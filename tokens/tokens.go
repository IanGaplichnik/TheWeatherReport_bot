package tokens

import (
	"flag"
	"log"
)

type Tokens struct {
	tgToken *string
	gcToken *string
}

func (t *Tokens) MustToken() {
	t.tgToken = flag.String("tg-bot-token", "", "telegram bot access token (from BotFather)")
	t.gcToken = flag.String("geocoding-token", "", "geocoding-token")
	flag.Parse()
	if t.tgToken == nil || *t.tgToken == "" {
		log.Fatal("No tgToken entered")
	}
	if t.gcToken == nil || *t.gcToken == "" {
		log.Fatal("No geocodingToken entered")
	}
}

func (t *Tokens) GeocodingToken() string {
	return *t.gcToken
}

func (t *Tokens) TgToken() string {
	return *t.tgToken
}
