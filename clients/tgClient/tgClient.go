package tgClient

import (
	"encoding/json"
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
)

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
