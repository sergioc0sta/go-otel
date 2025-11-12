package config

import (
	"github.com/spf13/viper"
)

var Cfg *Conf

type Conf struct {
	ViaCepAPI  string `mapstructure:"VIA_CEP_API"`
	WeatherAPI string `mapstructure:"WEATHER_API"`
}

func LoadConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	cfg := &Conf{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil
	}

	Cfg = cfg

	return nil
}
