package main

import (
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
		offset:    c.String(flagOffset),
		timeout:   c.Duration(flagTimeout),
		delay:     c.Duration(flagDelay),
		skipDelay: c.Bool(flagSkipDelay),
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
			"auto.offset.reset":     client.offset,
			"group.id":              "test",
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create consumer: %s\n", err)
	}

	log.Printf("Created Consumer %v\n", kafkaConsumer)

	if err = kafkaConsumer.SubscribeTopics(client.topics, nil); err != nil {
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

			ev := kafkaConsumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("Message on %s: %s", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					log.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				errors = append(errors, fmt.Sprintf("error: %v: %v\n", e.Code(), e))
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				log.Printf("Ignored %v\n", e)
			}
		}

		if !client.skipDelay {
			time.Sleep(client.delay)
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
