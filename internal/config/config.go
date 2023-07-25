package config

import (
	"github.com/spf13/viper"
)

// PostgresDBConfig its PostgresDB configuration struct that using in Config
type PostgresDBConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
	Port     string `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

// ServerConfig its Server configuration struct that using in Config
type ServerConfig struct {
	RESTPort     string `mapstructure:"restPort"`
	DataCenterID int    `mapstructure:"dataCenterId"`
	MashineID    int    `mapstructure:"mashineId"`
}

// Config its main config struct for viper
type Config struct {
	PostgresDB PostgresDBConfig `mapstructure:"postgresDB"`
	Server     ServerConfig     `mapstructure:"server"`
}

var (
	vp     *viper.Viper
	config *Config
)

// LoadConfig is a init method that find config file and initialize Config struct. Must be called in main.go
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

// GetConfig method provide geting already init config data
func GetConfig() *Config {
	return config
}
