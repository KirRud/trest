package config

import (
	"github.com/joho/godotenv"
	"os"
	"trest/internal/models"
)

func InitConfig() (*models.Config, error) {
	var cfg models.Config
	err := godotenv.Load("./configs/dev.env")
	if err != nil {
		return nil, err
	}

	cfg.DB = models.DB{
		DataBase: os.Getenv("DB_NAME"),
	}
	cfg.Secret = models.Secret{
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	return &cfg, nil
}
