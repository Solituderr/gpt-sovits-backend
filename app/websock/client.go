package websock

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type WebSocketService interface {
	NewWSClient()
	Send(data []byte)
	Receive() <-chan []byte
	Close()
}

type BasicService struct {
	Client *Client
	url    string
}

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	recv     chan []byte
	ExitChan chan bool
	close    bool
}

func (bs *BasicService) NewWSClient() {
	uri := url.URL{Scheme: "wss", Host: bs.url, Path: "/queue/join"}
	fmt.Println(uri.String())
	c, _, err := websocket.DefaultDialer.Dial(uri.String(), nil)
	if err != nil {
		log.Println(err)
	}
	client := &Client{
		conn:     c,
		send:     make(chan []byte),
		recv:     make(chan []byte),
		ExitChan: make(chan bool),
	}
	go client.receiveLoop()
	go client.sendLoop()
	bs.Client = client
}

func NewBasicService(url string) *BasicService {
	return &BasicService{url: url}
}
