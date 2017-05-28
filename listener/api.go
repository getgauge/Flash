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

package listener

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/getgauge/flash/event"
	m "github.com/getgauge/flash/gauge_messages"
	"github.com/golang/protobuf/proto"
)

type handlerFn func(*m.Message)

type apiListener struct {
	connection net.Conn
	handlers   map[m.Message_MessageType]handlerFn
	event      chan event.Event
}

func NewApiListener(host string, port string, e chan event.Event) (Listener, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}
	return &apiListener{connection: conn, handlers: map[m.Message_MessageType]handlerFn{
		m.Message_SuiteExecutionResult: func(msg *m.Message) {
			e <- event.NewEndEvent(msg.SuiteExecutionResult.SuiteResult.GetFailed())
		},
		m.Message_SpecExecutionStarting: func(msg *m.Message) {
			e <- event.NewSpecEvent(msg.SpecExecutionStartingRequest.CurrentExecutionInfo, true)
		},
		m.Message_SpecExecutionEnding: func(msg *m.Message) {
			e <- event.NewSpecEvent(msg.SpecExecutionEndingRequest.CurrentExecutionInfo, false)
		},
		m.Message_ScenarioExecutionStarting: func(msg *m.Message) {
			e <- event.NewScenarioEvent(msg.ScenarioExecutionStartingRequest.CurrentExecutionInfo, true)
		},
		m.Message_ScenarioExecutionEnding: func(msg *m.Message) {
			e <- event.NewScenarioEvent(msg.ScenarioExecutionEndingRequest.CurrentExecutionInfo, false)
		},
		m.Message_StepExecutionStarting: func(msg *m.Message) {
			e <- event.NewStepEvent(msg.StepExecutionStartingRequest.CurrentExecutionInfo, true)
		},
		m.Message_StepExecutionEnding: func(msg *m.Message) {
			e <- event.NewStepEvent(msg.StepExecutionEndingRequest.CurrentExecutionInfo, false)
		},
	}, event: e}, nil
}

func (l *apiListener) Start() {
	buffer := new(bytes.Buffer)
	data := make([]byte, 8192)
	for {
		n, err := l.connection.Read(data)
		if err != nil {
			return
		}
		buffer.Write(data[0:n])
		l.processMessages(buffer)
	}
}

func (l *apiListener) processMessages(buffer *bytes.Buffer) {
	for {
		messageLength, bytesRead := proto.DecodeVarint(buffer.Bytes())
		if messageLength > 0 && messageLength < uint64(buffer.Len()) {
			message := &m.Message{}
			messageBoundary := int(messageLength) + bytesRead
			err := proto.Unmarshal(buffer.Bytes()[bytesRead:messageBoundary], message)
			if err != nil {
				log.Printf("Failed to read proto message: %s\n", err.Error())
			} else {
				if message.MessageType == m.Message_KillProcessRequest {
					l.connection.Close()
					<-l.event
					os.Exit(0)
				} else {
					h := l.handlers[message.MessageType]
					if h != nil {
						h(message)
					}
				}
				buffer.Next(messageBoundary)
				if buffer.Len() == 0 {
					return
				}
			}
		} else {
			return
		}
	}
}
