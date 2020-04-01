package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/zjbztianya/poppy/conf"
)

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

func NewServices() (*Services, error) {
	dbConf := conf.Conf.Database
	connectionInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConf.User, dbConf.Password, dbConf.Host, dbConf.Name)
	db, err := gorm.Open(dbConf.Type, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(dbConf.LogMode)

	return &Services{
		Gallery: NewGalleryService(db),
		User:    NewUserService(db),
		Image:   NewImageService(),
		db:      db}, nil
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
