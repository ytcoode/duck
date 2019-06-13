package util

import (
	"io"
	"net"
)

const (
	relayClose = iota
	relayCloseWait
)

func Relay(c1, c2 *net.TCPConn) {
	ch := make(chan int, 2)
	go cpy(c1, c2, ch)
	go cpy(c2, c1, ch)

	for i := 0; i < 2; i++ {
		if v := <-ch; v == relayClose {
			break
		}
	}
	c1.Close()
	c2.Close()
}

func cpy(c1, c2 *net.TCPConn, ch chan int) {
	_, err := io.Copy(c2, c1)
	if err != nil {
		ch <- relayClose
		return
	}
	err = c2.CloseWrite()
	if err != nil {
		ch <- relayClose
		return
	}
	ch <- relayCloseWait
}
