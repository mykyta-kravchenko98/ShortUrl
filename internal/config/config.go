package config

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/mykyta-kravchenko98/ShortUrl/pkg/closeutil"
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

	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.AutomaticEnv()

	for _, key := range []string{
		"postgresDB.host", "postgresDB.user", "postgresDB.password",
		"postgresDB.dbName", "postgresDB.port", "postgresDB.sslmode",
		"server.restPort", "server.dataCenterId", "server.mashineId",
	} {
		_ = vp.BindEnv(key)
	}

	if err := vp.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return &Config{}, err
		}
	}

	if err := vp.Unmarshal(&config); err != nil {
		return &Config{}, err
	}

	return config, nil
}

// LoadConfigYAML is a init method that find config yaml file and initialize Config struct. Using when deploying on prod
func LoadConfigYAML() (*Config, error) {
	file, err := os.Open("config/config.yml")
	if err != nil {
		return nil, err
	}
	defer closeutil.Close(file)

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

// GetConfig method provide getting already init config data
func GetConfig() *Config {
	return config
}
