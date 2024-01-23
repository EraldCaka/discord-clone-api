package pkg

import "github.com/gorilla/websocket"

func Dialer() {
	websocket.DefaultDialer.Dial("ws:/oo", nil)
}

type Consumer interface {
	Start() error
}
