package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zjbztianya/poppy/models"
	"github.com/zjbztianya/poppy/util/context"
	"github.com/zjbztianya/poppy/views"
	"net/http"
	"strconv"
)

const (
	ShowGallery    = "show_gallery"
	IndexGalleries = "index_galleries"
	EditGallery    = "edit_gallery"
)

type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	is        models.ImageService
}

func NewGalleries(gs models.GalleryService, is models.ImageService, r *gin.Engine) *Galleries {
	return &Galleries{
		New:       views.NewView(r, "galleries_new", "galleries/new"),
		ShowView:  views.NewView(r, "galleries_show", "galleries/show"),
		EditView:  views.NewView(r, "galleries_edit", "galleries/edit"),
		IndexView: views.NewView(r, "galleries_index", "galleries/index"),
		gs:        gs,
		is:        is,
	}
}

type GalleryForm struct {
	Title string `scheme:"title"`
}

func (g *Galleries) Create(c *gin.Context) {
	var vd views.Response
	var form GalleryForm

	if err := parseForm(c, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(c, http.StatusInternalServerError, vd)
		return
	}

	user := context.User(c)
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(c, http.StatusInternalServerError, vd)
		return
	}

	url := fmt.Sprintf("/galleries/edit/%v", gallery.ID)
	c.Redirect(http.StatusFound, url)
}

func (g *Galleries) Show(c *gin.Context) {
	gallery, err := g.galleryByID(c)
	if err != nil {
		return
	}
	var vd views.Response
	vd.Data.Yield = gallery
	g.ShowView.Render(c, http.StatusOK, vd)
}

func (g *Galleries) galleryByID(c *gin.Context) (*models.Gallery, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		http.Error(c.Writer, "Invalid gallery ID", http.StatusFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(c.Writer, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(c.Writer, "Whoops! Something went wrong.",
				http.StatusInternalServerError)
		}
		return nil, err
	}
	images, _ := g.is.ByGalleryID(gallery.ID)
	gallery.Images = images
	return gallery, nil
}

func (g *Galleries) Edit(c *gin.Context) {
	gallery, err := g.galleryByID(c)
	if err != nil {
		return
	}

	user := context.User(c)
	if gallery.UserID != user.ID {
		http.Error(c.Writer, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}

	var vd views.Response
	vd.Data.Yield = gallery
	g.EditView.Render(c, http.StatusOK, vd)
}

func (g *Galleries) Update(c *gin.Context) {
	gallery, err := g.galleryByID(c)
	if err != nil {
		return
	}

	user := context.User(c)
	if gallery.UserID != user.ID {
		http.Error(c.Writer, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}

	var vd views.Response
	vd.Data.Yield = gallery
	var form GalleryForm
	if err := parseForm(c, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(c, http.StatusInternalServerError, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvSuccess,
			Message: "Gallery updated successfully!",
		}
	}
	g.EditView.Render(c, http.StatusOK, vd)
}

func (g *Galleries) Delete(c *gin.Context) {
	gallery, err := g.galleryByID(c)
	if err != nil {
		return
	}
	user := context.User(c)
	if gallery.UserID != user.ID {
		http.Error(c.Writer, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}

	var vd views.Response
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Data.Yield = gallery
		g.EditView.Render(c, http.StatusInternalServerError, vd)
		return
	}
	c.Redirect(http.StatusFound, "/galleries/index")
}

func (g *Galleries) Index(c *gin.Context) {
	user := context.User(c)
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(c.Writer, "Something went wrong.", http.StatusForbidden)
		return
	}

	var vd views.Response
	vd.Data.Yield = galleries
	g.IndexView.Render(c, http.StatusOK, vd)
}

func (g *Galleries) ImageUpload(c *gin.Context) {
	gallery, err := g.galleryByID(c)
	if err != nil {
		return
	}
	user := context.User(c)
	if gallery.UserID != user.ID {
		http.Error(c.Writer, "Gallery not found", http.StatusFound)
		return
	}

	var vd views.Response
	vd.Data.Yield = gallery
	form, err := c.MultipartForm()
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(c, http.StatusInternalServerError, vd)
		return
	}

	files := form.File["images"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(c, http.StatusInternalServerError, vd)
			return
		}
		defer file.Close()

		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(c, http.StatusInternalServerError, vd)
			return
		}
	}

	url := fmt.Sprintf("/galleries/edit/%v", gallery.ID)
	alert := views.Alert{
		Level:   views.AlertLvSuccess,
		Message: "Upload image success!",
	}
	views.RedirectAlert(c, url, http.StatusFound, alert)
}

func (g *Galleries) ImageDelete(c *gin.Context) {
	gallery, err := g.galleryByID(c)
	if err != nil {
		return
	}
	user := context.User(c)
	if gallery.UserID != user.ID {
		http.Error(c.Writer, "Gallery not found", http.StatusFound)
		return
	}

	filename := c.Param("filename")
	image := models.Image{
		GalleryID: gallery.ID,
		Filename:  filename,
	}
	err = g.is.Delete(&image)
	if err != nil {
		var vd views.Response
		vd.Data.Yield = gallery
		vd.SetAlert(err)
		g.EditView.Render(c, http.StatusInternalServerError, vd)
		return
	}

	url := fmt.Sprintf("/galleries/edit/%v", gallery.ID)
	c.Redirect(http.StatusFound, url)
}
