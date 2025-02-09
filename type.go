package main

import (
	"math/rand/v2"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"
)

type apiFunction func(http.ResponseWriter, *http.Request) error

type apiError struct {
	DefinationOfError string `json:"definationOfError"`
}

type APIserver struct {
	listenAddr string
	store      Storage
}

type Account struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
	Password  string `json:"-"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	TokenString string `json:"tokenstring"`
}

type CustomClaims struct {
	Number   int64
	Password string
	jwt.RegisteredClaims
}

func createAccount(firstname string, lastname string, password string) *Account {
	return &Account{
		FirstName: firstname,
		LastName:  lastname,
		Number:    int64(rand.IntN(545535243)),
		Password:  password,
	}
}

func makeServer(addr string, store Storage) *APIserver {
	return &APIserver{
		listenAddr: addr,
		store:      store,
	}
}
