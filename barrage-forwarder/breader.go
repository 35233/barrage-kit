package main

import (
	"fmt"
	"github.com/35233/barrage-kit/douyuclient"
	"log"
	"os"
)

var logger *log.Logger
var outputFactoryList []OutputFactory

func init() {
	outputFactoryList = make([]OutputFactory, 0, 5)
}

func initOutputs(config []settingOutput) []Output {
	outputs := make([]Output, 0)
	for _, settingOutput := range config {
		isFound := false
		for _, outputFactory := range outputFactoryList {
			if settingOutput.Type == outputFactory.Type() {
				if out := outputFactory.NewOutput(&settingOutput); out != nil {
					outputs = append(outputs, out)
				}
				isFound = true
				break
			}
		}
		if !isFound {
			logger.Printf("Unknown output type: %s\n", settingOutput.Type)
		}
	}
	return outputs
}

func main() {
	logger = log.New(os.Stdout, "[main]", log.LstdFlags|log.Lshortfile)
	logger.Println("starting")
	setting, err := loadSetting()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	f, err := os.OpenFile(setting.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	logger = log.New(f, "[main]", log.LstdFlags|log.Lshortfile)

	douyuclient.SetLoggerOut(f)
	dyClient := douyuclient.New("openbarrage.douyutv.com:8601", 50)
	if len(setting.Sources) <= 0 {
		logger.Println("Sources empty")
		return
	}
	for _, ss := range setting.Sources {
		if len(ss.RoomIds) <= 0 {
			logger.Println("RoomIds empty")
			return
		}
		if ss.Type != "douyu" {
			logger.Printf("Unknown source type: %s\n", ss.Type)
			return
		}
		for _, roomId := range ss.RoomIds {
			dyClient.AddRoom(roomId)
		}
	}
	msgChannel, err := dyClient.Start()
	if err != nil {
		logger.Println(err)
		return
	}

	outputs := initOutputs(setting.Output)
	logger.Println("started")
	for msg := range msgChannel {
		text := msg.Text()
		for _, out := range outputs {
			out.Emit(msg.NanoTimestamp, text)
		}
	}
}
