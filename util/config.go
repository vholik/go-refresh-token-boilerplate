package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUrl     string `mapstructure:"DB_URL"`
	JwtSecret string `mapstructure:"JWT_SECRET"`
	TokenTtl  int    `mapstructure:"TOKEN_TTL"`
}

func LoadConfig() (c Config, err error) {
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
