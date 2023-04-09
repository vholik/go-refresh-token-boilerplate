package user

import (
	"context"
	"database/sql"
)

type repository struct {
	db DBTX
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func NewUserRepository(db DBTX) UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) SignUp(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	query := "INSERT INTO users(name, password, email) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Name, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByCredentials(ctx context.Context, req *SignInReq, password string) (*User, error) {
	var user User
	var userId int64
	query := "SELECT id FROM users WHERE email = $1 AND password = $2"
	err := r.db.QueryRowContext(ctx, query, req.Email, password).Scan(&userId)
	if err != nil {
		return &User{}, err
	}

	user.ID = userId
	return &user, nil
}

func (r *repository) CreateRefreshSession(ctx context.Context, token RefreshSession) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)",
		token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (r *repository) GetRefreshSession(ctx context.Context, token string) (RefreshSession, error) {
	var t RefreshSession
	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	_, err = r.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE user_id=$1", t.UserID)

	return t, err
}
