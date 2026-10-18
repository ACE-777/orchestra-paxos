package storage

import (
	"orchestra-paxos/internal/domain/roles"
	"sync"
)

type Storage struct {
	values map[roles.HighestID]string
	lock   *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		values: make(map[roles.HighestID]string),
		lock:   &sync.RWMutex{},
	}
}

func (s *Storage) AddResult(value string, operationID roles.HighestID) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.values[operationID] = value
}

func (s *Storage) CheckResult() uint64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return uint64(len(s.values))
}
