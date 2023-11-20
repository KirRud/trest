package models

type Config struct {
	DB
	Secret
}

type DB struct {
	DataBase string
}

type Secret struct {
	SecretKey string
}
