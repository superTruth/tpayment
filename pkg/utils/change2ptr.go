package utils

import "time"

func ChangeTime2Ptr(t time.Timer) *time.Timer {
	return &t
}
