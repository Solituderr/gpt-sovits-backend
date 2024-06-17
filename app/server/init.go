package server

import (
	"github.com/gin-gonic/gin"
	"gptsovits/app/tts"
	"gptsovits/app/websock"
)

var r *gin.Engine

func Init() {
	s := NewServerService()
	s.svs.ws = websock.NewBasicService(s.conf["gradio_url"])
	s.svs.tts = tts.NewBasicService(s.conf["gradio_url"], s.svs.ws)
	r = gin.Default()
	SetRoutes(r, s)
}

func Run() {
	r.Run(":8078")
}
