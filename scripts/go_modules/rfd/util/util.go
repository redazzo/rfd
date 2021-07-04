package util

import (
	"fmt"
	"log"
)

func CheckFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func TraceLog(msg string) {
	fmt.Println(msg)
}
