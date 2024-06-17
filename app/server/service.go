package server

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gptsovits/app/tts"
	"gptsovits/app/websock"
	"log"
)

type Service struct {
	svs  BasicService
	db   *gorm.DB
	conf map[string]string
}

type BasicService struct {
	ws  *websock.BasicService
	tts *tts.BasicService
}

func NewServerService() *Service {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}
	var s Service
	s.conf = viper.GetStringMapString("tts")
	return &s
}
