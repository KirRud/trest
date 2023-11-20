package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type handlerInterface interface {
	getTokens(c *gin.Context)
	refreshTokens(c *gin.Context)
	revokeRefreshToken(c *gin.Context)
	revokeAllRefreshTokens(c *gin.Context)
}

func (h Handler) getTokens(c *gin.Context) {
	id := c.Params.ByName("id")

	tokenPair, err := h.CreateTokenPair(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken})
}

func (h Handler) refreshTokens(c *gin.Context) {
	refreshToken := c.Params.ByName("token")
	accessToken := c.Params.ByName("access_token")

	tokenPair, err := h.RefreshToken(refreshToken, accessToken)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken})
}

func (h Handler) revokeRefreshToken(c *gin.Context) {
	refreshToken := c.Params.ByName("token")
	id := c.Params.ByName("id")

	if err := h.RevokeRefreshToken(refreshToken, id); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "revoke token"})
}

func (h Handler) revokeAllRefreshTokens(c *gin.Context) {
	id := c.Params.ByName("id")

	if err := h.RevokeAllRefreshTokens(id); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "revoke all tokens"})

}
