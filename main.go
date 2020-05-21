/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/

package main

import (
	"fmt"
	"net"

	"github.com/getgauge/flash/event"
	gm "github.com/getgauge/flash/gauge_messages"
	"github.com/getgauge/flash/listener"
	"google.golang.org/grpc"
)

func startAPI(e chan event.Event) {
	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		panic("failed to start server.")
	}
	l, err := net.ListenTCP("tcp", address)
	if err != nil {
		panic("failed to start server.")
	}
	server := grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 1024))
	h := listener.NewHandler(server, e)
	gm.RegisterReporterServer(server, h)
	fmt.Printf("Listening on port:%d\n", l.Addr().(*net.TCPAddr).Port)
	server.Serve(l)
}

func main() {
	e := make(chan event.Event)
	go startAPI(e)
	listener.NewHttpListener(e).Start()
}
