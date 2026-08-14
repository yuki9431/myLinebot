package logger

import (
	"io"
	"log"
	"os"
)

type logger struct {
	out io.Writer
}

type Logger interface {
	Write(...interface{})
	Fatal(...interface{})
}

func New(w io.Writer) Logger {
	return &logger{out: w}
}

func (l *logger) Write(a ...interface{}) {
	log.SetOutput(l.out)
	log.Println(a...)
}

func (l *logger) Fatal(a ...interface{}) {
	l.Write(a...)
	os.Exit(1)
}
