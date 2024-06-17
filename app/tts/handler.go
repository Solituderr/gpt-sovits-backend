package tts

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TTSService interface {
	BindCharacter(req Req) error
	GetWavAudio(req Req) (string, error)
}

func TTSHandler(service TTSService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.BindJSON(&req); err != nil {
			_ = ctx.Error(err)
			return
		}
		if audio, err := service.GetWavAudio(req); err != nil {
			_ = ctx.Error(err)
		} else {
			ctx.String(http.StatusOK, audio)
		}
	}
}
