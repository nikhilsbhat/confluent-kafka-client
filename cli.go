package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

const (
	flagServerURL              = "server-url"
	flagTopic                  = "topic"
	flagTopics                 = "topics"
	flagTimeout                = "timeout"
	flagDelay                  = "delay"
	flagSkipDelay              = "skip-delay"
	flagOffset                 = "offset"
	flagGroupID                = "group-id"
	durationFlagTimeMultiplier = 30
	delayFlagTimeMultiplier    = 100
)

func cliApp() *cli.App {
	return &cli.App{
		Name:                 "confluent-kafka-client",
		Usage:                "Utility to test slow consumer behaviour",
		UsageText:            "confluent-kafka-client [flags]",
		EnableBashCompletion: true,
		HideHelp:             false,
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "version of the confluent-test",
				Action:  appVersion,
			},
			registerProducer(),
			registerConsumer(),
		},
	}
}

func registerProducer() *cli.Command {
	return &cli.Command{

		Name:            "producer",
		Aliases:         []string{"p"},
		Usage:           "publishes messages to topic until timer lapses",
		Action:          initProducer,
		HideHelp:        false,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagServerURL,
				Usage:    "kafka server url",
				Aliases:  []string{"s"},
				EnvVars:  []string{"SERVER_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    flagTopic,
				Usage:   "name of topic to which message has to be published",
				EnvVars: []string{"TOPIC"},
				Aliases: []string{"t"},
			},
			&cli.DurationFlag{
				Name:    flagTimeout,
				Usage:   "time duration until when the messages to be published to specified topic",
				EnvVars: []string{"TIMEOUT"},
				Aliases: []string{"to"},
				Value:   time.Second * durationFlagTimeMultiplier,
			},
			&cli.DurationFlag{
				Name:    flagDelay,
				Usage:   "delay to be introduced between every message that is published",
				EnvVars: []string{"DELAY"},
				Aliases: []string{"d"},
				Value:   time.Millisecond * delayFlagTimeMultiplier,
			},
		},
	}
}

func registerConsumer() *cli.Command {
	return &cli.Command{
		Name:            "consumer",
		Aliases:         []string{"c"},
		Usage:           "subscribe to the topic with the specified delay in between",
		Action:          initConsumer,
		HideHelp:        false,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagServerURL,
				Usage:    "kafka server url",
				Aliases:  []string{"s"},
				EnvVars:  []string{"SERVER_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    flagOffset,
				Usage:   "The offset id to consume from",
				Aliases: []string{"o"},
				EnvVars: []string{"OFFSET"},
				Value:   "earliest",
			},
			&cli.StringSliceFlag{
				Name:    flagTopics,
				Usage:   "lists of topic to which message has to be published",
				EnvVars: []string{"TOPICS"},
				Aliases: []string{"ts"},
			},
			&cli.StringFlag{
				Name:    flagTopic,
				Usage:   "name of topic to which message has to be published",
				EnvVars: []string{"TOPIC"},
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:    flagGroupID,
				Usage:   "a unique string that identifies the consumer group this consumer belongs to",
				EnvVars: []string{"GROUP_ID"},
				Aliases: []string{"g"},
				Value:   "test-group",
			},
			&cli.BoolFlag{
				Name:    flagSkipDelay,
				Usage:   "enable this flag to skip delay between the consumer call",
				EnvVars: []string{"SKIP_DELAY"},
				Value:   false,
			},
			&cli.DurationFlag{
				Name:    flagTimeout,
				Usage:   "timeout while waiting for receiving message",
				EnvVars: []string{"TIMEOUT"},
				Aliases: []string{"to"},
				Value:   time.Second * durationFlagTimeMultiplier,
			},
			&cli.DurationFlag{
				Name:    flagDelay,
				Usage:   "delay to be introduced between every subscriber poll",
				EnvVars: []string{"DELAY"},
				Aliases: []string{"d"},
				Value:   time.Second * durationFlagTimeMultiplier,
			},
		},
	}
}
