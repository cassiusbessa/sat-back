package postgres

import (
	"crypto/x509"
	"database/sql"
	"sync"

	"github.com/cassiusbessa/satback/entities"
)

type UserRepository struct {
	db *sql.DB
}

var (
	userRepository *UserRepository
	userRepoOnce   sync.Once
)

func NewUserRepository(db *sql.DB) *UserRepository {
	userRepoOnce.Do(func() {
		userRepository = &UserRepository{
			db: db,
		}
	})
	return userRepository
}

func (r *UserRepository) CreateUser(email string, password string, publicKey []byte) error {
	_, err := r.db.Exec(`
		INSERT INTO users (email, password, public_key) VALUES ($1, $2, $3)
	`, email, password, publicKey)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	var publicKey []byte

	err := r.db.QueryRow(`
		SELECT id, email, password, public_key FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.Password, &publicKey)
	if err != nil {
		return nil, err
	}

	// Converte a []byte publicKey para *rsa.PublicKey, se necessário.
	user.PublicKey, err = x509.ParsePKCS1PublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
