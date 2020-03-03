package main

import (
	"github.com/gorilla/mux"
	"github.com/zjbztianya/poppy/controllers"
	"net/http"
)

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":8080", r)
}
