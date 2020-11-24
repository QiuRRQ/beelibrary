package main

import (
	"city/controllers"
	d "city/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	db := d.GetDB()

	db.AutoMigrate(&d.Users{}, &d.Stock{}, &d.Borrowd{}, &d.Borrow{}, &d.Books{})
	db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/city/{id}", controllers.GetCity).Methods("GET")
	router.HandleFunc("/api/city", controllers.GetCities).Methods("GET")
	router.HandleFunc("/api/srep/{id}", controllers.GetSrep).Methods("GET")
	router.HandleFunc("/api/srep", controllers.GetSreps).Methods("GET")
	router.HandleFunc("/api/city", controllers.CreateCity).Methods("POST")
	router.HandleFunc("/api/srep", controllers.CreateSrep).Methods("POST")

	router.HandleFunc("/api/newestbook", controllers.GetNewestBook).Methods("GET")
	router.HandleFunc("/api/book/{id}", controllers.GetBookByID).Methods("GET")
	router.HandleFunc("/api/book", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/api/book/{id}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/api/borrow/{id}", controllers.BorrowDetail).Methods("GET")
	router.HandleFunc("/api/borrow", controllers.BorrowC).Methods("POST")
	router.HandleFunc("/api/borrow", controllers.Borrowing).Methods("GET")
	router.HandleFunc("/api/user/login", controllers.UserLoginController).Methods("POST")
	router.HandleFunc("/api/user", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/return/{id}", controllers.ReturningC).Methods("POST")
	// router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}

}
