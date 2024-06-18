package websock

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

func (c *Client) sendLoop() {
	for {
		select {
		case msg := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println(err)
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Client) receiveLoop() {
	for {
		if c.close {
			break
		}
		if _, p, err := c.conn.ReadMessage(); err != nil {
			if strings.Contains(err.Error(), "1000") {
				break
			}
			log.Println(err)
		} else {
			fmt.Println(string(p))
			c.recv <- p
		}
	}
}

func (bs *BasicService) Send(data []byte) {
	bs.Client.send <- data
}

func (bs *BasicService) Receive() <-chan []byte {
	return bs.Client.recv
}

func (bs *BasicService) Close() {
	if bs.Client.close {
		return
	}
	bs.Client.close = true
	bs.Client.conn.Close()
	close(bs.Client.ExitChan)
	close(bs.Client.recv)
	close(bs.Client.send)
}
