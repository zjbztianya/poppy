package controllers

import (
	"fmt"
	"github.com/zjbztianya/poppy/models"
	"github.com/zjbztianya/poppy/views"
	"net/http"
)

type Users struct {
	NewView *views.View
	us      *models.UserService
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name:  form.Name,
		Email: form.Email,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User is", user)
}
