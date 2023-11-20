package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
	"trest/internal/models"
)

type TokenHandleInterface interface {
	CreateTokenPair(id string) (*models.TokenPair, error)
	RefreshToken(refreshToken, accessToken string) (*models.TokenPair, error)
	RevokeRefreshToken(token, id string) error
	RevokeAllRefreshTokens(id string) error

	generateRefreshToken() string
	storeToken(id string, refreshToken, accessToken string, expireAt time.Time) error
	findToken(refreshToken, accessToken string) (*models.TokenDB, error)
	refreshToken(token *models.TokenDB, refreshToken, accessToken string, expireAt time.Time) error
	getAllTokensByID(id string) ([]models.TokenDB, error)
	deleteToken(refreshTokenHash string) error
	deleteTokens(id string) error
}

func (h Handler) CreateTokenPair(id string) (*models.TokenPair, error) {
	accessClaims := jwt.MapClaims{
		"id": id,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessString, err := accessToken.SignedString([]byte(h.secretKey))
	if err != nil {
		return nil, err
	}

	refreshToken := h.generateRefreshToken()
	refreshExpireAt := time.Now().Add(time.Hour * 24 * 30) // Refresh token expires in 30 days

	err = h.storeToken(id, refreshToken, accessString, refreshExpireAt)
	if err != nil {
		return nil, err
	}

	tokenPair := models.TokenPair{
		AccessToken:  accessString,
		RefreshToken: refreshToken,
	}
	return &tokenPair, nil
}

func (h Handler) RefreshToken(refreshToken, accessToken string) (*models.TokenPair, error) {
	token, err := h.findToken(refreshToken, accessToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid refresh token: %s", err.Error())
	}

	if token.ExpireAt.Before(time.Now()) {
		return nil, fmt.Errorf("Refresh token has expired")
	}

	accessClaims := jwt.MapClaims{
		"id": token.ID,
	}
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessString, err := newAccessToken.SignedString([]byte(h.secretKey))
	if err != nil {
		return nil, err
	}
	newRefreshToken := h.generateRefreshToken()
	refreshExpireAt := time.Now().Add(time.Hour * 24 * 30)

	if err := h.refreshToken(token, newRefreshToken, accessToken, refreshExpireAt); err != nil {
		return nil, err
	}

	return &models.TokenPair{AccessToken: accessString, RefreshToken: refreshToken}, nil
}

func (h Handler) RevokeRefreshToken(refreshToken, id string) error {
	tokens, err := h.getAllTokensByID(id)
	if err != nil {
		return err
	}

	token, err := h.findByRefreshToken(tokens, refreshToken)
	if err != nil {
		return err
	}
	return h.deleteToken(token.Token)
}

func (h Handler) RevokeAllRefreshTokens(id string) error {
	return h.deleteTokens(id)
}

func (h Handler) generateRefreshToken() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes)
}

func (h Handler) storeToken(id string, refreshToken, accessToken string, expireAt time.Time) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	token := models.TokenDB{
		ID:          id,
		Token:       string(hashedToken),
		AccessToken: accessToken,
		ExpireAt:    expireAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return h.db.CreateToken(&token)
}

func (h Handler) findToken(refreshToken, accessToken string) (*models.TokenDB, error) {
	tokens, err := h.db.GetByAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	return h.findByRefreshToken(tokens, refreshToken)
}

func (h Handler) findByRefreshToken(tokens []models.TokenDB, refreshToken string) (*models.TokenDB, error) {
	var token models.TokenDB
	for _, t := range tokens {
		if err := bcrypt.CompareHashAndPassword([]byte(t.Token), []byte(refreshToken)); err == nil {
			token = t
			break
		}
	}

	if token.ID == "" {
		return nil, fmt.Errorf("refresh token not found")
	}
	return &token, nil
}

func (h Handler) refreshToken(token *models.TokenDB, refreshToken, accessToken string, expireAt time.Time) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	oldAccessToken := token.AccessToken
	token.UpdateToken(string(hashedToken), accessToken, expireAt)

	return h.db.UpdateByRefreshToken(token, oldAccessToken)
}

func (h Handler) getAllTokensByID(id string) ([]models.TokenDB, error) {
	return h.db.GetByID(id)
}

func (h Handler) deleteToken(refreshTokenHash string) error {
	return h.db.DeleteByRefreshToken(refreshTokenHash)
}

func (h Handler) deleteTokens(id string) error {
	return h.db.DeleteByGUID(id)
}
