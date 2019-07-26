package app

import (
	"github.com/urfave/cli"
	"metrics/bridge/app/input"
	"metrics/bridge/app/logger"
	"metrics/bridge/app/processor"
	"metrics/bridge/app/storage"
)

func InitApp(cliContext *cli.Context) error {
	log, _ := logger.New(cliContext.String("log-mode"))

	queue, err := input.New(
		cliContext.String("queue-dsn"),
		cliContext.String("queue-tube"),
		log)

	if err != nil {
		log.Fatalf("can't initialize queue: %v", err)
	}

	defer queue.Disconnect()

    log.Info("Successfully connected to queue")

	mongo, err := storage.New(
		cliContext.String("mongo-dsn"),
		cliContext.String("mongo-db"))

	if err != nil {
		log.Fatalf("can't initialize storage: %v", err)
	}

	defer mongo.Disconnect()

    log.Info("Successfully connected to storage")

    go initMetrics(cliContext.String("metrics-addr"))

    log.Info("Successfully initialized pprof & metrics on port " + cliContext.String("metrics-port"))

	log.Info("Starting metrics processor...")

	p := processor.New(
		queue,
		mongo,
		cliContext.Int("batch-size"),
		cliContext.Int("batch-timeout"),
		log)

	p.Run()

	return nil
}
