package config

import (
	"io"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// PostgresDBConfig its PostgresDB configuration struct that using in Config
type PostgresDBConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	DBName   string `mapstructure:"dbName" yaml:"dbName"`
	Port     string `mapstructure:"port" yaml:"port"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"`
}

// ServerConfig its Server configuration struct that using in Config
type ServerConfig struct {
	RESTPort     string `mapstructure:"restPort" yaml:"restPort"`
	DataCenterID int    `mapstructure:"dataCenterId" yaml:"dataCenterId"`
	MashineID    int    `mapstructure:"mashineId" yaml:"mashineId"`
}

// Config its main config struct for viper
type Config struct {
	PostgresDB PostgresDBConfig `mapstructure:"postgresDB" yaml:"postgresDB"`
	Server     ServerConfig     `mapstructure:"server" yaml:"server"`
}

var (
	vp     *viper.Viper
	config *Config
)

// LoadConfigJSON is a init method that find config json file and initialize Config struct. Must be called in main.go. Using in local env
func LoadConfigJSON(env string) (*Config, error) {
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

// LoadConfigYAML is a init method that find config yaml file and initialize Config struct. Using when deploying on prod
func LoadConfigYAML() (*Config, error) {
	file, err := os.Open("config/config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Parse the YAML data into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, err
}

// GetConfig method provide geting already init config data
func GetConfig() *Config {
	return config
}
