package id

import (
	"strconv"
)

type UUID uint64

var IDGenr *SnowFlake

func init() {
	IDGenr, _ = NewSnowFlake(0, 0)
}

func New() uint64 {
	i, _ := IDGenr.Next()
	return i
}

func (c UUID) String() string {
	return strconv.FormatUint(uint64(c), 10)
	// return string(c)
}
