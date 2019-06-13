package server

import (
	"net"

	"bytes"

	"github.com/wangyuntao/duck/protocol"
)

func Start(addr, passwd string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := &server{
		ln:     ln,
		passwd: []byte(passwd),
		svcMap: make(map[int]*service),
		evtCh:  make(chan *evt, 128),
	}

	go s.eventLoop()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err) // TODO: Better way?
		}
		go s.handle(c)
	}
}

func (s *server) handle(conn net.Conn) {
	closeConn := true
	defer func() {
		if closeConn {
			conn.Close()
		}
	}()

	msg, err := protocol.ReadMsgFrom(conn)
	if err != nil {
		return
	}

	msgID := msg.ReadMsgID()
	passwd := msg.ReadBytes()

	if !bytes.Equal(passwd, s.passwd) {
		return
	}

	port := msg.ReadUint16()

	switch msgID {
	case protocol.C2S_Listen:
		s.evtCh <- &evt{evtListen, int(port), conn}
	case protocol.C2S_Relay:
		s.evtCh <- &evt{evtRelay, int(port), conn}
	default:
		return
	}

	closeConn = false
}
