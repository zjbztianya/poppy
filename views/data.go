package views

import (
	"github.com/zjbztianya/poppy/models"
	"log"
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

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

type PublicError interface {
	error
	Public() string
}

func (d *Data) SetAlert(err error) {
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

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvError,
		Message: msg,
	}
}
