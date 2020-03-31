package views

import (
	"github.com/gin-gonic/gin"
	"github.com/zjbztianya/poppy/models"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	AlertLvError    = "danger"
	AlertLvWarning  = "waring"
	AlertLvInfo     = "info"
	AlertLvSuccess  = "success"
	AlertMsgGeneric = "Something went wrong. Please try" +
		"again, and contact us if the problem persists."
)

type Alert struct {
	Level   string
	Message string
}

type Response struct {
	Alert *Alert
	User  *models.User
	Data  struct {
		Yield     interface{}
		CsrfField template.HTML
	}
}

type PublicError interface {
	error
	Public() string
}

func (d *Response) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvError,
		Message: msg,
	}
}

func (d *Response) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvError,
		Message: msg,
	}
}

func persistAlert(c *gin.Context, alert Alert) {
	expiresAt := time.Now().Add(5 * time.Minute)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &lvl)
	http.SetCookie(c.Writer, &msg)
}

func clearAlert(c *gin.Context) {
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &lvl)
	http.SetCookie(c.Writer, &msg)
}

func getAlert(c *gin.Context) *Alert {
	lvl, err := c.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := c.Cookie("alert_message")
	if err != nil {
		return nil
	}
	alert := Alert{
		Level:   lvl,
		Message: msg,
	}
	return &alert
}

func RedirectAlert(c *gin.Context, urlStr string, code int, alert Alert) {
	persistAlert(c, alert)
	c.Redirect(code, urlStr)
}
