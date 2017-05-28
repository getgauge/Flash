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
	"os"

	"github.com/getgauge/flash/event"
	"github.com/getgauge/flash/listener"
)

const (
	gaugeHost    = "localhost"
	gaugePortEnv = "plugin_connection_port"
)

func main() {
	e := make(chan event.Event)
	a, err := listener.NewApiListener(gaugeHost, os.Getenv(gaugePortEnv), e)
	if err != nil {
		fmt.Println("Could not create the gauge listener")
		os.Exit(1)
	}
	go a.Start()
	listener.NewHttpListener(e).Start()
}
