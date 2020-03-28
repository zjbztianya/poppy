package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/zjbztianya/poppy/controllers"
	"github.com/zjbztianya/poppy/middleware"
	"github.com/zjbztianya/poppy/models"
	"github.com/zjbztianya/poppy/rand"
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

	r := mux.NewRouter()
	userMw := middleware.User{UserService: services.User}
	requireUserMw := middleware.RequireUser{}
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")

	logoutUser := requireUserMw.ApplyFn(usersC.Logout)
	r.Handle("/signup", usersC.NewView).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.Handle("/logout", logoutUser).Methods("POST")

	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)
	editGallery := requireUserMw.ApplyFn(galleriesC.Edit)
	updateGallery := requireUserMw.ApplyFn(galleriesC.Update)
	deleteGallery := requireUserMw.ApplyFn(galleriesC.Delete)
	indexGallery := requireUserMw.ApplyFn(galleriesC.Index)
	deleteImage := requireUserMw.ApplyFn(galleriesC.ImageDelete)
	imageUploadGallery := requireUserMw.ApplyFn(galleriesC.ImageUpload)
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}",
		galleriesC.Show).Methods("GET").
		Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET").
		Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", deleteGallery).Methods("POST")
	r.HandleFunc("/galleries", indexGallery).Methods("GET").
		Name(controllers.IndexGalleries)
	r.HandleFunc("/galleries/{id:[0-9]+}/images", imageUploadGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", deleteImage).Methods("POST")

	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images").Handler(http.StripPrefix("/images", imageHandler))
	authKey, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(authKey, csrf.Secure(false))
	http.ListenAndServe(":8080", csrfMw(userMw.Apply(r)))
}
