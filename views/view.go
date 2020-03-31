package views

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/zjbztianya/poppy/context"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type View struct {
	Name string
}

var (
	LayoutDir   = "views/layouts/"
	TemplateDir = "views/"
	TemplateExt = ".gohtml"
)

func NewView(r *gin.Engine, name string, files ...string) *View {
	render, ok := r.HTMLRender.(multitemplate.Renderer)
	if !ok {
		panic(errors.New("not set multitemplate render!"))
	}
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	for i, file := range files {
		if strings.HasSuffix(file, "bootstrap.gohtml") {
			files[0], files[i] = files[i], files[0]
			break
		}
	}
	funMap := template.FuncMap{
		"pathEscape": func(s string) string {
			return url.PathEscape(s)
		},
	}
	render.AddFromFilesFuncs(name, funMap, files...)
	return &View{name}
}

func (v *View) Render(c *gin.Context, code int, data interface{}) {
	var vd Response
	switch d := data.(type) {
	case Response:
		vd = d
	default:
		vd = Response{
			Data: struct {
				Yield     interface{}
				CsrfField template.HTML
			}{Yield: data},
		}
	}
	if alert := getAlert(c); alert != nil {
		vd.Alert = alert
		clearAlert(c)
	}
	vd.User = context.User(c)
	vd.Data.CsrfField = csrf.TemplateField(c.Request)
	c.HTML(code, v.Name, vd)
}

func (v *View) HTML(c *gin.Context) {
	v.Render(c, http.StatusOK, nil)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
