package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/urfave/cli/v2"
)

func initProducer(c *cli.Context) error {
	newKafkaClient := &kafkaClient{
		serverUrl: c.String(flagServerUrl),
		topic:     c.String(flagTopic),
		timeout:   c.Duration(flagTimeout),
		delay:     c.Duration(flagDelay),
	}
	if err := newKafkaClient.producer(); err != nil {
		return err
	}
	return nil
}

func (client *kafkaClient) producer() error {
	kafkaProducer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": client.serverUrl,
		},
	)
	if err != nil {
		return fmt.Errorf("creating producer errored with: %s", err)
	}

	timeout := time.After(client.timeout)
	finish := make(chan bool)
	count := 1
	var errors []string
	go func() {
		for {
			select {
			case <-timeout:
				fmt.Println("timer lapsed, stopped producer")
				finish <- true
				return
			default:
				deliveryChan := make(chan kafka.Event)
				if err := kafkaProducer.Produce(
					&kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic: &client.topic, Partition: kafka.PartitionAny},
						Value:          getBytes(count),
					},
					deliveryChan); err != nil {
				}

				e := <-deliveryChan
				m := e.(*kafka.Message)

				if m.TopicPartition.Error != nil {
					errors = append(errors, fmt.Sprintf("delivery failed: %v\n", m.TopicPartition.Error))
				} else {
					log.Printf("delivered message %s of instance: %d to topic %s [%d] at offset %v\n",
						string(m.Value), count, *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}

				close(deliveryChan)
				count++
			}
			time.Sleep(client.delay)
		}
	}()

	<-finish

	if len(errors) != 0 {
		return fmt.Errorf("publishing message to topic %s failed with below errors: %v", client.topic, errors)
	}
	return nil
}