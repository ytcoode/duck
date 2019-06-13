package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/wangyuntao/duck/protocol"
	"github.com/wangyuntao/duck/util"
)

const (
	evtListen = iota
	evtConnect
	evtConnectErr
	evtRelay
	evtCloseService
)

type evt struct {
	t    int
	port int
	conn net.Conn
}

func (s *server) eventLoop() {
	// TODO tick check
	for {
		e := <-s.evtCh
		//fmt.Println("evt:", e.t, "port:", e.port)

		switch e.t {
		case evtListen:
			s.evtListen(e)
		case evtConnect:
			s.evtConnect(e)
		case evtConnectErr:
			s.evtConnectErr(e)
		case evtRelay:
			s.evtRelay(e)
		case evtCloseService:
			s.evtCloseService(e)
		default:
			panic(fmt.Sprintf("illegal event type: %d", e.t))
		}
	}
}

func (s *server) evtListen(e *evt) {
	if s.svcMap[e.port] != nil {
		e.conn.Close()
		return
	}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(e.port))
	if err != nil {
		e.conn.Close()
		return
	}

	svc := &service{
		ln:      ln,
		port:    e.port,
		conn:    e.conn,
		clients: []net.Conn{},
	}

	s.svcMap[e.port] = svc
	go svc.start(s.evtCh)
}

func (s *server) evtConnect(e *evt) {
	svc := s.svcMap[e.port]
	if svc == nil {
		e.conn.Close()
		return
	}

	msg := protocol.NewMsg(protocol.S2C_Connect)
	err := msg.WriteTo(svc.conn)
	if err != nil {
		e.conn.Close()
		return
	}
	svc.clients = append(svc.clients, e.conn)
}

func (s *server) evtConnectErr(e *evt) {
	svc := s.svcMap[e.port]
	if svc == nil || len(svc.clients) == 0 {
		return
	}
	svc.clients[0].Close()
	svc.clients = svc.clients[1:]
}

func (s *server) evtRelay(e *evt) {
	svc := s.svcMap[e.port]
	if svc == nil || len(svc.clients) == 0 {
		e.conn.Close()
		return
	}
	go util.Relay(e.conn.(*net.TCPConn), svc.clients[0].(*net.TCPConn))
	svc.clients = svc.clients[1:]
}

func (s *server) evtCloseService(e *evt) {
	svc := s.svcMap[e.port]
	if svc == nil {
		return
	}
	delete(s.svcMap, e.port)
	svc.ln.Close()
	svc.conn.Close()
	for _, c := range svc.clients {
		c.Close()
	}
}
