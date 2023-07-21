package config

import (
	"github.com/spf13/viper"
)

type PostgresDBConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
	Port     string `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

type ServerConfig struct {
	RESTPort     string `mapstructure:"restPort"`
	DataCenterId int    `mapstructure:"dataCenterId"`
	MashineId    int    `mapstructure:"mashineId"`
}
type Config struct {
	PostgresDB PostgresDBConfig `mapstructure:"postgresDB"`
	Server     ServerConfig     `mapstructure:"server"`
}

var (
	vp     *viper.Viper
	config *Config
)

func LoadConfig(env string) (*Config, error) {
	vp = viper.New()

	vp.SetConfigType("json")
	vp.SetConfigName(env)
	vp.AddConfigPath("../config/")
	vp.AddConfigPath("../../config/")
	vp.AddConfigPath("config/")

	err := vp.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return &Config{}, err
	}

	return config, err
}

func GetConfig() *Config {
	return config
}
