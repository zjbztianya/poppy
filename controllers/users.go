package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zjbztianya/poppy/models"
	"github.com/zjbztianya/poppy/util/context"
	"github.com/zjbztianya/poppy/util/rand"
	"github.com/zjbztianya/poppy/views"
	"net/http"
	"time"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

func NewUsers(us models.UserService, r *gin.Engine) *Users {
	return &Users{
		NewView:   views.NewView(r, "users_new", "users/new"),
		LoginView: views.NewView(r, "users_login", "users/login"),
		us:        us,
	}
}

type SignupForm struct {
	Name     string `form:"name" binding:"required,max=30"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,gte=8,lte=16"`
}

type LoginForm struct {
	Email    string `email:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,gte=8,lte=16"`
}

func (u *Users) Create(c *gin.Context) {
	var vd views.Response
	var form SignupForm
	if err := c.Bind(&form); err != nil {
		vd.SetAlert(models.ErrBadRequst)
		u.NewView.Render(c, http.StatusBadRequest, vd)
		return
	}
	fmt.Println(form)
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(c, http.StatusInternalServerError, vd)
		return
	}

	err := u.signIn(c, &user)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.Redirect(http.StatusFound, "/galleries/index")
}

func (u *Users) Login(c *gin.Context) {
	var vd views.Response
	var form LoginForm
	if err := c.Bind(&form); err != nil {
		vd.SetAlert(models.ErrBadRequst)
		u.LoginView.Render(c, http.StatusBadRequest, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(c, http.StatusInternalServerError, vd)
		return
	}

	err = u.signIn(c, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(c, http.StatusInternalServerError, vd)
		return
	}
	c.Redirect(http.StatusFound, "/galleries/index")
}

func (u *Users) CookieTest(c *gin.Context) {
	cookie, err := c.Cookie("remember_token")
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie)
	if err != nil {
		fmt.Println(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(c.Writer, user)
}

func (u *Users) signIn(c *gin.Context, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	return nil
}

func (u *Users) Logout(c *gin.Context) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	user := context.User(c)
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	c.Redirect(http.StatusFound, "/")
}
