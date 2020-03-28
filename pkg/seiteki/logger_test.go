package seiteki

import (
	"fmt"
	"testing"
)

type testLogger struct {
	buff string
	typ  string
}

func (l *testLogger) CheckBuffAndFlush(t *testing.T, exB, exT string) {
	if l.buff != exB {
		t.Errorf("unexpected content: expected '%s', was '%s'",
			l.buff, exB)
	}

	if l.typ != exT {
		t.Errorf("unexpected type: expected '%s', was '%s'",
			l.typ, exT)
	}

	l.buff = ""
	l.typ = ""
}

func (l *testLogger) Error(args ...interface{}) {
	l.buff = fmt.Sprint(args...)
	l.typ = "error"
}

func (l *testLogger) Errorf(format string, args ...interface{}) {
	l.buff = fmt.Sprintf(format, args...)
	l.typ = "error"
}

func (l *testLogger) Fatal(args ...interface{}) {
	l.buff = fmt.Sprint(args...)
	l.typ = "fatal"
}

func (l *testLogger) Fatalf(format string, args ...interface{}) {
	l.buff = fmt.Sprintf(format, args...)
	l.typ = "fatal"
}

func (l *testLogger) Info(args ...interface{}) {
	l.buff = fmt.Sprint(args...)
	l.typ = "info"
}

func (l *testLogger) Infof(format string, args ...interface{}) {
	l.buff = fmt.Sprintf(format, args...)
	l.typ = "info"
}

func (l *testLogger) Warning(args ...interface{}) {
	l.buff = fmt.Sprint(args...)
	l.typ = "warn"
}

func (l *testLogger) Warningf(format string, args ...interface{}) {
	l.buff = fmt.Sprintf(format, args...)
	l.typ = "warn"
}

func TestNewLoggerWrapper(t *testing.T) {
	lw := newLogegrWrapper(new(testLogger))
	if lw == nil {
		t.Fatalf("new loggerWrapper was nil")
	}
	if lw.logger == nil {
		t.Fatalf("loggerWrapper inner was nil")
	}
}

func TestError(t *testing.T) {
	tl := new(testLogger)
	lw := newLogegrWrapper(tl)

	lw.Error(123, []byte("test"))
	tl.CheckBuffAndFlush(t, "123 [116 101 115 116]", "error")
}

func TestErrorf(t *testing.T) {
	tl := new(testLogger)
	lw := newLogegrWrapper(tl)

	lw.Errorf("%d, %s", 123, []byte("test"))
	tl.CheckBuffAndFlush(t, "123, test", "error")
}

func TestFatal(t *testing.T) {
	tl := new(testLogger)
	lw := newLogegrWrapper(tl)

	lw.Fatal(123, []byte("test"))
	tl.CheckBuffAndFlush(t, "123 [116 101 115 116]", "fatal")
}

func TestFatalf(t *testing.T) {
	tl := new(testLogger)
	lw := newLogegrWrapper(tl)

	lw.Fatalf("%d, %s", 123, []byte("test"))
	tl.CheckBuffAndFlush(t, "123, test", "fatal")
}

func TestInfo(t *testing.T) {
	tl := new(testLogger)
	lw := newLogegrWrapper(tl)

	lw.Info(123, []byte("test"))
	tl.CheckBuffAndFlush(t, "123 [116 101 115 116]", "info")
}

func TestInfof(t *testing.T) {
	tl := new(testLogger)
	lw := newLogegrWrapper(tl)

	lw.Infof("%d, %s", 123, []byte("test"))
	tl.CheckBuffAndFlush(t, "123, test", "info")
}

func TestWarning(t *testing.T) {
	tl := new(testLogger)
	lw := newLogegrWrapper(tl)

	lw.Warning(123, []byte("test"))
	tl.CheckBuffAndFlush(t, "123 [116 101 115 116]", "warn")
}

func TestWarningf(t *testing.T) {
	tl := new(testLogger)
	lw := newLogegrWrapper(tl)

	lw.Warningf("%d, %s", 123, []byte("test"))
	tl.CheckBuffAndFlush(t, "123, test", "warn")
}
