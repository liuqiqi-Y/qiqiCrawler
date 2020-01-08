package util

import (
	"log"
	"os"
)
var Trace *log.Logger
func init() {
	Trace = log.New(os.Stderr, "Trace: ", log.Ldate|log.Ltime|log.Lshortfile)
}