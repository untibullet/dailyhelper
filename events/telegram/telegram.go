package telegram

import (
	"errors"

	"github.com/untibullet/dailyhelper/clients/telegram"
	"github.com/untibullet/dailyhelper/events"
	"github.com/untibullet/dailyhelper/storage"
	"github.com/untibullet/dailyhelper/tools/elog"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Processor struct {
	client  *telegram.Client
	offset  int
	storage storage.PageStorer
}

type Meta struct {
	ChatID   int
	Username string
}

func NewProcessor(client *telegram.Client, storage storage.PageStorer) *Processor {
	return &Processor{
		client:  client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.client.GetUpdates(p.offset, limit)
	if err != nil {
		return nil, elog.Wrap("cannot get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, toEvent(u))
		p.offset = u.ID + 1 // change offset by last update ID
	}

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return elog.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return elog.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return elog.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, elog.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func toEvent(update telegram.Update) events.Event {
	uType := fetchType(update)

	res := events.Event{
		Type: uType,
		Text: fetchText(update),
	}

	if uType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}

	return events.Message
}
