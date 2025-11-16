package config

import (
	"github.com/spf13/viper"
)

var Cfg *Conf

type Conf struct {
	ViaCepAPI            string `mapstructure:"VIA_CEP_API"`
	WeatherAPI           string `mapstructure:"WEATHER_API"`
	ServiceBPort         string `mapstructure:"SERVICE_B_PORT"`
	ServiceAPort         string `mapstructure:"SERVICE_A_PORT"`
	ServiceAPI           string `mapstructure:"SERVICE_API"`
	ServiceNameA         string `mapstructure:"SERVICE_A"`
	ServiceNameB         string `mapstructure:"SERVICE_B"`
	OTelExporterEndpoint string `mapstructure:"OTEL_EXPORTER_ZIPKIN_ENDPOINT"`
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
		return err
	}

	Cfg = cfg

	return nil
}
