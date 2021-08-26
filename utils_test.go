package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getBytes(t *testing.T) {
	t.Run("should return the expected json", func(t *testing.T) {
		randomMessage := "nikhilbhat"
		expected := `{"message":"nikhilbhat","instance":1}`

		actual := getBytes(1, randomMessage)
		assert.Equal(t, expected, string(actual))
	})
}

func Test_randStringBytes(t *testing.T) {
	t.Run("", func(t *testing.T) {
		expected := randStringBytes(20)
		actual := randStringBytes(20)
		assert.Equal(t, expected, actual)
	})
}