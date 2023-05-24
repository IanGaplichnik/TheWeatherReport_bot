package telegram

import (
	"context"
	"errors"

	"main.go/clients/events"
	"main.go/clients/geocoding"
	"main.go/clients/tgClient"
	"main.go/lib/e"
	"main.go/storage"
)

type Processor struct {
	tg        *tgClient.TgClient
	geocoding *geocoding.GeocodingClient
	offset    int
	storage   storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("event type unknown")
	ErrUnknownMetaType  = errors.New("meta mype unknown")
)

func New(client *tgClient.TgClient, storage storage.Storage, geocoding *geocoding.GeocodingClient) *Processor {
	return &Processor{
		tg:        client,
		storage:   storage,
		geocoding: geocoding,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, update := range updates {
		res = append(res, event(update))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.Location:
		return p.processLocation(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}
func (p *Processor) processLocation(event events.Event) error {

}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(context.Background(), event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd tgClient.Update) events.Event {
	res := events.Event{
		Type:     fetchType(upd),
		Text:     fetchText(upd),
		Location: fetchLocation(upd),
	}

	if res.Type == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchLocation(upd tgClient.Update) *events.Coordinates {
	if upd.Message == nil {
		return nil
	}
	if upd.Message.Location != nil {
		return &events.Coordinates{
			Longitude: upd.Message.Location.Lon,
			Latitude:  upd.Message.Location.Lat,
		}
	}
	return nil
}

func fetchText(upd tgClient.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd tgClient.Update) events.Type {
	if upd.Message.Location != nil {
		return events.Location
	}
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
