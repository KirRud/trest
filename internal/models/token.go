package models

import (
	"time"
)

type TokenDB struct {
	ID          string    `gorm:"column:id"`
	Token       string    `gorm:"column:token"`
	AccessToken string    `gorm:"column:access_token"`
	ExpireAt    time.Time `gorm:"column:expire_at"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `grom:"column:updated_at"`
}

func (t *TokenDB) UpdateToken(hashedToken, accessToken string, expireAt time.Time) {
	t.AccessToken = accessToken
	t.Token = hashedToken
	t.ExpireAt = expireAt
	t.UpdatedAt = time.Now()
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
