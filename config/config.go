package config

import (
	"encoding/json"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or env variables.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	RedisHost            string        `mapstructure:"REDIS_HOST"`
	RedisPort            string        `mapstructure:"REDIS_PORT"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	StrCorsOrigins       string        `mapstructure:"CORS_ORIGIN"`
	CorsOrigins          []string
	Auth0Domain          string `mapstructure:"AUTH0_DOMAIN"`
	Auth0ClientID        string `mapstructure:"AUTH0_CLIENT_ID"`
	Auth0ClientSecret    string `mapstructure:"AUTH0_CLIENT_SECRET"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv() // overwrites config file if variables are specified in current env.

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	json.Unmarshal([]byte(config.StrCorsOrigins), &(config.CorsOrigins))

	return
}
