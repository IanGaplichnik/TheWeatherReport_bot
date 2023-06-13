package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode"

	"main.go/clients/events"
	"main.go/clients/tgClient"
	"main.go/lib/e"
	"main.go/storage"
	"main.go/storage/sqlite"
)

const (
	HelpCmd      = "/help"
	StartCmd     = "/start"
	CheckWeather = "/checkweather"
	CurrentCity  = "/currentcity"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("command %s from user %s", text, username)

	userdata := &storage.Userdata{
		UserName: username,
		ChatId:   chatID,
	}

	isKeyboard, err := p.storage.IsKeyboardMenu(ctx, userdata)
	if err != nil {
		return e.Wrap("can't do cmd", err)
	}

	if isCity(text) {
		return p.setCity(ctx, text, userdata)
	}

	if len(text) > 1 && unicode.IsDigit(rune(text[0])) && isKeyboard {
		return p.executeKeyboardInput(ctx, userdata, text)
	}

	switch text {
	case HelpCmd:
		return p.tg.SendMessage(chatID, helpMsg)
	case StartCmd:
		return p.tg.SendMessage(chatID, helloMsg)
	case CheckWeather:
		return p.checkWeather(ctx, userdata)
	case CurrentCity:
		return p.checkCurrentCity(ctx, userdata)
	default:
		return p.tg.SendMessage(chatID, unknownCommandMsg)
	}
}

func (p *Processor) executeKeyboardInput(ctx context.Context, userdata *storage.Userdata, text string) error {
	if len(text) < 3 {
		return e.Wrap("can't execute Keyboard input", errors.New("input too short"))
	}

	reply := []rune(text)
	reply[2] = ','
	text = string(reply)

	stats := strings.Split(text, ",")
	for i := 0; i < len(stats); i++ {
		stats[i] = strings.Trim(stats[i], " ")
		println(stats[i])
	}
	var result string = stats[1]

	for i := 2; i < len(stats); i++ {
		result = result + "," + stats[i]
	}
	println(result)

	cities, err := p.geocoding.FetchCitiesByCityName(result)
	if err != nil {
		return e.Wrap("can't execute keyboard input", err)
	}
	fmt.Printf("amount of cities keyboard input = %d\n", len(cities))

	userdata.Country = cities[0].Country
	userdata.City = cities[0].CityName
	userdata.State = cities[0].State

	p.storage.Save(ctx, userdata)
	return nil
}

func (p *Processor) checkWeather(ctx context.Context, userdata *storage.Userdata) error {

	cityData, err := p.storage.RetrieveLocation(ctx, userdata)
	if err != nil {
		return e.Wrap("can't check if it'll rain", err)
	}
	if cityData == nil {
		return p.tg.SendMessage(userdata.ChatId, "You need to set the city first :)")
	}

	var state string = ""
	if cityData.State != "" {
		state = "," + cityData.State
	}
	city := cityData.CityName + state + "," + cityData.Country

	coord, err := p.geocoding.FetchCoordsByCity(city)
	if err != nil {
		return e.Wrap("can't check if rain", err)
	}

	message, err := p.geocoding.FetchWeather(coord.Latitude, coord.Longitude)
	if err != nil {
		return e.Wrap("can't check if rain", err)
	}

	p.tg.SendMessage(userdata.ChatId, message)

	return nil
}

func (p *Processor) checkCurrentCity(ctx context.Context, userdata *storage.Userdata) error {
	cityData, err := p.storage.RetrieveLocation(ctx, userdata)
	if err != nil {
		return e.Wrap("can't retreive city", err)
	}
	if cityData == nil {
		p.tg.SendMessage(userdata.ChatId, "You need to set the city first :)")
		return nil
	}

	var state string = ""
	if cityData.State != "" {
		state = ", " + cityData.State
	}
	cityName := cityData.CityName + state + ", " + cityData.Country
	p.tg.SendMessage(userdata.ChatId, "Your city is "+cityName)
	return nil
}

func (p *Processor) setCity(ctx context.Context, text string, userdata *storage.Userdata) error {

	cities, err := p.geocoding.FetchCitiesByCityName(text)
	if err != nil {
		return e.Wrap("error occured getting the city", err)
	}

	if len(cities) == 0 {
		if err := p.tg.SendMessage(userdata.ChatId, msgNoCity); err != nil {
			return e.Wrap("can't set city", err)
		}
		return nil
	}

	if len(cities) > 1 {
		if err := p.handleMultipleCities(ctx, cities, userdata.ChatId, userdata.UserName); err != nil {
			return e.Wrap("can't set city", err)
		}
		return nil
	}

	if len(cities) == 1 {
		userdata.City = cities[0].CityName
		userdata.Country = cities[0].Country
		if cities[0].State == "" {
			userdata.State = cities[0].State
		}
		userdata.Menu = sqlite.MenuNot
		if err := p.saveCityToDB(ctx, userdata); err != nil {
			return e.Wrap("can't save city to db", err)
		}
	}

	return nil
}

func (p *Processor) handleMultipleCities(ctx context.Context, cities []events.CityData, chatID int, username string) error {
	kbMarkup := buildCitiesKeyboard(cities)

	err := p.tg.SendKeyboard(chatID, *kbMarkup)
	if err != nil {
		return e.Wrap("can't handle multiple cities", err)
	}

	userData := storage.Userdata{
		UserName: username,
		ChatId:   chatID,
		Menu:     sqlite.MenuKeyboard,
		City:     "",
	}

	if err := p.storage.Save(ctx, &userData); err != nil {
		return e.Wrap("can't handle multiple cities", err)
	}

	return nil
}

func buildCitiesKeyboard(cities []events.CityData) *tgClient.ReplyKeboardMarkup {

	kbMarkup := new(tgClient.ReplyKeboardMarkup)
	kbMarkup.KeyboardButtons = make([][]tgClient.KeyboardButton, len(cities))
	kbMarkup.IsOnetime = true

	for i, city := range cities {
		kbMarkup.KeyboardButtons[i] = make([]tgClient.KeyboardButton, 1)
		if city.State == "" {
			kbMarkup.KeyboardButtons[i][0].Text = fmt.Sprintf("%d. %s, %s", i+1, city.CityName, city.Country)
		} else {
			kbMarkup.KeyboardButtons[i][0].Text = fmt.Sprintf("%d. %s, %s, %s", i+1, city.CityName, city.State, city.Country)
		}
	}

	return kbMarkup
}

func (p *Processor) saveCityToDB(ctx context.Context, userdata *storage.Userdata) error {
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

	err = p.storage.Save(ctx, userdata)
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
		if !unicode.IsLetter(symbol) && symbol != ' ' && symbol != '-' {
			fmt.Println("Not a city!")
			return false
		}
	}

	return true
}
