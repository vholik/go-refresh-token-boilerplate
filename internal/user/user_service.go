package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	UserRepository
	hasher     PasswordHasher
	hmacSecret []byte
	tokenTtl   int
}

func NewService(repository UserRepository, hasher PasswordHasher, secret []byte, tokenTtl int) UserService {
	return &service{
		repository,
		hasher,
		secret,
		tokenTtl,
	}
}

func (s *service) SignUp(c context.Context, req *SignUpReq) (*SignUpRes, error) {
	password, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	r, err := s.UserRepository.SignUp(c, u)
	if err != nil {
		return nil, err
	}

	res := &SignUpRes{
		ID:    strconv.Itoa(int(r.ID)),
		Name:  r.Name,
		Email: r.Email,
	}

	return res, nil
}

func (s *service) SignIn(c context.Context, req *SignInReq) (string, string, error) {
	password, err := s.hasher.Hash(req.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.UserRepository.GetUserByCredentials(c, req, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrUserNotFound
		}

		return "", "", err
	}

	return s.generateTokens(c, user.ID)
}

func (s *service) generateTokens(ctx context.Context, userId int64) (string, string, error) {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * time.Duration(s.tokenTtl))},
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	})

	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.UserRepository.CreateRefreshSession(ctx, RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.UserRepository.GetRefreshSession(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("session expired")
	}

	return s.generateTokens(ctx, session.UserID)
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
