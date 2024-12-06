package config

import (
    "log"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
}

type JWTConfig struct {
    Secret string
}

func LoadConfig() Config {
    viper.SetConfigFile("configs/config.yaml")
    viper.AutomaticEnv()

    err := viper.ReadInConfig()
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    var config Config
    err = viper.Unmarshal(&config)
    if err != nil {
        log.Fatalf("Unable to decode into struct: %v", err)
    }

    return config
}
