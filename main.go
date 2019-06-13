package main

import (
	"flag"
	"fmt"
	"os"

	"strconv"

	"github.com/wangyuntao/duck/client"
	"github.com/wangyuntao/duck/server"
)

func main() {
	var (
		serverMode = flag.Bool("l", false, "server mode")
		serverAddr = flag.String("addr", ":9990", "server address")
		password   = flag.String("p", "", "password")
	)

	flag.Parse()

	if len(*serverAddr) == 0 || len(*password) == 0 {
		flag.Usage()
		return
	}

	if *serverMode {
		err := server.Start(*serverAddr, *password)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("Must specify listen ports")
		return
	}

	ports := make([]uint16, flag.NArg())
	for i, s := range flag.Args() {
		p, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("illegal listen port")
			return
		}
		ports[i] = uint16(p)
	}

	client.Start(*serverAddr, ports, *password)
}
