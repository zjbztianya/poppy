package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

func parseForm(c *gin.Context, dst interface{}) error {
	if err := c.Request.ParseForm(); err != nil {
		return err
	}

	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	return dec.Decode(dst, c.Request.PostForm)
}
