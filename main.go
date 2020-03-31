package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/zjbztianya/poppy/models"
	"github.com/zjbztianya/poppy/router"
	"github.com/zjbztianya/poppy/util/rand"
	"net/http"
)

const (
	user     = "root"
	password = "123456"
	dbname   = "test"
)

func main() {
	sqlInfo := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", user, password, dbname)
	services, err := models.NewServices(sqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	r := router.InitRouter(services)
	authKey, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(authKey, csrf.Secure(false))
	http.ListenAndServe(":8080", csrfMw(r))
}
