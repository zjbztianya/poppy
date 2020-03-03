package controllers

import "github.com/zjbztianya/poppy/views"

type Static struct {
	Home    *views.View
	Contact *views.View
	FAQ     *views.View
}

func NewStatic() *Static {
	return &Static{
		views.NewView(
			"bootstrap", "static/home"),
		views.NewView(
			"bootstrap", "static/contact"),
		views.NewView(
			"bootstrap", "static/faq"),
	}
}
