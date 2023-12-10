package main

import (
	"testing"
)

func TestNtpServerIsCorrect(t *testing.T) {
	ntpServer = "Hello, Gopher!"
	err := getTime()
	if err == nil {
		t.Error(err)
	}
}
