package logger

import (
	"log"
	"os"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func Init(debug bool) {
	Info = log.New(os.Stdout, "[info] ", log.LstdFlags)
	Error = log.New(os.Stdout, "[error] ", log.LstdFlags)

	if debug {
		Debug = log.New(os.Stdout, "[debug] ", log.LstdFlags|log.Lshortfile)
	} else {
		Debug = log.New(os.Stdout, "", 0)
		Debug.SetOutput(os.Stderr)
	}
}
