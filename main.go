package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meinantoyuriawan/gobank-v2/auth"
	"github.com/meinantoyuriawan/gobank-v2/controller/accountcontroller"
	"github.com/meinantoyuriawan/gobank-v2/controller/transactioncontroller"
	"github.com/meinantoyuriawan/gobank-v2/models"
)

func main() {
	models.ConnectDB()
	r := mux.NewRouter()

	r.HandleFunc("/login", accountcontroller.Login).Methods("POST")
	r.HandleFunc("/register", accountcontroller.Registeruser).Methods("POST")
	r.HandleFunc("/logout", accountcontroller.Logout).Methods("GET")
	authwall := r.PathPrefix("/authwall").Subrouter()
	authwall.HandleFunc("/create", accountcontroller.Createaccount).Methods("POST")
	authwall.Use(auth.ValidateJWT)

	r.HandleFunc("/saldo", transactioncontroller.Getsaldo).Methods("GET")
	r.HandleFunc("/transfer", transactioncontroller.Transfer).Methods("POST")
	r.HandleFunc("/recent", transactioncontroller.Getrecent).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
