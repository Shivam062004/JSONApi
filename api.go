package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
)

func (s *APIserver) run() {
	fmt.Println("http server is statred")

	router := mux.NewRouter()
	router.HandleFunc("/createAccount", createHTTPHandlerFunction(s.handleCreateAccount))
	router.HandleFunc("/admin", createHTTPHandlerFunction(s.handleGetAllAccounts))
	router.HandleFunc("/getAccount/{id}", withJWTauth(s.handleGetAccountByID, s.store))
	router.HandleFunc("/deleteAccount/{id}", withJWTauth(s.handleDeleteAccount, s.store))
	router.HandleFunc("/login", createHTTPHandlerFunction(s.handleLogin))
	http.ListenAndServe(s.listenAddr, router)
}

func (s * APIserver) handleLogin(w http.ResponseWriter, r * http.Request) error {
	loginRequestObject := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(loginRequestObject)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	accountForChecking, err:= s.store.getAccountByID(int(loginRequestObject.ID))
	if err != nil {
		return err
	}
	if accountForChecking.Password != loginRequestObject.Password {
		writeJSON(w, http.StatusForbidden, apiError{
			DefinationOfError: "access denied",
		})
	}

	tokenString, err := createJWTtoken(accountForChecking)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusOK, LoginResponse{
		TokenString: tokenString,
	})
	return nil;
}

func (s *APIserver) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id ,err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return fmt.Errorf("not a valid number")
	}
	acc, err := s.store.getAccountByID(id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, acc)

}

func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accountReqObject := &CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(accountReqObject)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	account := createAccount(accountReqObject.FirstName, accountReqObject.LastName, accountReqObject.Password)
	id, err2 := s.store.storeAccount(account)
	if err2 != nil {
		return err2
	}
	account.ID = id
	return writeJSON(w, http.StatusOK, account)
}

func (s *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id ,err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return fmt.Errorf("not a valid number")
	}
	err2 := s.store.deleteAccountByID(id)
	if err2 != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, map[string]int{"deleted" : id})
}

func (s *APIserver) handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	res, err := s.store.getAccounts()
	if err != nil {
		log.Fatal(err)
	}
	return writeJSON(w, http.StatusOK, res)
}
