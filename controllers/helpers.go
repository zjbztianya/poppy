package controllers

import (
	"github.com/gorilla/schema"
	"net/http"
)

const maxMultipartMem = 1 << 20

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	return dec.Decode(dst, r.PostForm)
}
