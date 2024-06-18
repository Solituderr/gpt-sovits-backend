package server

import (
	"github.com/gin-gonic/gin"
	"gptsovits/app/tts"
	"gptsovits/app/websock"
)

var r *gin.Engine

func Init() {
	var err error
	s := NewServerService()
	s.svs.ws = websock.NewBasicService(s.conf["gradio_url"])
	s.svs.tts, err = tts.NewBasicService(s.conf["gradio_url"], s.svs.ws)
	if err != nil {
		panic(err)
	}
	r = gin.Default()
	SetRoutes(r, s)
}

func Run() {
	r.Run(":8078")
}
