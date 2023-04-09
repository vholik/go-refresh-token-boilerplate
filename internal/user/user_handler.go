package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService UserService
}

func NewHandler(s UserService) *Handler {
	return &Handler{
		UserService: s,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	var req SignUpReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.UserService.SignUp(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) SignIn(c *gin.Context) {
	var req SignInReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.UserService.SignIn(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))

	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *Handler) Refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.UserService.RefreshTokens(c, cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))

	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}
