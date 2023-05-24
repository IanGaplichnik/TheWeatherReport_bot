package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode"

	"main.go/lib/e"
	"main.go/storage"
)

const (
	HelpCmd     = "/help"
	StartCmd    = "/start"
	SetCity     = "/setcity"
	GetWeather  = "/getweather"
	CheckRain   = "/checkrain"
	CurrentCity = "/currentcity"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("command %s from user %s", text, username)

	if isCity(text) {
		return p.setCity(ctx, text, chatID, username)
	}

	switch text {
	case HelpCmd:
		return p.tg.SendMessage(chatID, helpMsg)
	case StartCmd:
		return p.tg.SendMessage(chatID, helloMsg)
	// case GetWeather: //TODO
	default:
		return p.tg.SendMessage(chatID, unknownCommandMsg)
	}

}

func (p *Processor) setCity(ctx context.Context, text string, chatID int, username string) error {
	cities, err := p.geocoding.GeoCoordinates(text)
	if err != nil {
		return e.Wrap("error occured getting the city", err)
	}

	if len(cities) == 0 {
		p.tg.SendMessage(chatID, "Couldn't find such city!")
		return nil
	}

	userdata := storage.Userdata{
		UserName: username,
		City:     cities[0].Name,
	}

	if len(cities) == 1 {
		if userdata.City != text {
			p.tg.SendMessage(chatID, "Couldn't find such city!")
			return nil
		}
		p.saveCityToDB(ctx, userdata)
		message := fmt.Sprintf("City %s is succesfully set!", userdata.City)
		p.tg.SendMessage(chatID, message)
	}

	return nil
}

func (p *Processor) saveCityToDB(ctx context.Context, userdata storage.Userdata) error {
	exists, err := p.storage.Exists(ctx, userdata)
	if err != nil {
		return e.Wrap("can't check if user exists", err)
	}

	if exists {
		err := p.storage.Remove(ctx, userdata)
		if err != nil {
			return e.Wrap("can't remove entry from db", err)
		}
	}

	err = p.storage.Save(ctx, &userdata)
	if err != nil {
		return e.Wrap("can't save userdata to storage", err)
	}
	return nil
}

func isCity(text string) bool {
	if len(text) == 0 {
		return false
	}

	for _, symbol := range text {
		if !unicode.IsLetter(symbol) {
			return false
		}
	}

	return true
}
