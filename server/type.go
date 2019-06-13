package server

import (
	"net"
)

type server struct {
	ln     net.Listener
	passwd []byte
	svcMap map[int]*service
	evtCh  chan *evt
}

type service struct {
	ln      net.Listener
	port    int
	conn    net.Conn
	clients []net.Conn
}
