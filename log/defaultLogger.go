package log

import (
	"log"
)

type DefaultLogger struct {
}

func (logger DefaultLogger) Debug(args ...interface{}) {
	log.Println(args)
}

func (logger DefaultLogger) Info(args ...interface{}) {
	log.Println(args)
}

func (logger DefaultLogger) Warn(args ...interface{}) {
	log.Println(args)
}

func (logger DefaultLogger) Error(args ...interface{}) {
	log.Println(args)
}

func (logger DefaultLogger) Fatal(args ...interface{}) {
	log.Println(args)
}
