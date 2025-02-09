package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
	// we can also use 
	// bodyToWrite, err := json.Marshal(v)
	// if err != nil {
	// 	return err
	// }
	// _, err2 := w.Write(bodyToWrite)
	// return err2
}


func createHTTPHandlerFunction(handlerFunction apiFunction) http.HandlerFunc {
	fmt.Println("called the wrapper")
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("called the handler")
		err := handlerFunction(w, r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, apiError {
				DefinationOfError: err.Error(),
			})
			return
		}
	}
}
