package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Grishun/curate/internal/config"
	log "github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/provider/coindesk"
	"github.com/Grishun/curate/internal/service"
	"github.com/Grishun/curate/internal/storage/memory"
	"github.com/Grishun/curate/internal/transport/http"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "curate",
		Usage: "get rates",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "http-port",
				Value:   "8080",
				Sources: cli.EnvVars("CURATE_HTTP_PORT"),
			},

			&cli.DurationFlag{
				Name:    "provider-interval",
				Value:   time.Minute,
				Sources: cli.EnvVars("CURATE_WORKER_INTERVAL"),
			},

			&cli.StringSliceFlag{
				Name:    "currencies",
				Value:   []string{"BTC", "ETH", "TRX"},
				Sources: cli.EnvVars("CURATE_CURRENCIES"),
			},

			&cli.StringFlag{
				Name:    "quote",
				Value:   "USD",
				Sources: cli.EnvVars("CURATE_QUOTE"),
			},

			&cli.IntFlag{
				Name:    "rate-history-limit",
				Value:   10,
				Sources: cli.EnvVars("RATE_HISTORY_LIMIT"),
			},

			&cli.StringFlag{
				Name:    "coindesk-url",
				Value:   "https://min-api.cryptocompare.com",
				Sources: cli.EnvVars("CURATE_COINDESK_URL"),
			},

			&cli.StringFlag{
				Name:    "coindesk-token",
				Value:   "",
				Sources: cli.EnvVars("CURATE_COINDESK_TOKEN"),
			},
		},

		Action: run,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, c *cli.Command) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg := config.New(c)

	logger := log.NewSlog(log.WithEncoderJSON(slog.LevelDebug))

	provider := coindesk.New(
		coindesk.WithURI(cfg.CoindeskURL),
		coindesk.WithToken(cfg.CoindeskToken),
		coindesk.WithQuote(cfg.Quote),
		coindesk.WithCurrencies(cfg.Currencies...),
		coindesk.WithLogger(logger),
		//TODO : client http
	)

	storage := memory.New()

	svc := service.New(
		service.WithProviders(provider),
		service.WithStorage(storage),
		service.WithLogger(logger),
		service.WithPollingInterval(cfg.WorkerInterval),
		service.WithQuote(cfg.Quote),
	)

	go svc.Start(ctx)

	rest := http.NewRouter(svc)
	errCh := make(chan error, 1)

	go func() {
		logger.Info("starting http server")
		if err := rest.Listen("127.0.0.1:" + cfg.HTTPPort); err != nil {
			errCh <- err
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		sCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := rest.ShutdownWithContext(sCtx); err != nil {
			logger.Error("failed to shutdown http server", "error", err)
		}
		if err := <-errCh; err != nil && ctx.Err() == nil {
			return err
		}
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
