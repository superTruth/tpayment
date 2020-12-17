package id

/*
github.com/twitter/snowflake in golang

id =>  timestamp retain center worker sequence
          40      4       5      5      10

*/

import (
	"fmt"
	"sync"
	"time"
)

// var partnerId uint32

func init() {

}

const (
	nano = 1000 * 1000
)

const (
	RetainedBits = 4
	CenterBits   = 5
	MaxCenter    = -1 ^ (-1 << CenterBits) // center mask
	WorkerBits   = 5
	MaxWorker    = -1 ^ (-1 << WorkerBits)   // worker mask
	SequenceBits = 10                        // sequence
	MaxSequence  = -1 ^ (-1 << SequenceBits) // sequence mask
)

var (
	Since = time.Date(2017, 5, 1, 0, 0, 0, 0, time.Local).UnixNano() / nano
	//poolMu = sync.RWMutex{}
	//pool   = make(map[uint64]*SnowFlake)
)

type SnowFlake struct {
	lastTimestamp uint64
	retain        uint32
	center        uint32
	worker        uint32
	sequence      uint32
	lock          sync.Mutex
}

func (sf *SnowFlake) uint64() uint64 {
	return (sf.lastTimestamp << (RetainedBits + CenterBits + WorkerBits + SequenceBits)) |
		(uint64(sf.retain) << (CenterBits + WorkerBits + SequenceBits)) |
		(uint64(sf.center) << (WorkerBits + SequenceBits)) |
		(uint64(sf.worker) << SequenceBits) |
		uint64(sf.sequence)
}

func (sf *SnowFlake) Next() (uint64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	ts := timestamp()
	if ts == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & MaxSequence
		if sf.sequence == 0 {
			ts = tilNextMillis(ts)
		}
	} else {
		sf.sequence = 0
	}

	if ts < sf.lastTimestamp {
		return 0, fmt.Errorf("Invalid timestamp: %v - precedes %v", ts, sf)
	}
	sf.lastTimestamp = ts
	return sf.uint64(), nil
}

func NewSnowFlake(centerID uint32, workerID uint32) (*SnowFlake, error) {
	if centerID > MaxCenter {
		return nil, fmt.Errorf("CenterID %v is invalid", centerID)
	} else if workerID > MaxWorker {
		return nil, fmt.Errorf("WorkerID %v is invalid", workerID)
	}
	return &SnowFlake{
		worker: workerID,
		center: centerID,
	}, nil
}

func timestamp() uint64 {
	return uint64(time.Now().UnixNano()/nano - Since)
}

func tilNextMillis(ts uint64) uint64 {
	i := timestamp()
	for i <= ts {
		i = timestamp()
	}
	return i
}
