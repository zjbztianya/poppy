package router

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/zjbztianya/poppy/controllers"
	"github.com/zjbztianya/poppy/middleware"
	"github.com/zjbztianya/poppy/models"
	"net/http"
)

func InitRouter(services *models.Services) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.User(services.User))
	r.HTMLRender = multitemplate.NewRenderer()
	staticC := controllers.NewStatic(r)
	usersC := controllers.NewUsers(services.User, r)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	r.StaticFS("/images", http.Dir("./images/"))
	r.GET("/", staticC.Home.HTML)
	r.GET("/contact", staticC.Contact.HTML)
	r.GET("/faq", staticC.FAQ.HTML)
	r.GET("/signup", usersC.NewView.HTML)
	r.POST("/signup", usersC.Create)
	r.GET("/login", usersC.LoginView.HTML)
	r.POST("/login", usersC.Login)
	r.GET("/cookietest", usersC.CookieTest)
	r.POST("/logout", middleware.User(services.User), usersC.Logout)

	gallery := r.Group("/galleries", middleware.RequireUser())
	gallery.GET("/new", galleriesC.New.HTML)
	gallery.POST("/create", galleriesC.Create)
	gallery.GET("/show/:id", galleriesC.Show)
	gallery.GET("/edit/:id", galleriesC.Edit)
	gallery.POST("/update/:id", galleriesC.Update)
	gallery.POST("/delete/:id", galleriesC.Delete)
	gallery.GET("/index", galleriesC.Index)
	gallery.POST("/images/:id", galleriesC.ImageUpload)
	gallery.POST("/delete/:id/images/:filename", galleriesC.ImageDelete)
	return r
}
