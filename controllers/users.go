package controllers

import (
	"fmt"
	"github.com/zjbztianya/poppy/views"
	"net/http"
)

type Users struct {
	NewView *views.View
}

func NewUers() *Users {
	return &Users{NewView: views.NewView("bootstrap", "views/users/new.gohtml")}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Email is ", form.Email)
	fmt.Fprintln(w, "Password is ", form.Password)
}
