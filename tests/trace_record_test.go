package trace_record_test						

import (
	"testing"
	// "fmt"

	ct "github.com/metacraft-labs/trace_record"
)


func TestTraceRecordHelpers(t *testing.T) {
	record := ct.MakeTraceRecord()

	workdir := "/home/alexander92/example_workdir/"
	path0 := workdir + "path0.txt"
	path1 := workdir + "path1.txt"
	if record.RegisterPathWithNewId(path0) != ct.PathId(0) {
		t.Errorf("expected PathId 0 for path0")
	}
	if record.RegisterPathWithNewId(path1) != ct.PathId(1) {
		t.Errorf("expected PathId 1 for path1")
	}

	record.RegisterCall("<toplevel>", path0, ct.Line(0))
	record.RegisterStep(path0, ct.Line(1))
	record.RegisterCall("example", path1, ct.Line(1))

	if record.RegisterTypeWithNewId("Int", ct.NewSimpleTypeRecord(ct.INT_TYPE_KIND, "Int")) != ct.TypeId(0) {
		t.Errorf("expected TypeId 0 for type Int")
	}

	record.RegisterReturn(ct.IntValue(1, ct.TypeId(0)))

	directory := "trace/"
	err := record.ProduceTrace(directory, workdir)
	if err != nil {
		t.Errorf("producing trace error: %v", err)
	}
	// fmt.Println(record)
}

