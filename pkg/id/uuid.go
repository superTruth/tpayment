package id

import "sync"

var factory *SnowFlake
var initOnce sync.Once

func New() uint64 {
	initOnce.Do(func() {
		factory, _ = NewSnowFlake(0, 0)
	})
	i, _ := factory.Next()
	return i
}
