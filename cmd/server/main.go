package main

import (
	"context"
	"flag"
	"os"

	"github.com/docker/mayday/pkg/db"
	"github.com/docker/mayday/pkg/mayday"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const (
	defaultAppConfig = "configs/mayday.toml"
)

var (
	appConfig = flag.String("config", defaultAppConfig, "application config file")
)

func main() {

	flag.Parse()

	file, err := os.Open(*appConfig)
	if err != nil {
		panic(err)
	}

	cfg, err := mayday.NewConfig(file)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()

	logger.SetLevel(cfg.LogLevel())

	baseLogger := logger.WithFields(logrus.Fields{
		"environment": cfg.AppEnvironment(),
	})

	databaseConnection := db.NewConnection(cfg.DatabaseAdapter(), cfg.DatabaseDSN())

	databaseConnection.Open()

	maydayServer := mayday.NewServer(cfg, baseLogger, databaseConnection)

	baseLogger.Info("application starting")
	if err := maydayServer.Start(context.TODO()); err != nil {
		baseLogger.Errorf("failed to serve: %v", err)
	}
}
