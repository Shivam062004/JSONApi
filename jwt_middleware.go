package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"os"
	jwt "github.com/golang-jwt/jwt/v5"
)

var SECRET string = os.Getenv("JWT_SECRET")

func createJWTtoken(account *Account) (string, error) {
	claims := CustomClaims{
		account.Number,
		account.Password,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			ID:        strconv.Itoa(int(account.ID)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRET))
	fmt.Println(tokenString)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWTtoken(tokenString string, store Storage) bool {
	claims := &CustomClaims{}
	fmt.Println(tokenString)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})
	if err != nil {
		fmt.Println(err)
		return false
	}
	extractedClaims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return false
	}
	id, err := strconv.Atoi(extractedClaims.ID)
	if err != nil {
		return false
	}
	account, err := store.getAccountByID(id)
	if err != nil {
		return false
	}
	if !(extractedClaims.Number == account.Number && extractedClaims.Password == account.Password) {
		return false
	}
	return true
}

func withJWTauth(handlerFunction apiFunction, store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("jwt middleware called")
		tokenString := r.Header.Get("Authorization")
		flag := validateJWTtoken(tokenString, store)
		if flag {
			err := handlerFunction(w, r)
			if err != nil {
				writeJSON(w, http.StatusForbidden, apiError{
					DefinationOfError: err.Error(),
				})
			}
		} else {
			writeJSON(w, http.StatusForbidden, apiError{
				DefinationOfError: "access denied",
			})
		}
	}
}
