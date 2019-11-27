package main

import (
	"crypt-coin-payment/app"
	"crypt-coin-payment/blockchain"
	"crypt-coin-payment/controllers"
	"crypt-coin-payment/subscriber"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
	"net/http"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	subscribe := subscriber.SubscriberFactory(1)
	go subscribe.Subscribe()

	cronTab := cron.New()
	cronTab.AddFunc("*/1 * * * *", func() {
		blockchain.ScanBlock(1)
	})
	cronTab.Start()
	//blockchain.ScanBlock(1)

	err := http.ListenAndServe(":" + port, handler) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}