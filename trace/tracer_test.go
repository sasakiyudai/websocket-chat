package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("New: nil error")
	} else {
		tracer.Trace("hello from trace pkg!")
		if buf.String() != "hello from trace pkg!\n" {
			t.Errorf("wrong string '%s' was output", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("data")
}