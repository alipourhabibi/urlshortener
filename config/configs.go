package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Confs = Config{}

type Config struct {
	Port       string  `mapstructure:"Port"`
	PostgresDB Posgres `mapstructure:"PostgresDB"`
	Redis      Redis   `mapstructure:"RedisDB"`
	RedisJWT   Redis   `mapstructure:"RedisJWT"`
	Auth       Auth    `mapstructure:"Auth"`
}

func (g *Config) Load(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		return viper.Unmarshal(&Confs)
	}
	return fmt.Errorf("file not exists")

}
