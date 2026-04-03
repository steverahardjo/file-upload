package config

import (
	"os"
)

type Config struct {
	DB       DBConfig
	ObjStore ObjStoreConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ObjStoreConfig struct {
	Location   string
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	Secure     bool
}

func LoadConfig() *Config {
	return &Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		ObjStore: ObjStoreConfig{},
	}
}
