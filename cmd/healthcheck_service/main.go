package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/Quantum12k/healthcheck-service/internal/config"
	"github.com/Quantum12k/healthcheck-service/internal/logger"
)

const (
	SettingsFlagName = "settings"
	URLsFlagName     = "urls"
	ModeFlagName     = "mode"

	DefaultSettingsConfigName = "../configs/settings.yml"
	DefaultURLsConfigName     = "../configs/urls.yml"

	SingleMode = "single"
	WithDBMode = "withDB"
	ServerMode = "server"
	APIMode    = "api"
)

type (
	App struct {
		Log *zap.SugaredLogger
		Cfg *config.Config
		cli *cli.App
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			cancel()
		}
	}()

	app := App{}

	appCli := &cli.App{
		Name:  "HealthCheck-service",
		Usage: "Сервис для проверки работоспособности сайтов",
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Запускает приложение",
				Action: app.run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  SettingsFlagName,
						Usage: "особый путь до файла настроек",
						Value: DefaultSettingsConfigName,
					},
					&cli.StringFlag{
						Name:  URLsFlagName,
						Usage: "особый путь до конфигурационного файла URL",
						Value: DefaultURLsConfigName,
					},
					&cli.StringFlag{
						Name:  ModeFlagName,
						Usage: "режим работы утилиты",
						Value: SingleMode,
					},
				},
			},
		},
	}

	app.cli = appCli

	if err := app.cli.RunContext(ctx, os.Args); err != nil {
		println("error run app: ", err.Error())
		os.Exit(1)
	}
}

func (a *App) run(cliCtx *cli.Context) error {
	cfg, err := config.New(
		cliCtx.String(SettingsFlagName),
		cliCtx.String(URLsFlagName),
	)
	if err != nil {
		return fmt.Errorf("new config: %v", err)
	}

	a.Cfg = cfg
	a.Log = logger.New(cfg.Logger)

	mode := cliCtx.String(ModeFlagName)

	ctx := cliCtx.Context

	switch mode {
	case SingleMode:
		return a.single(ctx)
	case WithDBMode:
		return a.withDB(ctx)
	case ServerMode:
		return a.server(ctx)
	case APIMode:
		return a.api(ctx)
	default:
		a.Log.Infof("no valid mode provided, got %s, exiting", mode)
		return nil
	}
}
