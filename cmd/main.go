package main

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/controller/configure"
	"backend-vpn/pkg/transport/handlers"
	"backend-vpn/pkg/transport/tgbot"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jessevdk/go-flags"
)

const (
	appName = "tg-bot"
	//	fqdnAppName = "tg_bot"
)

func main() {

	var env config.Environment
	if _, err := flags.Parse(&env); err != nil {
		panic(err)
	}
	logger := config.InitLog(appName, env.LogLevel)

	//ctx := context.Background()
	tgBot := tgbot.NewTgBot(env.TgToken, logger)

	var wg sync.WaitGroup
	var err error

	errs := make(chan error, 4)
	go waitInterruptSignal(errs)

	ctrl := controller.New()
	logger.Debug().Msg("configure VPN backend")
	if err := configure.MainController(
		ctrl,
	); err != nil {
		logger.Panic().Err(err).Msg("failed to configure VPN backend")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- tgBot.Listen(handlers.GetHandlers())
	}()

	logger.Info().Msg("started")
	err = <-errs
	logger.Err(err).Msg("trying to shutdown gracefully")

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tgBot.Shutdown()
		logger.Err(err).Msg("telegram bot stopped")
	}()
}

func waitInterruptSignal(errs chan<- error) {
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
	signal.Stop(c)
}
