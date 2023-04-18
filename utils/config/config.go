package config_utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	AuthPassword        string        `mapstructure:"AUTH_PASSWORD"`
	SymmetricKey        string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RPCAddress          string        `mapstructure:"RPC_ADDRESS"`
	MigrationURL        string        `mapstructure:"MIGRATION_URL"`
}

func LoadConfig(name, path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
