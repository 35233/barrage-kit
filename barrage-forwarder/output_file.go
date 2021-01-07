package main

import (
	"fmt"
	"os"
	"time"
)

type FileOutputFactory struct {
}

type fileOutput struct {
	file *os.File
}

func init() {
	outputFactoryList = append(outputFactoryList, &FileOutputFactory{})
}

func (factory *FileOutputFactory) Type() string {
	return "file"
}

func (factory *FileOutputFactory) NewOutput(config *settingOutput) Output {
	logger.Println("FileOutputFactory NewOutput", config)
	path := config.Path
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Println("OpenFile error", err)
		return nil
	}
	return &fileOutput{
		file: f,
	}
}

func (output *fileOutput) Emit(messageTime int64, data string) {
	currentTime := time.Unix(messageTime/1e9, 0)
	if _, err := output.file.WriteString(fmt.Sprintf("[%s]%s\n", currentTime.Format(time.RFC3339), data)); err != nil {
		logger.Println("fileOutput.Emit error", err)
	}
}
