package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func SetUuidNode(i int) {
	newUuid().setNodeId(i)
}

func GenStrUuid() string {
	for i := 0; i < 3; i++ {
		uuid, err := newUuid().generateId()
		if err == nil {
			return fmt.Sprintf("%d", uuid)
		}
	}
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

func GenUuid() int64 {
	for i := 0; i < 3; i++ {
		uuid, err := newUuid().generateId()
		if err == nil {
			return uuid
		}
	}
	return time.Now().UnixMilli()
}

func UuidToTime(uuid int64) time.Time {
	return newUuid().uuidToTime(uuid)
}

func TimeToUuid(ts time.Time) int64 {
	return newUuid().timeToUuid(ts)
}

var (
	rotateId       = 15
	nodeId         = 6
	T2018    int64 = 1514736000000
)

type Uuid struct {
	lock          sync.Mutex
	timeId        int64
	nodeId        int
	rotateId      int
	nodeIdWidth   int
	nodeIdMask    int
	rotateIdWidth int
	rotateIdMask  int
	T201801010000 int64
}

var midUtil *Uuid

func newUuid() *Uuid {
	if midUtil == nil {
		midUtil = &Uuid{
			timeId:        0,
			lock:          sync.Mutex{},
			nodeIdWidth:   nodeId,
			nodeIdMask:    0x3F,
			rotateIdWidth: rotateId,
			rotateIdMask:  0x7FFF,
			T201801010000: T2018,
			nodeId:        0,
		}
	}
	return midUtil
}

func (m *Uuid) setNodeId(i int) {
	m.nodeId = i
	return
}

func (m *Uuid) generateId() (int64, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.rotateId = (m.rotateId + 1) & m.rotateIdMask
	mid := time.Now().UnixMilli() - m.T201801010000
	if mid > m.timeId {
		m.timeId = mid
		m.rotateId = 1
	} else if mid == m.timeId {
		if m.rotateId == (m.rotateIdMask - 1) {
			for {
				if mid <= m.timeId {
					mid = time.Now().UnixMilli() - m.T201801010000
				} else {
					return m.generateId()
				}
			}
		}
	} else {
		if m.rotateId > (m.rotateIdMask-1)*9/10 {
			if m.timeId-mid < 3000 {
				for {
					if mid < m.timeId {
						mid = time.Now().UnixMilli() - m.T201801010000
					}
					return m.generateId()
				}
			} else {
				return 0, errors.New(fmt.Sprintf("Time turn back %d ms, it too long!!!", m.timeId-mid))
			}
		} else {
			mid = m.timeId
		}
	}
	mid <<= m.nodeIdWidth
	mid += int64(m.nodeId & m.nodeIdMask)
	mid <<= m.rotateIdWidth
	mid += int64(m.rotateId)
	return mid, nil
}

func (m *Uuid) uuidToTime(uuid int64) time.Time {
	uuid >>= m.nodeId + m.rotateId
	dateNum := uuid + T2018
	return time.UnixMilli(dateNum)
}

func (m *Uuid) timeToUuid(ts time.Time) int64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.rotateId = (m.rotateId + 1) & m.rotateIdMask
	mid := ts.UnixMilli() - m.T201801010000
	mid <<= m.nodeIdWidth
	mid += int64(m.nodeId & m.nodeIdMask)
	mid <<= m.rotateIdWidth
	mid += int64(m.rotateId)
	return mid
}
