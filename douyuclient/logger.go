package douyuclient

import (
	"io"
	"log"
	"os"
)

var (
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "[DYClient]", log.LstdFlags|log.Lshortfile)
}

func startSpan(msg string) string {
	logger.Println("start", msg)
	return msg
}

func SetLoggerOut(out io.Writer) {
	logger = log.New(out, "[DYClient]", log.LstdFlags|log.Lshortfile)
}
