package tts

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type WebSocketService interface {
	NewWSClient()
	Send(data []byte)
	Receive() <-chan []byte
	Close()
}

type BasicService struct {
	ws  WebSocketService
	url string
}

func (bs *BasicService) BindCharacter(req Req) error {
	//TODO implement me
	panic("implement me")
}

func (bs *BasicService) GetWavAudio(req Req) (string, error) {
	//bb := GetRandStr()
	bb := "d1h25ac2lm"
	tsr := TTSSendReq{
		Data: []interface{}{
			MetaData{
				Data: "ata:audio/wav;base64," + AudioToBase64(),
				Name: "jay参考.wav",
			},
			"哎后来他听到我的歌，他说，你这些歌曲别人不用干脆你自己唱唱看好了",
			"中文",
			req.Text,
			req.Language,
			"凑四句一切",
			5,
			1,
			1,
			false,
		},
		SessionHash: bb,
		FnIndex:     3,
	}
	bs.ws.NewWSClient()
	defer bs.ws.Close()
	data, _ := json.Marshal(tsr)
	var audio string
	for {
		select {
		case msg := <-bs.ws.Receive():
			var trr TTSRecvResp
			if err := json.Unmarshal(msg, &trr); err != nil {
				return "", err
			}
			switch trr.Msg {
			case "send_data":
				bs.ws.Send(data)
			case "send_hash":
				tmp := map[string]interface{}{"fn_index": 3, "session_hash": bb}
				tt, err := json.Marshal(tmp)
				if err != nil {
					return "", err
				}
				bs.ws.Send(tt)
			case "process_completed":
				audio = fmt.Sprintf("https://%s/file=%s", bs.url, trr.Output.Data[0].Name)
				return audio, nil
			}
		}
	}
}

func GetRandStr() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := 11
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func AudioToBase64() string {
	// 读取 WAV 文件
	data, err := os.ReadFile("./jay参考.wav")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	// 将二进制数据编码为 Base64
	encodedData := base64.StdEncoding.EncodeToString(data)

	// 打印或使用 Base64 编码后的数据
	return encodedData
}

func NewBasicService(url string, ws WebSocketService) *BasicService {
	return &BasicService{
		url: url,
		ws:  ws,
	}
}
