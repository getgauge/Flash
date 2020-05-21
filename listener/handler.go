/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/
package listener

import (
	"context"
	"os"

	"github.com/getgauge/flash/event"
	gm "github.com/getgauge/flash/gauge_messages"
	"google.golang.org/grpc"
)

type handler struct {
	server *grpc.Server
	e      chan event.Event
}

func NewHandler(s *grpc.Server, e chan event.Event) *handler {
	return &handler{server: s, e: e}
}

func (h *handler) NotifyExecutionStarting(c context.Context, m *gm.ExecutionStartingRequest) (*gm.Empty, error) {
	return &gm.Empty{}, nil
}
func (h *handler) NotifySpecExecutionStarting(c context.Context, m *gm.SpecExecutionStartingRequest) (*gm.Empty, error) {
	h.e <- event.NewSpecEvent(m.CurrentExecutionInfo, true)
	return &gm.Empty{}, nil
}
func (h *handler) NotifyScenarioExecutionStarting(c context.Context, m *gm.ScenarioExecutionStartingRequest) (*gm.Empty, error) {
	h.e <- event.NewScenarioEvent(m.CurrentExecutionInfo, true)
	return &gm.Empty{}, nil
}
func (h *handler) NotifyStepExecutionStarting(c context.Context, m *gm.StepExecutionStartingRequest) (*gm.Empty, error) {
	h.e <- event.NewStepEvent(m.CurrentExecutionInfo, true)
	return &gm.Empty{}, nil
}
func (h *handler) NotifyStepExecutionEnding(c context.Context, m *gm.StepExecutionEndingRequest) (*gm.Empty, error) {
	h.e <- event.NewStepEvent(m.CurrentExecutionInfo, false)
	return &gm.Empty{}, nil
}
func (h *handler) NotifyScenarioExecutionEnding(c context.Context, m *gm.ScenarioExecutionEndingRequest) (*gm.Empty, error) {
	h.e <- event.NewScenarioEvent(m.CurrentExecutionInfo, false)
	return &gm.Empty{}, nil
}
func (h *handler) NotifySpecExecutionEnding(c context.Context, m *gm.SpecExecutionEndingRequest) (*gm.Empty, error) {
	h.e <- event.NewSpecEvent(m.CurrentExecutionInfo, false)
	return &gm.Empty{}, nil
}
func (h *handler) NotifyExecutionEnding(c context.Context, m *gm.ExecutionEndingRequest) (*gm.Empty, error) {
	return &gm.Empty{}, nil
}

func (h *handler) NotifySuiteResult(c context.Context, m *gm.SuiteExecutionResult) (*gm.Empty, error) {
	h.e <- event.NewEndEvent(m.SuiteResult.GetFailed())
	return &gm.Empty{}, nil
}

func (h *handler) Kill(c context.Context, m *gm.KillProcessRequest) (*gm.Empty, error) {
	defer h.stopServer()
	return &gm.Empty{}, nil
}

func (h *handler) stopServer() {
	h.server.Stop()
	os.Exit(0)
}
