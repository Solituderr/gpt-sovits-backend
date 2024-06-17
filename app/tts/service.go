package tts

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	bb := GetRandStr()
	//tsr := TTSSendReq{Data: []}
	for {
		select {
		case msg := <-bs.ws.Receive():
			var trr TTSRecvResp
			if err := json.Unmarshal(msg, &trr); err != nil {
				return err
			}
			switch trr.Msg {
			case "send_data":
				bs.ws.Send([]byte{})
			case "send_hash":
				tmp := map[string]interface{}{"fn_index": 1, "session_hash": bb}
				tt, err := json.Marshal(tmp)
				if err != nil {
					return err
				}
				bs.ws.Send(tt)
			case "process_completed":
				return nil
			}
		}
	}
}

func (bs *BasicService) GetWavAudio(req Req) (string, error) {
	//bb := GetRandStr()
	bb := "smwof9n92f"
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
		FnIndex:     1,
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
				tmp := map[string]interface{}{"fn_index": 1, "session_hash": bb}
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

func (bs *BasicService) GetTTSInfo() error {
	uri := url.URL{Scheme: "https", Host: bs.url, Path: "/info"}
	resp, err := http.Get(uri.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var info TTSInfo
	if err = json.Unmarshal(body, &info); err != nil {
		return err
	}
	if len(info.UnnamedEndpoints) <= 0 {
		return errors.New("获取info失败")
	}
	sovits := info.UnnamedEndpoints["1"].Parameters[0].PythonType.Description
	sovitsModels := ExtractModel(sovits)
	gpt := info.UnnamedEndpoints["2"].Parameters[0].PythonType.Description
	gptModels := ExtractModel(gpt)
	fmt.Println(sovitsModels)
	fmt.Println(gptModels)
	return nil
}

func ExtractModel(fullString string) []string {
	// 找到第一个冒号后的字符串，然后去掉前面的 "Option from: " 和两边的空格
	cutString := strings.TrimPrefix(fullString, "Option from: ")
	cutString = strings.TrimSpace(cutString)

	// 去除两端的中括号
	cutString = strings.Trim(cutString, "['']")

	// 使用 "', '" 分割字符串得到路径数组
	modelPaths := strings.Split(cutString, "', '")

	// 打印结果
	var s []string
	for _, path := range modelPaths {
		s = append(s, path)
	}
	return s
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
