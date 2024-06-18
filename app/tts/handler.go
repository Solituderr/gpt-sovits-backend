package tts

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TTSService interface {
	BindCharacter(model string, hash string, fnIndex int) error
	GetWavAudio(req Req) (string, error)
	GetTTSInfo() (sovitsModels []string, gptModels []string, err error)
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

func InfoHandler(service TTSService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, _, err := service.GetTTSInfo()
		fmt.Println(err)
	}
}
