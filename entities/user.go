package entities

import (
	"crypto/rsa"
)

type User struct {
	ID         int             `json:"-"`
	Email      string          `json:"Email"`
	Password   string          `json:"password"`
	PublicKey  *rsa.PublicKey  `json:"publicKey"`
	PrivateKey *rsa.PrivateKey `json:"-"`
}

func (u *User) PublicKeyMatches(clientPublicKey *rsa.PublicKey) bool {
	// Compare as chaves públicas.
	return u.PublicKey.N.Cmp(clientPublicKey.N) == 0 && u.PublicKey.E == clientPublicKey.E
}
