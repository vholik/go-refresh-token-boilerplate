package util

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type middleware struct {
	secret string
}

func NewAuthMiddleware(secret string) *middleware {
	return &middleware{
		secret: secret,
	}
}

func (s *middleware) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := s.getTokenFromRequest(c.Request)
		if err != nil {
			respondWithError(c, 401, "API token required")
			return
		}

		userId, err := s.parseToken(token)
		if err != nil {
			respondWithError(c, 401, "Can not authorize user")
			return
		}

		ctx := context.WithValue(c.Request.Context(), userId, userId)
		c.Request.WithContext(ctx)

		c.Next()
	}
}

func (s *middleware) getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("empty auth header")
	}

	return headerParts[1], nil
}

func (s *middleware) parseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return 0, errors.New("invalid token")
		}

		subject, ok := claims["sub"].(string)
		if !ok {
			return 0, errors.New("invalid subject")
		}

		id, err := strconv.Atoi(subject)
		if err != nil {
			return 0, errors.New("invalid subject")
		}

		return int64(id), nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return 0, errors.New("invalid token")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		return 0, errors.New("token expired")
	} else {
		return 0, errors.New("error occured")
	}
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}
