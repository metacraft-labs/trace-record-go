package main

import "fmt"

type FunctionId uint64
type CallId uint64
type VariableId uint64
type StepId uint64
type PathId uint64
type Line int64

type RecordEvent interface {
	isRecordEvent() bool
}

type StepRecord struct {
	pathId PathId
	line   Line
}

func (s StepRecord) isRecordEvent() bool {
	return true
}

// ====

type FunctionRecord struct {
	name   string
	pathId PathId
	line   Line
}

func (r FunctionRecord) isRecordEvent() bool {
	return true
}

// ====

type ValueRecord interface {
	isValueRecord() bool
}

type NilValueRecord struct {
}

func (n NilValueRecord) isValueRecord() bool {
	return true
}

type ArgRecord struct {
	name  string
	value ValueRecord
}

type CallRecord struct {
	functionId FunctionId
	args       []ArgRecord
}

func (c CallRecord) isRecordEvent() bool {
	return true
}

// ====

type ReturnRecord struct {
	returnValue ValueRecord
}

func (r ReturnRecord) isRecordEvent() bool {
	return true
}

// ====

type PathRecord string

func (p PathRecord) isRecordEvent() bool {
	return true
}

// ====

type TraceRecord struct {
	events    []RecordEvent
	functions map[string]FunctionId
	paths     map[string]PathId
}

func MakeTraceRecord() TraceRecord {
	events := make([]RecordEvent, 0)
	functions := make(map[string]FunctionId, 0)
	paths := make(map[string]PathId, 0)
	return TraceRecord{events, functions, paths}
}

func (t *TraceRecord) Register(event RecordEvent) {
	t.events = append(t.events, event)
}

func (t *TraceRecord) RegisterStep(pathId PathId, line Line) {
	step := StepRecord{pathId, line}
	t.Register(step)
}

func (t *TraceRecord) RegisterCall(name string, functionPathId PathId, functionStartLine Line) {
	functionId := t.EnsureFunctionId(name, functionPathId, functionStartLine)
	call := CallRecord{functionId, make([]ArgRecord, 0)}
	t.Register(call)
}

func (t *TraceRecord) RegisterFunctionWithNewId(name string, pathId PathId, line Line) FunctionId {
	// doesn't check if name is already registered, if you want the check, use `EnsureFunctionId` !
	r := FunctionRecord{name, pathId, line}
	t.Register(r)
	newFunctionId := FunctionId(len(t.functions))
	t.functions[name] = newFunctionId
	return newFunctionId
}

func (t *TraceRecord) EnsureFunctionId(name string, pathId PathId, line Line) FunctionId {
	functionId, ok := t.functions[name]
	if !ok {
		functionId = t.RegisterFunctionWithNewId(name, pathId, line)
	}
	return functionId
}

func (t *TraceRecord) RegisterReturn(returnValue ValueRecord) {
	r := ReturnRecord{returnValue}
	t.Register(r)
}

func (t *TraceRecord) RegisterPathWithNewId(path string) PathId {
	newPathId := PathId(len(t.paths))
	t.paths[path] = newPathId
	p := PathRecord(path)
	t.Register(p)
	return newPathId
}

func (t *TraceRecord) EnsurePathId(path string) PathId {
	pathId, ok := t.paths[path]
	if !ok {
		pathId = t.RegisterPathWithNewId(path)
	}
	return pathId
}

func main() {
	record := MakeTraceRecord()

	if record.RegisterPathWithNewId("path0") != PathId(0) {
		panic("expected PathId 0 for path0")
	}
	if record.RegisterPathWithNewId("path1") != PathId(1) {
		panic("expected PathId 1 for path1")
	}
	record.RegisterStep(PathId(0), Line(1))
	record.RegisterCall("example", PathId(1), Line(1))
	record.RegisterReturn(NilValueRecord{})

	fmt.Println(record)
}
