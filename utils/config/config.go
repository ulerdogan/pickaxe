package config_utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBUsername    string `mapstructure:"DB_USERNAME"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	MigrationURL  string `mapstructure:"MIGRATION_URL"`
	RPCAddress    string `mapstructure:"RPC_ADDRESS"`
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
