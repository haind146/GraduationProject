package main

import (
	"crypt-coin-payment/app"
	"crypt-coin-payment/controllers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/register", controllers.CreateRegisterUser).Methods("POST")
	router.HandleFunc("/api/user/accept-register", controllers.AcceptRegisterUser).Methods("POST")
	router.HandleFunc("/api/user/import-master-pubkey", controllers.CreateMasterPublicKey).Methods("POST")
	router.HandleFunc("/api/user/application/create", controllers.CreateApplication).Methods("POST")
	router.HandleFunc("/api/user/order/create", controllers.CreateOrder).Methods("POST")
	//router.HandleFunc("/api/user/wallet/change", controllers.AcceptRegisterUser).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":" + port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}