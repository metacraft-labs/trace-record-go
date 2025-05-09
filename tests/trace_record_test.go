package trace_record_test

import (
	"fmt"
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

	intTypeRecord := ct.NewSimpleTypeRecord(ct.INT_TYPE_KIND, "Int")
	if record.RegisterTypeWithNewId("Int", intTypeRecord) != ct.TypeId(0) {
		t.Errorf("expected TypeId 0 for type Int")
	}

	record.RegisterCall("<toplevel>", path0, ct.Line(0), []ct.FullValueRecord{})
	record.RegisterStep(path0, ct.Line(1))
	record.RegisterVariable("a", ct.IntValue(1, ct.TypeId(0)))

	arg1 := record.Arg("a", ct.IntValue(1, ct.TypeId(0)))
	record.RegisterCall("example", path1, ct.Line(1), []ct.FullValueRecord{arg1})

	record.RegisterRecordEvent(ct.EventKindWriteOther, "write_bytes: #52", "0000x")

	record.RegisterReturn(ct.IntValue(1, record.EnsureTypeId("Int", intTypeRecord)))

	callsCount := record.CurrentCallsCount()
	if callsCount != 2 {
		t.Errorf("expected 2 calls, not %d", callsCount)
	}

	directory := "trace/"
	err := record.ProduceTrace(directory, "test_program", workdir)
	if err != nil {
		t.Errorf("producing trace error: %v", err)
	}
	fmt.Println(record)
}
