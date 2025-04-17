package trace_record

import "fmt"
import "bytes"
import "os"
import "path/filepath"
import "encoding/json"

type FunctionId uint64
type CallId uint64
type VariableId uint64
type StepId uint64
type PathId uint64
type Line int64

type RecordEvent interface {
	isRecordEvent() bool
	MarshalJson() ([]byte, error)
}

/// steps
type StepRecord struct {
	PathId PathId `json:"pathId"`
	Line   Line   `json:"line"`
}

func (s StepRecord) isRecordEvent() bool {
	return true
}

type RawStepRecord struct {
	Step StepRecord
}

func (receiver StepRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(RawStepRecord { receiver })
}

// functionrrecords

type FunctionRecord struct {
	Name   string `json:"name"`
	PathId PathId `json:"pathId"`
	Line   Line   `json:"line"`
}

func (r FunctionRecord) isRecordEvent() bool {
	return true
}

type RawFunctionRecord struct {
	Function FunctionRecord
}

func (receiver FunctionRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(RawFunctionRecord { receiver })
}


type ArgRecord struct {
	Name  string       `json:"name"`
	Value ValueRecord  `json:"value"`
}

func (receiver ArgRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(receiver)
}

type CallRecord struct {
	FunctionId FunctionId    `json:"functionId"`
	Args       []ArgRecord   `json:"args"`
}

func (c CallRecord) isRecordEvent() bool {
	return true
}

type RawCallRecord struct {
	Call CallRecord
}


func (receiver CallRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(RawCallRecord { receiver })
}

// ====

type ReturnRecord struct {
	ReturnValue ValueRecord    `json:"return_value"`
}

type RawReturnRecord struct {
	Return ReturnRecord
}

func (r ReturnRecord) isRecordEvent() bool {
	return true
}

func (receiver ReturnRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(RawReturnRecord { receiver })
}

// ====

type PathRecord string

type RawPathRecord struct {
	Path PathRecord
}

func (p PathRecord) isRecordEvent() bool {
	return true
}

func (receiver PathRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(RawPathRecord { receiver })
}

// ====

type RawTypeRecord struct {
	Type TypeRecord
}

func (r RawTypeRecord) isRecordEvent() bool {
	return true
}

func (receiver RawTypeRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(receiver)
}

type TraceRecord struct {
	events    []RecordEvent
	functions map[string]FunctionId
	paths     map[string]PathId
	types     map[string]TypeId
}

func MakeTraceRecord() TraceRecord {
	events := make([]RecordEvent, 0)
	functions := make(map[string]FunctionId, 0)
	paths := make(map[string]PathId, 0)
	types := make(map[string]TypeId, 0)
	return TraceRecord{ events, functions, paths, types }
}

func (t *TraceRecord) Register(event RecordEvent) {
	t.events = append(t.events, event)
}

func (t *TraceRecord) RegisterStepWithPathId(pathId PathId, line Line) {
	step := StepRecord{pathId, line}
	t.Register(step)
}

func (t *TraceRecord) RegisterStep(path string, line Line) {
	pathId := t.EnsurePathId(path)
	t.RegisterStepWithPathId(pathId, line)
}

// naming copied from DWARF: definition path and definition line
func (t *TraceRecord) RegisterCall(name string, definitionPath string, definitionLine Line) {
	definitionPathId := t.EnsurePathId(definitionPath)
	t.RegisterCallWithPathId(name, definitionPathId, definitionLine)
}

func (t *TraceRecord) RegisterCallWithPathId(name string, definitionPathId PathId, definitionLine Line) {
	functionId := t.EnsureFunctionId(name, definitionPathId, definitionLine)
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


func (t *TraceRecord) RegisterTypeWithNewId(name string, typeRecord TypeRecord) TypeId {
	newTypeId := TypeId(len(t.types))
	t.types[name] = newTypeId
	t.Register(RawTypeRecord { typeRecord })
	return newTypeId
}

func (record *TraceRecord) SerializeEventsToJson() ([]byte, error) {
	var jsonEvents bytes.Buffer
	jsonEvents.WriteString("[\n")
	for i, event := range record.events {
		raw, err := event.MarshalJson()
		if err != nil {
			var empty []byte
			return empty, err
		} else {
			text := string(raw[:])
			jsonEvents.WriteString("    ")
			jsonEvents.WriteString(text)
			if i < len(record.events) - 1 {
				jsonEvents.WriteString(",\n")
			} else {
				jsonEvents.WriteString("\n")
			}
		}
	}
	jsonEvents.WriteString("]\n")
	// jsonText := jsonEvents.String()
	jsonBytes := jsonEvents.Bytes()
	return jsonBytes, nil
	// fmt.Println(jsonText)
	
}


func (record *TraceRecord) ProduceTrace(directory string) error { 
	// TODO : augment errors, instead of printing

	jsonBytes, err := record.SerializeEventsToJson()
	if err != nil {
		return err
	}

	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println("error: couldn't ensure trace directory exists or make it: ", err)
		return err
	}

	err = os.WriteFile(filepath.Join(directory, "trace.json"), jsonBytes, 0644)
	if err != nil {
		fmt.Println("error: couldn't write trace.json: ", err)
		return err
	}

	fmt.Println("generated trace in ", directory)
	return nil
}
