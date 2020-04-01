package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/zjbztianya/poppy/conf"
	"github.com/zjbztianya/poppy/models"
	"github.com/zjbztianya/poppy/router"
	"github.com/zjbztianya/poppy/util/rand"
	"net/http"
)

func main() {
	conf.Init()
	gin.SetMode(conf.Conf.Server.RunMode)
	services, err := models.NewServices()
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	r := router.InitRouter(services)
	authKey, err := rand.Bytes(conf.Conf.App.AuthKeyBytes)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(authKey, csrf.Secure(conf.Conf.App.CsrfSecure))
	endPoint := fmt.Sprintf(":%d", conf.Conf.Server.HttpPort)
	http.ListenAndServe(endPoint, csrfMw(r))
}
