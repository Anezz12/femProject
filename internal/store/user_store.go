package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
)

type PasswordHash struct {
	plaintText string
	hash       []byte
}

func (p *PasswordHash) SetPassword(plaintTextpassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintTextpassword), 12)
	if err != nil {
		return err
	}
	p.hash = hash
	return nil
}

func (p *PasswordHash) Matches(plaintTextpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintTextpassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type User struct {
	ID           int          `json:"id"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	PasswordHash PasswordHash `json:"-"`
	Bio          string       `json:"bio"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type postgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *postgresUserStore {
	return &postgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) error
	// GetUserByID(id int64) (*User, error)
	// GetUserByEmail(email string) (*User, error)
	// UpdateUser(*User) error
	// DeleteUser(id int64) error
}

func (s *postgresUserStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, bio)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *postgresUserStore) GetUserByName(username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, bio, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &User{
		PasswordHash: PasswordHash{},
	}
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *postgresUserStore) UpdateUser(user *User) error {
	// bisa juuga menggunakan current_timestamp
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, bio = $4, updated_at = NOW()
		WHERE id = $5
	`

	result, err := s.db.Exec(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
