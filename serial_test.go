package main

import (
	"strings"
	"testing"
	"time"
)

func _TestReadMessageStrings(t *testing.T) {
	msgStrCh := make(chan string, 0)
	//reader, err := os.Open("test_data.bin")

	//go MessageToPatchConverter()

	var str string

	select {
	case str = <-msgStrCh:
		// OK
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout…")
	}

	if len(str) < 21 || len(str) > 25 || len(strings.Split(str, ",")) != 6 || len(strings.Split(str, "*")) != 2 {
		t.Errorf("Wrong string sendt back, got '%s'", str)
	}
}
