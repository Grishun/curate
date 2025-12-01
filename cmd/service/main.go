package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	rest "github.com/Grishun/curate/internal/clients/rest"
	cli "github.com/urfave/cli/v3"

	config "github.com/Grishun/curate/internal/config"
	log "github.com/Grishun/curate/internal/log"
	coindesk "github.com/Grishun/curate/internal/provider/coindesk"
	service "github.com/Grishun/curate/internal/service"
	memory "github.com/Grishun/curate/internal/storage/memory"
	http "github.com/Grishun/curate/internal/transport/http"
)

// namedEnv constructs a cli.ValueSourceChain with environment variables prefixed by "CURATE_"
// and appends any given envs.
func namedEnv(envs ...string) cli.ValueSourceChain {
	resultEnvs := cli.EnvVars()

	for _, env := range envs {
		resultEnvs.Append(cli.EnvVars("CURATE_" + env))
	}

	return resultEnvs
}

func main() {
	cmd := &cli.Command{
		Name:        "Curate",
		Description: "Crypto currency rates service",
		Usage:       "Run the service",
		Action:      run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "rest-host", // TODO: need to transfer it to config
				Value:   "127.0.0.1",
				Sources: namedEnv("HTTP_HOST"),
			},
			&cli.StringFlag{
				Name:    "rest-port",
				Value:   "8080",
				Sources: namedEnv("HTTP_PORT"),
			},
			&cli.DurationFlag{
				Name:    "polling-interval",
				Value:   time.Minute,
				Sources: namedEnv("POOLING_INTERVAL"),
			},
			&cli.StringSliceFlag{
				Name:    "currencies",
				Value:   []string{"BTC", "ETH", "TRX"},
				Sources: namedEnv("CURRENCIES"),
			},
			&cli.StringFlag{
				Name:    "quote",
				Value:   "USD",
				Sources: namedEnv("QUOTE"),
			},
			&cli.UintFlag{
				Name:    "history-limit",
				Value:   10,
				Sources: namedEnv("HISTORY_LIMIT"),
			},
			&cli.StringFlag{
				Name:    "coindesk-url",
				Value:   "https://min-api.cryptocompare.com",
				Sources: namedEnv("COINDESK_URL"),
			},
			&cli.StringFlag{
				Name:    "coindesk-token",
				Value:   "",
				Sources: cli.EnvVars("COINDESK_TOKEN"),
			},
		},
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

	httpClient := rest.NewClient(
		rest.WithLogger(logger),
		rest.WithContext(ctx),
		rest.WithTimeout(time.Minute),
	)

	provider := coindesk.New(
		coindesk.WithURI(cfg.CoindeskURL),
		coindesk.WithToken(cfg.CoindeskToken),
		coindesk.WithQuote(cfg.Quote),
		coindesk.WithCurrencies(cfg.Currencies...),
		coindesk.WithLogger(logger),
		coindesk.WithHTTPClient(httpClient),
	)

	storage := memory.New(cfg.HistoryLimit)

	svc := service.New(
		service.WithProviders(provider),
		service.WithStorage(storage),
		service.WithLogger(logger),
		service.WithPollingInterval(cfg.PollingInterval),
		service.WithQuote(cfg.Quote),
	)

	go func() {
		if err := svc.Start(ctx); err != nil {
			logger.Error("service failed to start", "error", err)
		}
	}()

	rest := http.NewRouter(svc)
	errCh := make(chan error, 1)

	go func() {
		addr := net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort)

		logger.Info("starting rest server",
			"addr", addr,
		)

		if err := rest.Listen(addr); err != nil {
			errCh <- err
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		sCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := rest.ShutdownWithContext(sCtx); err != nil {
			logger.Error("failed to shutdown rest server", "error", err)
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
