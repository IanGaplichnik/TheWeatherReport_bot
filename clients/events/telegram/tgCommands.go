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
	// case CheckRain:
	case CurrentCity:
		return p.checkCurrentCity(ctx, chatID, username)
	default:
		return p.tg.SendMessage(chatID, unknownCommandMsg)
	}
}

func (p *Processor) checkCurrentCity(ctx context.Context, chatId int, username string) error {
	userdata := storage.Userdata{
		ChatId:   chatId,
		UserName: username,
	}
	exists, err := p.storage.Exists(ctx, userdata)
	if err != nil {
		e.Wrap("can't check if user exists", err)
	}

	if !exists {
		return p.tg.SendMessage(chatId, "No city is set :(")
	}

	cityname, err := p.storage.RetrieveCity(ctx, userdata)
	if err != nil {
		return e.Wrap("can't retreive city", err)
	}

	p.tg.SendMessage(chatId, "Your city is "+cityname)
	return nil
}

func (p *Processor) setCity(ctx context.Context, text string, chatID int, username string) error {

	cities, err := p.geocoding.FetchCity(text)
	if err != nil {
		return e.Wrap("error occured getting the city", err)
	}

	if len(cities) == 0 {
		p.tg.SendMessage(chatID, msgNoCity)
		return nil
	}

	userdata := storage.Userdata{
		UserName: username,
		ChatId:   chatID,
		City:     cities[0].CityName,
	}

	if len(cities) == 1 {
		if userdata.City != text {
			if err := p.tg.SendMessage(chatID, msgNoCity); err != nil {
				return e.Wrap("can't send message", err)
			}
			return nil
		}

		fmt.Print("HEREEEE")
		if err := p.saveCityToDB(ctx, userdata); err != nil {
			return e.Wrap("can't save city to db", err)
		}
		fmt.Print("HEREEEE2")
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

	message := fmt.Sprintf("City of %s is succesfully set!🙌\nYou can choose the next action from the menu.\nSend a message or a location to set a new city", userdata.City)
	if err := p.tg.SendMessage(userdata.ChatId, message); err != nil {
		return e.Wrap("can't send message", err)
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
