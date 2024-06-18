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

const (
	fnIndex1 = 1
	fnIndex2 = 2
	fnIndex3 = 3
)

type WebSocketService interface {
	NewWSClient()
	Send(data []byte)
	Receive() <-chan []byte
	Close()
}

type BasicService struct {
	ws           WebSocketService
	url          string
	characters   []string
	characterMap map[string]*Character
}

func (bs *BasicService) BindCharacter(model string, hash string, fnIndex int) error {
	bs.ws.NewWSClient()
	defer bs.ws.Close()
	tsr := TTSSendReq{Data: []interface{}{model}, FnIndex: fnIndex, SessionHash: hash}
	data, err := json.Marshal(tsr)
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-bs.ws.Receive():
			var trr TTSRecvResp
			if err = json.Unmarshal(msg, &trr); err != nil {
				return err
			}
			switch trr.Msg {
			case "send_data":
				bs.ws.Send(data)
			case "send_hash":
				tmp := map[string]interface{}{"fn_index": 1, "session_hash": hash}
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
	bb := GetRandStr()
	cInfo := bs.characterMap[req.Character]
	if err := bs.BindCharacter(cInfo.ModelSoVits, bb, fnIndex1); err != nil {
		return "", err
	}
	if err := bs.BindCharacter(cInfo.ModelGPT, bb, fnIndex2); err != nil {
		return "", err
	}
	tsr := TTSSendReq{
		Data: []interface{}{
			MetaData{
				Data: cInfo.AudioBase64,
				Name: cInfo.Name,
			},
			cInfo.Words,
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
	data, err := json.Marshal(tsr)
	if err != nil {
		return "", err
	}
	var audio string
	for {
		select {
		case msg := <-bs.ws.Receive():
			var trr TTSRecvResp
			if err = json.Unmarshal(msg, &trr); err != nil {
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

func (bs *BasicService) GetTTSInfo() (sovitsModels []string, gptModels []string, err error) {
	uri := url.URL{Scheme: "https", Host: bs.url, Path: "/info"}
	resp, err := http.Get(uri.String())
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var info TTSInfo
	if err = json.Unmarshal(body, &info); err != nil {
		return nil, nil, err
	}
	if len(info.UnnamedEndpoints) <= 0 {
		return nil, nil, errors.New("获取info失败")
	}
	sovits := info.UnnamedEndpoints["1"].Parameters[0].PythonType.Description
	sovitsModels = ExtractModel(sovits)
	gpt := info.UnnamedEndpoints["2"].Parameters[0].PythonType.Description
	gptModels = ExtractModel(gpt)
	return sovitsModels, gptModels, nil
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

func AudioToBase64(path string) string {
	// 读取 WAV 文件
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	// 将二进制数据编码为 Base64
	encodedData := base64.StdEncoding.EncodeToString(data)

	// 打印或使用 Base64 编码后的数据
	return encodedData
}

func NewBasicService(url string, ws WebSocketService) (*BasicService, error) {
	bs := &BasicService{
		url:          url,
		ws:           ws,
		characterMap: make(map[string]*Character),
	}
	// 获取 model 信息
	sovits, gpts, err := bs.GetTTSInfo()
	if err != nil {
		return nil, err
	}
	// 初始化 角色信息
	audioFilePath := "./audio"
	files, err := os.ReadDir(audioFilePath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			path := audioFilePath + "/" + name
			lis := strings.Split(name, "+")
			character, words := lis[0], lis[1]
			bs.characters = append(bs.characters, character)
			var sovit, gpt string
			for i := 0; i < len(sovits); i++ {
				if strings.Contains(sovits[i], character) {
					sovit = sovits[i]
					break
				}
			}
			for i := 0; i < len(gpts); i++ {
				if strings.Contains(gpts[i], character) {
					gpt = gpts[i]
					break
				}
			}
			if sovit == "" || gpt == "" {
				return nil, errors.New("没有数据找到模型")
			}
			bs.characterMap[character] = &Character{
				Name:        name,
				ModelSoVits: sovit,
				ModelGPT:    gpt,
				Words:       words,
				AudioBase64: "ata:audio/wav;base64," + AudioToBase64(path),
				Hash:        "",
			}
		}
	}
	return bs, nil
}
