package main

import (
	"time"
)

type kafkaClient struct {
	serverURL string
	topic     string
	topics    []string
	offset    string
	timeout   time.Duration
	delay     time.Duration
	skipDelay bool
	groupID   string
}
