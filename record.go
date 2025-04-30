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
	PathId PathId `json:"path_id"`
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
	PathId PathId `json:"path_id"`
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
	FunctionId FunctionId    `json:"function_id"`
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

/// ===

type VariableNameRecord struct {
	VariableName string
}

func (v VariableNameRecord) isRecordEvent() bool {
	return true
}

func (receiver VariableNameRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(receiver)
}

// ===

type FullValueRecord struct {
	VariableId VariableId `json:"variable_id"`
	Value ValueRecord `json:"value"`
}

type RawValueRecord struct {
	Value FullValueRecord
}
func (r FullValueRecord) isRecordEvent() bool {
	return true
}

func (receiver FullValueRecord) MarshalJson() ([]byte, error) {
	r := RawValueRecord{receiver}
	return json.Marshal(r)
}

// ====

type RecordEventKind int

const (
	EventKindWrite RecordEventKind = iota
	EventKindWriteFile
	EventKindWriteOther 
	EventKindRead
	EventKindReadFile
	EventKindReadOther
	// not generated yet in most recorders
	EventKindReadDir
	EventKindOpenDir
	EventKindCloseDir
	EventKindSocket
	EventKindOpen
	// errors/exceptions/signals
	EventKindError
	// used for trace events
	EventKindTraceLogEvent
	// TODO others
)

type RecordEventRecord struct {
	Kind RecordEventKind `json:"kind"`
	Metadata string `json:"metadata"`
	Content string `json:"content"`
}

type RawRecordEventRecord struct {
	Event RecordEventRecord
}

func (r RecordEventRecord) isRecordEvent() bool {
	return true
}

func (receiver RecordEventRecord) MarshalJson() ([]byte, error) {
	return json.Marshal(RawRecordEventRecord { receiver })
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
	variables map[string]VariableId
	types     map[string]TypeId
}

func MakeTraceRecord() TraceRecord {
	events := make([]RecordEvent, 0)
	functions := make(map[string]FunctionId, 0)
	paths := make(map[string]PathId, 0)
	variables := make(map[string]VariableId, 0)
	types := make(map[string]TypeId, 0)
	return TraceRecord{ events, functions, paths, variables, types }
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

func (t *TraceRecord) RegisterVariableNameWithNewId(name string) VariableId {
	r := VariableNameRecord{name}
	t.Register(r)
	newVariableId := VariableId(len(t.variables))
	t.variables[name] = newVariableId
	return newVariableId
}

func (t *TraceRecord) EnsureVariableId(name string) VariableId {
	variableId, ok := t.variables[name]
	if !ok {
		variableId = t.RegisterVariableNameWithNewId(name)
	}
	return variableId
}

func (t *TraceRecord) RegisterFullValue(variableId VariableId, value ValueRecord) {
	r := FullValueRecord{variableId, value}
	t.Register(r)
}

func (t *TraceRecord) RegisterVariable(name string, value ValueRecord) {
	variableId := t.EnsureVariableId(name)
	t.RegisterFullValue(variableId, value)
}

func (t *TraceRecord) RegisterRecordEvent(kind RecordEventKind, metadata string, content string) {
	event := RecordEventRecord {kind, metadata, content}
	t.Register(event)
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

func (t *TraceRecord) EnsureTypeId(name string, typeRecord TypeRecord) TypeId {
	typeId, ok := t.types[name]
	if !ok {
		typeId = t.RegisterTypeWithNewId(name, typeRecord)
	}
	return typeId
}

type TraceMetadata struct {
	Workdir string `json:"workdir"`
	Program string `json:"program"`
	Args []string `json:"args"`
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


func (record *TraceRecord) ProduceTrace(traceDirectory string, programName string, workdir string) error { 
	// TODO : augment errors, instead of printing

	jsonBytes, err := record.SerializeEventsToJson()
	if err != nil {
		return err
	}

	err = os.MkdirAll(traceDirectory, os.ModePerm)
	if err != nil {
		fmt.Println("error: couldn't ensure trace directory exists or make it: ", err)
		return err
	}

	err = os.WriteFile(filepath.Join(traceDirectory, "trace.json"), jsonBytes, 0644)
	if err != nil {
		fmt.Println("error: couldn't write trace.json: ", err)
		return err
	}

	var args []string = make([]string, 0)
	traceMetadata := TraceMetadata {workdir, programName, args }
	traceMetadataJson, err := json.Marshal(traceMetadata)
	if err != nil {
		fmt.Println("error: encoding trace metadata: ", err)
		return err
	}
	// fmt.Println(traceMetadataJson)
	err = os.WriteFile(filepath.Join(traceDirectory, "trace_metadata.json"), traceMetadataJson, 0644)
	if err != nil {
		fmt.Println("error: couldn't write trace_metadata.json: ", err)
		return err
	}

	paths := make([]string, 0)
	for _, event := range record.events {
		switch event.(type) {
		case PathRecord:
			pathRecord, _ := event.(PathRecord) // fmt.Fprint("%v", event)
			path := string(pathRecord)
			paths = append(paths, path)
		default:
			// nothing
		}
	}
	pathsJson, err := json.Marshal(paths)
	if err != nil {
		fmt.Println("error: encoding trace paths: ", err)
		return err
	}
	err = os.WriteFile(filepath.Join(traceDirectory, "trace_paths.json"), pathsJson, 0644)
	if err != nil {
		fmt.Println("error: couldn't write trace_paths.json: ", err)
		return err
	}

	fmt.Println("generated trace in ", traceDirectory)
	return nil
}
