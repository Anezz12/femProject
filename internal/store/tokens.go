package store

import (
	"database/sql"
	"time"

	"github.com/Anezz12/femProject/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateToken(UserID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(UserID int, scope string) error
}

func (t *PostgresTokenStore) CreateToken(UserID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(UserID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`
	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(UserID int, scope string) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1 AND scope = $2`
	_, err := t.db.Exec(query, UserID, scope)
	return err
}
