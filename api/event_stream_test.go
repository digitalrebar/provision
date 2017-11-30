package api

import (
	"testing"

	"github.com/digitalrebar/provision/models"
)

func testFunc(t *testing.T, m interface{}, fn TestFunc, result, haveError bool) {
	t.Helper()

	b, e := fn(m)
	if haveError && e == nil {
		t.Errorf("Should have had an error: %v %v", b, e)
	} else if !haveError && e != nil {
		t.Errorf("Should not have had an error: %v %v", b, e)
	} else if !haveError && e == nil {
		if b != result {
			t.Errorf("Result mismatch: E,R (%v, %v)\n", result, b)
		}
	}
}

func TestEqualItem(t *testing.T) {
	machine := &models.Machine{CurrentTask: 3, Tasks: []string{"t1", "t2"}}

	fnRunnable := EqualItem("Runnable", true)
	fnCT := EqualItem("CurrentTask", 1)
	fnTask := EqualItem("Tasks", []string{"t4"})
	fnOr := OrItems(fnRunnable, fnCT)
	fnAnd := AndItems(fnRunnable, fnCT)
	fnNot := NotItem(fnRunnable)

	testFunc(t, machine, fnRunnable, false, false)
	testFunc(t, machine, fnNot, true, false)
	testFunc(t, machine, fnCT, false, false)
	testFunc(t, machine, fnTask, false, false)
	testFunc(t, machine, fnOr, false, false)
	testFunc(t, machine, fnAnd, false, false)

	machine.CurrentTask = 1
	machine.Runnable = true
	machine.Tasks = []string{"t4"}

	testFunc(t, machine, fnRunnable, true, false)
	testFunc(t, machine, fnNot, false, false)
	testFunc(t, machine, fnCT, true, false)
	testFunc(t, machine, fnTask, true, false)
	testFunc(t, machine, fnOr, true, false)
	testFunc(t, machine, fnAnd, true, false)

	machine.CurrentTask = 1
	machine.Runnable = false
	testFunc(t, machine, fnOr, true, false)
	testFunc(t, machine, fnAnd, false, false)

	machine.CurrentTask = 3
	machine.Runnable = true
	testFunc(t, machine, fnOr, true, false)
	testFunc(t, machine, fnAnd, false, false)
}
