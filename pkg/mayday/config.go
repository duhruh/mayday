package mayday

import (
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type app struct {
	Environment string `toml:"environment"`
	Name        string `toml:"name"`
}
type log struct {
	Level string `toml:"level"`
}

type grpcConfig struct {
	Port string `toml:"port"`
}

type config struct {
	GRPC grpcConfig `toml:"grpc"`
	Log  log        `toml:"log"`
	App  app        `toml:"app"`
}

// Config -
type Config interface {
	GRPCPort() string

	LogLevel() logrus.Level

	AppEnvironment() string
	IsProduction() bool
	AppName() string
}

const (
	defaultLogLevel       = "debug"
	productionEnvironment = "production"
)

// NewConfig - generates a config from a toml file.
func NewConfig(file io.Reader) (Config, error) {
	var c config
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return c, err
	}

	if _, err := toml.Decode(string(bytes), &c); err != nil {
		return c, err
	}

	return c, nil
}

func (c config) GRPCPort() string {
	return c.GRPC.Port
}

func (c config) LogLevel() logrus.Level {
	if c.Log.Level == defaultLogLevel {
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}
func (c config) AppEnvironment() string {
	return c.App.Environment
}

func (c config) IsProduction() bool {
	return c.AppEnvironment() == productionEnvironment
}

func (c config) AppName() string {
	return c.App.Name
}
