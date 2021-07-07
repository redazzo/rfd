package main

import (
	"fmt"
	"log"
)

type Trace interface {
	traceLog(msg string)
}

type TraceLog struct {
	// In future, log this to a file

}

func CheckFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func (t TraceLog) traceLog(msg string) {
	fmt.Println(msg)
	log.Print(msg)
}
