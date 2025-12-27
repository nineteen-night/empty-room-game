package config

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Database             DatabaseConfig     `yaml:"database"`
    Kafka                KafkaConfig        `yaml:"kafka"`
    AuthServiceSettings  ServiceSettings    `yaml:"AuthServiceSettings"`
    GameServiceSettings  ServiceSettings    `yaml:"GameServiceSettings"`
}


type DatabaseConfig struct {
    AuthDB DatabaseInstance `yaml:"auth_db"`
    GameDB DatabaseInstance `yaml:"game_db"`
}

type DatabaseInstance struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    DBName   string `yaml:"name"`
    SSLMode  string `yaml:"ssl_mode"`
}

type KafkaConfig struct {
    Host                     string `yaml:"host"`
    Port                     int    `yaml:"port"`
    PartnershipEventsTopic   string `yaml:"partnership_events_topic"`
    RoomCompletedEventsTopic string `yaml:"room_completed_events_topic"`
}

type ServiceSettings struct {
    HTTPPort   int `yaml:"http_port"`
    GRPCPort   int `yaml:"grpc_port"`
    MinNameLen int `yaml:"min_name_len"`
    MaxNameLen int `yaml:"max_name_len"`
}

func LoadConfig(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
    }

    return &config, nil
}