package event

import m "github.com/getgauge/flash/gauge_messages"

type status string
type eventType string

const (
	pass     status    = "pass"
	fail     status    = "fail"
	progress status    = "progress"
	spec     eventType = "spec"
	scenario eventType = "scenario"
	step     eventType = "step"
	end      eventType = "end"
)

type Event interface {
}

type specEvent struct {
	Name     string
	Status   status
	FileName string
	Type     eventType
}

type scenarioEvent struct {
	Name         string
	Status       status
	SpecFileName string
	Type         eventType
}

type stepEvent struct {
	Name         string
	Status       status
	ScenarioName string
	SpecFileName string
	Type         eventType
}

type endEvent struct {
	Status status
	Type   eventType
}

func NewSpecEvent(i *m.ExecutionInfo, hasStarted bool) specEvent {
	s := pass
	if i.CurrentSpec.GetIsFailed() {
		s = fail
	}
	if hasStarted {
		s = progress
	}
	return specEvent{
		Name:     i.CurrentSpec.Name,
		FileName: i.CurrentSpec.FileName,
		Status:   s,
		Type:     spec,
	}
}

func NewScenarioEvent(i *m.ExecutionInfo, hasStarted bool) scenarioEvent {
	s := pass
	if i.CurrentScenario.GetIsFailed() {
		s = fail
	}
	if hasStarted {
		s = progress
	}
	return scenarioEvent{
		Name:         i.CurrentScenario.Name,
		SpecFileName: i.CurrentSpec.FileName,
		Status:       s,
		Type:         scenario,
	}
}

func NewStepEvent(i *m.ExecutionInfo, hasStarted bool) stepEvent {
	s := pass
	if i.CurrentStep.GetIsFailed() {
		s = fail
	}
	if hasStarted {
		s = progress
	}
	return stepEvent{
		Name:         i.CurrentStep.Step.ActualStepText,
		Status:       s,
		ScenarioName: i.CurrentScenario.Name,
		SpecFileName: i.CurrentSpec.FileName,
		Type:         step,
	}
}

func NewEndEvent(isFailed bool) endEvent {
	s := pass
	if isFailed {
		s = fail
	}
	return endEvent{Status: s, Type: end}
}
