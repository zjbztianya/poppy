package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zjbztianya/poppy/views"
)

type Static struct {
	Home    *views.View
	Contact *views.View
	FAQ     *views.View
}

func NewStatic(r *gin.Engine) *Static {
	return &Static{
		views.NewView(r,
			"static_home", "static/home"),
		views.NewView(r,
			"static_contact", "static/contact"),
		views.NewView(r,
			"static_faq", "static/faq"),
	}
}
