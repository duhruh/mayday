package main

import (
	"flag"
	"os"

	"github.com/docker/mayday/pkg/mayday"
	"github.com/sirupsen/logrus"
)

const (
	defaultAppConfig = "mayday.toml"
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

	maydayServer := mayday.NewServer(cfg, baseLogger)

	baseLogger.Info("application starting")
	if err := maydayServer.Start(); err != nil {
		baseLogger.Errorf("failed to serve: %v", err)
	}
}
