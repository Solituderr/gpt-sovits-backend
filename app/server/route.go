package server

import (
	"github.com/gin-gonic/gin"
	"gptsovits/app/tts"
)

func SetRoutes(r *gin.Engine, s *Service) {
	r.Use(tts.Cors())
	r.POST("/tts", tts.TTSHandler(s.svs.tts))
}
