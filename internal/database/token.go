package database

import (
	"fmt"
	"trest/internal/models"
)

type TokenRepo interface {
	CreateToken(model *models.TokenDB) error
	GetByAccessToken(access string) ([]models.TokenDB, error)
	GetByID(id string) ([]models.TokenDB, error)
	UpdateByRefreshToken(model *models.TokenDB, oldAccessToken string) error
	DeleteByRefreshToken(refreshToken string) error
	DeleteByGUID(guid string) error
}

func (s DataBase) CreateToken(model *models.TokenDB) error {
	if err := s.db.Create(&model).Error; err != nil {
		return err
	}
	return nil
}

func (s DataBase) GetByAccessToken(access string) ([]models.TokenDB, error) {
	var tokens []models.TokenDB
	if err := s.db.Where("access_token = ?", access).Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("error getting token data: %v", err)
	}
	return tokens, nil
}

func (s DataBase) GetByID(id string) ([]models.TokenDB, error) {
	var tokens []models.TokenDB
	if err := s.db.Where("id = ?", id).Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("error getting token data: %v", err)
	}
	return tokens, nil
}

func (s DataBase) UpdateByRefreshToken(model *models.TokenDB, oldAccessToken string) error {
	err := s.db.Model(&models.TokenDB{}).Where("id = ? AND access_token", model.ID, oldAccessToken).
		Updates(models.TokenDB{Token: model.Token, AccessToken: model.AccessToken, UpdatedAt: model.UpdatedAt, ExpireAt: model.ExpireAt}).Error
	if err != nil {
		return fmt.Errorf("error updating app data: %v", err)
	}
	return nil
}

func (s DataBase) DeleteByRefreshToken(refreshToken string) error {
	if err := s.db.Where("token = ?", refreshToken).Delete(&models.TokenDB{}).Error; err != nil {
		return err
	}
	return nil
}

func (s DataBase) DeleteByGUID(guid string) error {
	if err := s.db.Where("id = ?", guid).Delete(&models.TokenDB{}).Error; err != nil {
		return err
	}
	return nil
}
