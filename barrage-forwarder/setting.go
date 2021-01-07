package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type setting struct {
	LogPath string `yaml:"log.path"`
	Sources []struct {
		Type    string `yaml:"type"`
		RoomIds []string
	}
	Output []settingOutput
}

type settingOutput struct {
	Type                        string `yaml:"type"`
	Path, Index, Brokers, Topic string
	BulkActions                 int           `yaml:"bulkActions"`
	BulkSize                    int           `yaml:"bulkSize"`
	FlushInterval               time.Duration `yaml:"flushInterval"`
	Urls                        []string      `yaml:"urls"`
}

func loadSetting() (*setting, error) {
	pp, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Println("current directory:", pp)

	if len(os.Args) < 2 {
		return nil, errors.New("Please set a config file name.")
	}
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	buf := make([]byte, 1024)
	m := bytes.Buffer{}
	for {
		n, err := f.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		m.Write(buf[:n])
	}

	t := setting{}

	err = yaml.Unmarshal(m.Bytes(), &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
