package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	postgres "github.com/cassiusbessa/satback/db"
	"github.com/cassiusbessa/satback/entities"
)

// User representa a estrutura de dados do usuário.

// Token representa a estrutura de dados do token de autenticação.
type Token struct {
	UserID    int       `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
}

var db, _ = postgres.Connect()
var userRepo = postgres.NewUserRepository(db)

// signupHandler lida com as requisições de signup.
func signupHandler(w http.ResponseWriter, r *http.Request) {
	var newUser entities.User

	// Decodifica o JSON do corpo da requisição para obter os detalhes do usuário.
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Verifica se o usuário já existe no banco de dados.
	_, err = userRepo.GetUserByEmail(newUser.Email)
	if err == nil {
		http.Error(w, "Usuário já existe", http.StatusConflict)
		return
	}

	if len(newUser.Password) < 6 {
		http.Error(w, "Senha deve ter 6 dígitos", http.StatusBadRequest)
		return
	}

	// Gera um par de chaves para o novo usuário.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		http.Error(w, "Erro ao gerar chave privada", http.StatusInternalServerError)
		return
	}

	// Armazena a chave privada no objeto do usuário (normalmente, seria armazenada de forma mais segura).
	newUser.PrivateKey = privateKey

	// Converte a chave pública para o formato PKCS#1
	publicKeyPKCS1 := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	// Armazena o usuário no db.
	err = userRepo.CreateUser(newUser.Email, newUser.Password, publicKeyPKCS1)
	if err != nil {
		println(err.Error())
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	// Retorna a chave pública e privada ao cliente.
	response := struct {
		PublicKey  string `json:"publicKey"`
		PrivateKey string `json:"privateKey"`
	}{
		PublicKey:  base64.StdEncoding.EncodeToString(publicKeyPKCS1),
		PrivateKey: base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey)),
	}

	// Codifica a resposta em JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		PublicKey string `json:"publicKey"`
	}

	// Decodifica o JSON do corpo da requisição para obter as credenciais do usuário.
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Recupera o usuário do banco de dados usando o email.
	user, err := userRepo.GetUserByEmail(credentials.Email)
	if err != nil {
		println(err.Error())
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	// Converte a chave pública do cliente de base64 para *rsa.PublicKey.
	clientPublicKeyBytes, err := base64.StdEncoding.DecodeString(credentials.PublicKey)
	if err != nil {
		http.Error(w, "Erro ao decodificar chave pública", http.StatusBadRequest)
		return
	}

	clientPublicKey, err := x509.ParsePKCS1PublicKey(clientPublicKeyBytes)
	if err != nil {
		http.Error(w, "Erro ao converter chave pública", http.StatusBadRequest)
		return
	}

	// Verifica se a chave pública do cliente corresponde à chave pública do usuário no banco de dados.
	if !user.PublicKeyMatches(clientPublicKey) {
		http.Error(w, "Chave pública inválida", http.StatusUnauthorized)
		return
	}

	// As credenciais são válidas. Gere um token de autenticação.
	token := Token{
		UserID:    user.ID,
		Timestamp: time.Now(),
	}

	// Codifica o token em JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
