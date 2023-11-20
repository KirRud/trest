package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"trest/internal/models"
)

type DataBase struct {
	db *gorm.DB
}

type DBInterface interface {
	TokenRepo
}

func InitDB(conf *models.Config) (*DataBase, error) {
	conn := sqlite.Open(fmt.Sprintf("./db/%s.db", conf.DataBase))
	db, err := gorm.Open(conn)
	if err != nil {
		log.Printf("Failed to initialize GORM with SQLite dialect: %v", err)
		return nil, err
	}

	if !db.Migrator().HasTable(&models.TokenDB{}) {
		db.Migrator().CreateTable(&models.TokenDB{})
	}
	sl := &DataBase{db: db}

	return sl, nil
}
