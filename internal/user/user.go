package user

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID       int64  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Name     string `json:"name" db:"name"`
}

type SignUpReq struct {
	Name     string `json:"name" binding:"required,gte=2,lte=50"`
	Email    string `json:"email" binding:"required,gte=2,lte=50"`
	Password string `json:"password" binding:"required,gte=8,lte=25"`
}

type RefreshSession struct {
	ID        int64
	UserID    int64
	ExpiresAt time.Time
	Token     string
}

type SignUpRes struct {
	ID    string `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
}

type SignInReq struct {
	Email    string `json:"email" binding:"required,gte=2,lte=50"`
	Password string `json:"password" binding:"required,gte=8,lte=25"`
}

var ErrUserNotFound = errors.New("user with such credentials not found")

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	SignUp(ctx context.Context, user *User) (*User, error)
	GetUserByCredentials(ctx context.Context, req *SignInReq, password string) (*User, error)
	CreateRefreshSession(ctx context.Context, token RefreshSession) error
	GetRefreshSession(ctx context.Context, token string) (RefreshSession, error)
}

type UserService interface {
	SignUp(c context.Context, req *SignUpReq) (*SignUpRes, error)
	SignIn(c context.Context, req *SignInReq) (string, string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}
