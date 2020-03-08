package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zjbztianya/poppy/models"
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
	//user:password@tcp(localhost:5555)/dbname?
	sqlInfo := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", user, password, dbname)
	us, err := models.NewUserService(sqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	user := models.User{
		Name:  "zjbztianya",
		Email: "zjbztianya@163.com",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}

	user.Name = "hanhan"
	if err := us.Update(&user); err != nil {
		panic(err)
	}

	foundUser, err := us.ByEmail("zjbztianya@163.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser)

	if err := us.Delete(foundUser.ID); err != nil {
		panic(err)
	}
	_, err = us.ById(foundUser.ID)
	if err != models.ErrNotFound {
		panic("user not delete!")
	}
}
