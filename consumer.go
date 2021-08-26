package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/urfave/cli/v2"
)

func initConsumer(c *cli.Context) error {
	newKafkaClient := &kafkaClient{
		serverUrl: c.String(flagServerUrl),
		topics:    c.StringSlice(flagTopics),
		topic:     c.String(flagTopic),
		offset:    c.String(flagOffset),
		timeout:   c.Duration(flagTimeout),
		delay:     c.Duration(flagDelay),
		skipDelay: c.Bool(flagSkipDelay),
		groupID:   c.String(flagGroupID),
	}
	if err := newKafkaClient.consumer(); err != nil {
		return err
	}
	return nil
}

func (client *kafkaClient) consumer() error {
	kafkaConsumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":     client.serverUrl,
			"broker.address.family": "v4",
			"session.timeout.ms":    6000,
			"queued.min.messages":   5,
			"auto.offset.reset":     client.offset,
			"group.id":              client.groupID,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create consumer: %s\n", err)
	}

	log.Printf("Created Consumer: %v\n", kafkaConsumer)

	if err = kafkaConsumer.Subscribe(client.topic, nil); err != nil {
		return err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	var errors []string
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:

			if !client.skipDelay {
				time.Sleep(client.delay)
			}

			msg, err := kafkaConsumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				log.Println(err)
				continue
			}

			recordKey := string(msg.Key)
			recordValue := msg.Value
			data := fakeMessage{}
			err = json.Unmarshal(recordValue, &data)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err))
				continue
			}

			log.Printf("Consumed record with key %s and value %s on %v\n", recordKey, recordValue, msg.TopicPartition)
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("subscribing to a topic %s failed with below errors: %v", client.topic, errors)
	}

	log.Printf("closing consumer\n")
	if err := kafkaConsumer.Close(); err != nil {
		return err
	}

	return nil
}
