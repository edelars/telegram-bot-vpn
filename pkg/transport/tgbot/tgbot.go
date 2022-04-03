package tgbot

import (
	"backend-vpn/internal/dto"
	"bytes"
	"time"

	"backend-vpn/pkg/transport"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

const (
	timeout     = 10 * time.Second
	buttonInRow = 2
)

type tgBot struct {
	logger zerolog.Logger
	pref   tele.Settings
	bot    *tele.Bot
}

type HandlerTgBot struct {
	Endpoint string
	Handler  func() string
}

func NewTgBot(token string, logger zerolog.Logger) *tgBot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: timeout},
	}
	return &tgBot{
		pref:   pref,
		logger: logger,
	}
}

func (t *tgBot) withHandler(endpoint interface{}, handler func(data transport.HandlerData) interface{}, menu []transport.MenuI) {

	if t == nil {
		return
	}

	if endpoint == "" || handler == nil {
		return
	}

	opts := &tele.ReplyMarkup{ResizeKeyboard: false}

	var buttonsRow []tele.Row
	var buttons tele.Row

	for _, v := range menu {
		btn := opts.Data(v.Data())

		if len(buttons) == buttonInRow {
			buttonsRow = append(buttonsRow, buttons)
			buttons = nil
		}
		buttons = append(buttons, btn)

		t.withMenuHandler(&btn, v)
	}

	if len(buttons) > 0 {
		buttonsRow = append(buttonsRow, buttons)
	}
	opts.Inline(
		buttonsRow...,
	)
	t.bot.Handle(endpoint, func(c tele.Context) error {
		data := transport.HandlerData{
			Username: c.Sender().Username,
			Message:  c.Text(),
			Id:       c.Sender().ID,
		}

		return c.Send(handler(data), opts)
	})

}

func (t *tgBot) withMenuHandler(btn *tele.Btn, v transport.MenuI) {
	t.bot.Handle(btn, func(c tele.Context) error {

		data := transport.HandlerData{
			Username: c.Sender().Username,
			Id:       c.Sender().ID,
		}

		what := v.Handler()(data)

		switch object := what.(type) {
		case string:
			return c.Send(object)
		case dto.TgFile:
			reader := bytes.NewReader(object.File)
			td := tele.Document{File: tele.FromReader(reader), FileName: object.Filename, Caption: object.Caption}
			return c.Send(&td)
		default:
			return c.Send(what)
		}

	})
}

func (t *tgBot) Listen(handlers []transport.HandlerI) (err error) {

	t.bot, err = tele.NewBot(t.pref)
	if err != nil {
		return err
	}

	for _, v := range handlers {
		t.withHandler(v.Endpoint(), v.Handler(), v.Menu())
	}
	t.bot.Start()
	return nil
}

func (t *tgBot) Shutdown() error {
	t.bot.Stop()
	return nil
}

func (t tgBot) Send(tgUserId int64, message string) error {

	rec, err := t.bot.ChatByID(tgUserId)
	if err != nil {
		return err
	}
	_, err = t.bot.Send(rec, message)
	if err != nil {
		return err
	}

	return nil
}
