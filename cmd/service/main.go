package main

import (
	"context"
	"github.com/Grishun/curate/internal/config"
	"log"
	"os"
	"time"

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
				Name:    "worker-interval",
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
		log.Fatal(err)
	}
}

func run(ctx context.Context, c *cli.Command) error {
	_ = config.New(c)

	//x := 12345e-3
	//fmt.Printf("%.20f\nf", x)
	//fmt.Print(reflect.ValueOf(x).Kind())

	return nil
}
