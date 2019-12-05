// Copyright 2015 ThoughtWorks, Inc.

// This file is part of getgauge/flash.

// getgauge/flash is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// getgauge/flash is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with getgauge/flash.  If not, see <http://www.gnu.org/licenses/>.

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
	server := grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 1024 * 10))
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
