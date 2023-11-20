package controller

import (
	"github.com/gin-gonic/gin"
	"trest/internal/database"
)

type Handler struct {
	db        database.DBInterface
	secretKey string
}

func NewHandler(db database.DBInterface, key string) *Handler {
	return &Handler{
		db:        db,
		secretKey: key,
	}
}

func InitRoutes(handler handlerInterface) *gin.Engine {
	router := gin.Default()
	router.GET("/token/:id", handler.getTokens)
	router.GET("/refresh/:token/:access_token", handler.refreshTokens)
	router.DELETE("/token/:token/:id", handler.revokeRefreshToken)
	router.DELETE("/tokens/:id", handler.revokeAllRefreshTokens)
	return router
}
