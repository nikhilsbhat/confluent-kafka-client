package main

import (
	"time"
)

type kafkaClient struct {
	serverUrl string
	topic     string
	topics    []string
	offset    string
	timeout   time.Duration
	delay     time.Duration
	skipDelay bool
}
