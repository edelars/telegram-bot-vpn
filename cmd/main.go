package main

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/controller/configure"
	"backend-vpn/pkg/transport/handlers"
	"backend-vpn/pkg/transport/tgbot"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	appName = "backend-vpn"
	//	fqdnAppName = "tg_bot"
)

func main() {
	var wg sync.WaitGroup
	var err error
	var env config.Environment

	if _, err := flags.Parse(&env); err != nil {
		panic(err)
	}
	logger := config.InitLog(appName, env.LogLevel)

	//ctx := context.Background()
	tgBot := tgbot.NewTgBot(env.TgToken, logger)

	logger.Debug().Msgf("connecting to mysql")
	conStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", env.MysqlEnvironment.DBUser, env.MysqlEnvironment.DBPassword, env.MysqlEnvironment.DBHost, env.MysqlEnvironment.DBName)
	db, err := sqlx.Connect("mysql", conStr)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to connect to mysql")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to properly close mysql conn")
		}
	}()
	logger.Debug().Msgf("connected to mysql")

	errs := make(chan error, 4)
	go waitInterruptSignal(errs)

	ctrl := controller.New()
	logger.Debug().Msg("configure VPN backend")
	if err := configure.MainController(
		ctrl,
		db,
	); err != nil {
		logger.Panic().Err(err).Msg("failed to configure VPN backend")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- tgBot.Listen(handlers.GetHandlers(ctrl, &logger, env))
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
