package server

import (
	"github.com/wangyuntao/duck/protocol"
)

func (s *service) start(evtCh chan<- *evt) {
	defer func() {
		evtCh <- &evt{
			t:    evtCloseService,
			port: s.port,
		}
	}()

	go func() {
		for {
			msg, err := protocol.ReadMsgFrom(s.conn)
			if err != nil {
				break
			}
			msgID := msg.ReadMsgID()
			switch msgID {
			case protocol.C2S_ConnectErr:
				evtCh <- &evt{
					t:    evtConnectErr,
					port: s.port,
				}
			default:
				panic("illegal msgID")
			}
		}
		s.ln.Close()
	}()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			break // TODO
		}
		evtCh <- &evt{
			t:    evtConnect,
			port: s.port,
			conn: conn,
		}
	}
}
