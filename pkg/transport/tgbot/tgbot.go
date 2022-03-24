package tgbot

import (
	"time"

	"backend-vpn/pkg/transport"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

const (
	timeout = 10 * time.Second
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

func (t *tgBot) withHandler(endpoint interface{}, handler func() interface{}, menu []transport.MenuI) {

	if t == nil {
		return
	}

	if endpoint == "" || handler == nil {
		return
	}

	opts := &tele.ReplyMarkup{ResizeKeyboard: true}
	for _, v := range menu {
		opts.Inline(
			opts.Row(opts.Data(v.Data())),
		)
	}

	t.bot.Handle(endpoint, func(c tele.Context) error {
		return c.Send(handler(), opts)
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
