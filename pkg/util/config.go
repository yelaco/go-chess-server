package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host            string        `mapstructure:"address.host"`
	Port            string        `mapstructure:"address.port"`
	MatchingTimeout time.Duration `mapstructure:"game.matching_timeout"`
	DBName          string        `mapstructure:"database.name"`
	DBHost          string        `mapstructure:"database.host"`
	DBUser          string        `mapstructure:"database.user"`
	DBPassword      string        `mapstructure:"database.password"`
}

// LoadConfig reads configurations from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("server")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
