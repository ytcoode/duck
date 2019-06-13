package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/wangyuntao/duck/protocol"
	"github.com/wangyuntao/duck/util"
)

func Start(serverAddr string, ports []uint16, passwd string) {
	for _, port := range ports {
		go start(serverAddr, port, []byte(passwd))
	}

	for {
		// Better idea?
		time.Sleep(time.Hour)
	}
}

func start(serverAddr string, port uint16, passwd []byte) {
	defer func() {
		time.Sleep(time.Second)
		go start(serverAddr, port, passwd)
	}()

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	msg := protocol.NewMsg(protocol.C2S_Listen)
	msg.WriteBytes(passwd)
	msg.WriteUint16(port)

	err = msg.WriteTo(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for {
		msg, err := protocol.ReadMsgFrom(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}

		msgID := msg.ReadMsgID()
		switch msgID {
		case protocol.S2C_Connect:
			go connect(serverAddr, port, passwd, conn)
		}
	}
}

func connect(serverAddr string, port uint16, passwd []byte, conn net.Conn) {
	localAddr := "127.0.0.1:" + strconv.Itoa(int(port))
	c1, err := net.Dial("tcp", localAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		msg := protocol.NewMsg(protocol.C2S_ConnectErr)
		msg.WriteTo(conn)
		return
	}

	c2, err := net.Dial("tcp", serverAddr)
	if err != nil {
		c1.Close()
		return
	}

	msg := protocol.NewMsg(protocol.C2S_Relay)
	msg.WriteBytes(passwd)
	msg.WriteUint16(port)

	err = msg.WriteTo(c2)
	if err != nil {
		c1.Close()
		c2.Close()
		fmt.Fprintln(os.Stderr, err)
		return
	}

	go util.Relay(c1.(*net.TCPConn), c2.(*net.TCPConn))
}
