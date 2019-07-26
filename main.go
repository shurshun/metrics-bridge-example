package main

import (
	"github.com/urfave/cli"
	"metrics/bridge/app"
	"log"
	"os"
)

var (
	Version = "0.1.0"
	cliApp  = cli.NewApp()
)

func init() {
	cliApp.Version = Version
	cliApp.Name = "metrics-bridge"
	cliApp.Usage = "Move metrics from queue to MongoDB"

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "mongo-dsn",
			Value:  "mongodb://test:test@127.0.0.1:27017/admin?connecttimeoutms=1000",
			Usage:  "MongoDB DSN",
			EnvVar: "MONGO_DSN",
		},
		cli.StringFlag{
			Name:   "mongo-db",
			Value:  "days",
			Usage:  "Mongo database",
			EnvVar: "MONGO_DB",
		},
		cli.StringFlag{
			Name:   "queue-dsn",
			Value:  "127.0.0.1:11300",
			Usage:  "Local queue server ([mqtt://]ip:port)",
			EnvVar: "QUEUE_DSN",
		},
		cli.StringFlag{
			Name:   "queue-tube",
			Value:  "metrics",
			Usage:  "Tube with metrics",
			EnvVar: "QUEUE_TUBE",
		},
		cli.IntFlag{
			Name:   "batch-size",
			Value:  1000,
			Usage:  "Batch size",
			EnvVar: "BATCH_SIZE",
		},
		cli.IntFlag{
			Name:   "batch-timeout",
			Value:  30,
			Usage:  "Batch timeout (seconds)",
			EnvVar: "BATCH_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "log-mode",
			Value:  "prod",
			Usage:  "prod/develop",
			EnvVar: "LOG_MODE",
		},
		cli.StringFlag{
			Name:   "metrics-addr",
			Value:  ":8090",
			Usage:  "Metrics addr",
			EnvVar: "METRICS_ADDR",
		},
	}
}

func main() {
	cliApp.Action = app.InitApp
	err := cliApp.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
