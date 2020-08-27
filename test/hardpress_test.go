package test

import (
	"testing"
	"time"
)

func TestHardPress(t *testing.T) {

	for i := 0; i < 1000; i++ {
		go TestLoop(t)
	}

	time.Sleep(time.Minute * 30)
}

func TestLoop(t *testing.T) {

	for {
		TestHeartBeat(t)
	}

}
