package tgClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"main.go/lib/e"
)

const (
	sendMessageMethod = "sendMessage"
	getUpdatesMethod  = "getUpdates"
	sendJson          = "application/json"
)

func (c *TgClient) SendKeyboard(chatId int, kbMarkup ReplyKeboardMarkup) error {

	keyboardMsg := OutcomingKeyboard{
		ChatId:      chatId,
		ReplyMarkup: kbMarkup,
		Text:        "Choose the correct city",
	}

	payload, err := json.MarshalIndent(keyboardMsg, "", " ")
	if err != nil {
		return e.Wrap("can't send keyboard", err)
	}

	if err := c.doPostRequest(payload); err != nil {
		return e.Wrap("can't send keyboard", err)
	}

	return nil
}

func (c *TgClient) doPostRequest(payload []byte) error {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, sendMessageMethod),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(payload))
	if err != nil {
		return e.Wrap("can't build a request", err)
	}
	fmt.Println(req)

	resp, err := http.Post(u.String(), sendJson, bytes.NewReader(payload))
	if err != nil {
		return e.Wrap("can't client.Do request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return e.Wrap("can't ReadAll", err)
	}

	return nil
}

func (c *TgClient) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfError("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *TgClient) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func New(host, token string) *TgClient {
	return &TgClient{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *TgClient) doRequest(methodFunction string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme:   "https",
		Host:     c.host,
		Path:     path.Join(c.basePath, methodFunction),
		RawQuery: query.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap("can't build a request", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap("can't client.Do request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't ReadAll", err)
	}

	return body, err
}
