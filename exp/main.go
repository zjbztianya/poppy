package main

import (
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("exp/hello.gohtml")
	if err != nil {
		panic(err)
	}
	data := struct {
		Name string
	}{"Jone Smith"}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
