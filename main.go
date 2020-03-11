package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zjbztianya/poppy/controllers"
	"github.com/zjbztianya/poppy/models"
	"net/http"
)

const (
	user     = "root"
	password = "123456"
	dbname   = "test"
)

func main() {
	sqlInfo := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", user, password, dbname)
	us, err := models.NewUserService(sqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("Get")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	http.ListenAndServe(":8080", r)
}
