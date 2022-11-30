package producers

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Data struct {
	Id          int
	Message     string
	Source      string
	Destination string
}

func Produce(p *kafka.Producer, messages string, topicName string) error {
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
		time.Sleep(time.Second * 5) //wait for 5 seconds before sending another batch of messages
	}()

	topic := topicName
	data := []byte(messages)
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
