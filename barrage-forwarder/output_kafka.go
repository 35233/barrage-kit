package main

import (
	"encoding/json"
	"github.com/35233/barrage-kit/stt"
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

type KafkaOutputFactory struct {
}

type kafkaOutput struct {
	topic    string
	producer sarama.AsyncProducer
}

func init() {
	outputFactoryList = append(outputFactoryList, &KafkaOutputFactory{})
}

func (factory *KafkaOutputFactory) Type() string {
	return "kafka"
}

func (factory *KafkaOutputFactory) NewOutput(config *settingOutput) Output {
	logger.Println("KafkaOutputFactory NewOutput", config)
	broker := config.Brokers

	return &kafkaOutput{
		topic:    config.Topic,
		producer: newAccessLogProducer(strings.Split(broker, ",")),
	}
}

func (output *kafkaOutput) Emit(messageTime int64, data string) {
	jsonStr, err := json.Marshal(map[string]interface{}{
		"messageTime": messageTime,
		"data":        stt.Decode(data),
	})
	if err != nil {
		logger.Println("kafkaOutPut Marshal error", err)
		return
	}
	output.producer.Input() <- &sarama.ProducerMessage{
		Topic: output.topic,
		Value: sarama.ByteEncoder(jsonStr),
	}
}

func newAccessLogProducer(brokerList []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		logger.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			logger.Println("Failed to write access log entry:", err)
		}
	}()

	return producer
}
