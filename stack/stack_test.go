package stack

import (
	"fmt"
	"regexp"
	"testing"
)

type someT struct{}

func (s someT) level2() *Trace {
	return NewTrace(0)
}

func (s someT) level1() *Trace {
	return s.level2()
}

func (s someT) level0() *Trace {
	return s.level1()
}

func level2() *Trace {
	return NewTrace(0)
}

func level1() *Trace {
	return level2()
}

func level0() *Trace {
	return level1()
}

func assertFrame(t *testing.T, frame *Frame, file string, line int, fun string) {
	testName := fmt.Sprintf("[file: %s line: %d fun: %s]", file, line, fun)
	if !regexp.MustCompile(file).MatchString(frame.File) {
		t.Errorf("assertFrame: %s didn't match file in %v", testName, frame)
	}

	if frame.LineNumber != line {
		t.Errorf("assertFrame: %s didn't match line in %v", testName, frame)
	}

	if frame.Name != fun {
		t.Errorf("assertFrame: %s didn't match function name in %v", testName, frame)
	}
}
