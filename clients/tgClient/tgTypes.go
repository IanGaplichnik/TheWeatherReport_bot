package tgClient

import (
	"net/http"
)

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text     string    `json:"text"`
	From     From      `json:"from"`
	Chat     Chat      `json:"chat"`
	Location *Location `json:"location"`
}

type From struct {
	Username string `json:"first_name"`
}

type Chat struct {
	ID int `json:"id"`
}

type Location struct {
	Lon float32 `json:"longitude"`
	Lat float32 `json:"latitude"`
}

type TgClient struct {
	host     string
	basePath string
	client   http.Client
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyKeboardMarkup struct {
	KeyboardButtons [][]KeyboardButton `json:"keyboard"`
	IsOnetime       bool               `json:"one_time_keyboard"`
}

type OutcomingKeyboard struct {
	ChatId      int                `json:"chat_id"`
	ReplyMarkup ReplyKeboardMarkup `json:"reply_markup"`
	Text        string             `json:"text"`
}
