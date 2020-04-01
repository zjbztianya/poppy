package conf

import (
	"github.com/go-ini/ini"
	"log"
)

var Conf = &Config{}

type Config struct {
	App      App
	Database Database
	Server   Server
}

type App struct {
	UserPwPepper       string
	HmacSecretKey      string
	UserKey            string
	RememberTokenBytes int
	AuthKeyBytes       int
	CsrfSecure         bool
}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
	LogMode  bool
}

type Server struct {
	RunMode  string
	HttpPort int
}

var cfg *ini.File

func Init() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("config Init err: %v", err)
	}
	mapTo("app", &Conf.App)
	mapTo("server", &Conf.Server)
	mapTo("database", &Conf.Database)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("config mapTo %s err: %v", section, err)
	}
}
