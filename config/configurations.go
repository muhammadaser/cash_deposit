package config

import (
	"io"
	"os"

	"github.com/spf13/viper"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	// Env s
	Env string = "dev"
)

// Configuration s
type Configuration struct {
	Address string
	Log     Log
	Pg      Pg
}

// Log s
type Log struct {
	FileName   string
	MaxSize    int
	MaxBackups int
}

// Pg s
type Pg struct {
	Addr     string
	Port     string
	Username string
	Password string
	Database string
}

var (
	// Config s
	Config Configuration
	// LogOutput k
	LogOutput io.Writer
)

// InitConfig k
func InitConfig(env string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// unmarshal config ke dalam struct
	err = viper.UnmarshalKey(env, &Config)

	// set Output log, untuk production set ke dalam file.
	if env == "prod" {
		LogOutput = &lumberjack.Logger{
			Filename:   Config.Log.FileName,
			MaxSize:    Config.Log.MaxSize,
			MaxBackups: Config.Log.MaxBackups,
		}
	} else {
		LogOutput = os.Stderr
	}

	return err
}
