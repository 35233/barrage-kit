package main

import (
	"context"
	"github.com/35233/barrage-kit/stt"
	"github.com/olivere/elastic"
	"strings"
	"time"
)

type ElasticSearchOutputFactory struct {
}

type elasticSearchOutput struct {
	esClient    *elastic.Client
	p           *elastic.BulkProcessor
	index       string
	docType     string
	messageChan chan *EsMessage
}

type eslogger byte

type EsMessage struct {
	messageTime int64
	data        string
}

func (eslogger) Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func init() {
	outputFactoryList = append(outputFactoryList, &ElasticSearchOutputFactory{})
}

func (factory *ElasticSearchOutputFactory) Type() string {
	return "elasticsearch"
}

func (factory *ElasticSearchOutputFactory) NewOutput(config *settingOutput) Output {
	logger.Println("ElasticSearchOutputFactory NewOutput", config)
	if len(strings.TrimSpace(config.Index)) == 0 {
		logger.Println("ElasticSearchOutputFactory Index can not empty")
		return nil
	}

	client, err := elastic.NewClient(
		elastic.SetInfoLog(eslogger(1)),
		elastic.SetURL(config.Urls...),
	)
	if err != nil {
		logger.Println("Error creating the client:", err)
		return nil
	}

	// Setup a bulk processor
	s := client.BulkProcessor().Name("EsBackgroundWorker-1").Workers(2)
	bulkActions := 1000
	if config.BulkActions != 0 {
		s.BulkActions(config.BulkActions)
		bulkActions = config.BulkActions
	}
	if config.BulkSize != 0 {
		s.BulkSize(config.BulkSize)
	}
	if config.FlushInterval != 0 {
		s.FlushInterval(config.FlushInterval * time.Millisecond)
	}
	p, err := s.Do(context.Background())
	if err != nil {
		logger.Println("Error creating the BulkProcessorService:", err)
		return nil
	}

	output := &elasticSearchOutput{
		client,
		p,
		config.Index,
		"type",
		make(chan *EsMessage, bulkActions*5),
	}
	go output.handleMessage()
	return output
}

func (output *elasticSearchOutput) Emit(messageTime int64, data string) {
	output.messageChan <- &EsMessage{messageTime, data}
}

func (output *elasticSearchOutput) handleMessage() {
	for message := range output.messageChan {
		msgObj := stt.Decode(message.data)
		msgMap, ok := msgObj.(map[string]interface{})
		if !ok {
			msgObj = map[string]interface{}{
				"data": msgObj,
			}
		}
		msgMap["@timestamp"] = message.messageTime / 1e6
		msgMap["rawMessage"] = message.data

		r := elastic.NewBulkIndexRequest().Index(output.index).Type(output.docType).Doc(msgMap)

		output.p.Add(r)
	}
}
