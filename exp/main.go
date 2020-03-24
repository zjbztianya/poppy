package main

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	user     = "root"
	password = "123456"
	dbname   = "test"
)

type User struct {
	gorm.Model
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Orders []Order
}

type Order struct {
	gorm.Model
	UserId      uint
	Amount      int
	Description string
}

func main() {
	ctx := context.TODO()
	type privateKey string
	var a privateKey = "user"
	ctx = context.WithValue(ctx, "user", "123")
	ctx = context.WithValue(ctx, a, 456)
	fmt.Println(ctx.Value("user"))
}
