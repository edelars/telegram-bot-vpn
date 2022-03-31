package main

import (
	"backend-vpn/pkg/api/http"
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

	logger.Debug().Msgf("connecting to main mysql")
	conStrMain := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", env.MysqlEnvironment.DBUser, env.MysqlEnvironment.DBPassword, env.MysqlEnvironment.DBHost, env.MysqlEnvironment.DBName)
	mainDb, err := sqlx.Connect("mysql", conStrMain)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to connect to main mysql")
	}
	defer func() {
		if err := mainDb.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to properly close main mysql conn")
		}
	}()
	logger.Debug().Msgf("connected to main mysql")

	logger.Debug().Msgf("connecting to strongswan mysql")
	conStrSw := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", env.MysqlEnvironment.DBUserSW, env.MysqlEnvironment.DBPasswordSW, env.MysqlEnvironment.DBHostSW, env.MysqlEnvironment.DBNameSW)
	swanDb, err := sqlx.Connect("mysql", conStrSw)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to connect to strongswan mysql")
	}
	defer func() {
		if err := mainDb.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to properly close strongswan mysql conn")
		}
	}()
	logger.Debug().Msgf("connected to strongswan mysql")

	errs := make(chan error, 4)
	go waitInterruptSignal(errs)

	ctrl := controller.New()
	logger.Debug().Msg("configure VPN backend")
	if err := configure.MainController(
		ctrl,
		mainDb,
		swanDb,
		env,
		&logger,
	); err != nil {
		logger.Panic().Err(err).Msg("failed to configure VPN backend")
	}

	serverHttp, err := http.NewServer("", env.HttpPort, ctrl)
	if err != nil {
		logger.Panic().Err(err).Msg("unable to create http server")
	}
	logger.Debug().Msgf("New hhtp Server")

	logger.Info().Msg("starting telegram bot #1")
	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- tgBot.Listen(handlers.GetHandlers(ctrl, &logger, env))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- serverHttp.Serve()
	}()

	logger.Info().Msg("started")
	err = <-errs
	logger.Err(err).Msg("trying to shutdown gracefully")

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := serverHttp.Shutdown()
		logger.Err(err).Msg("http server stopped")
	}()

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
