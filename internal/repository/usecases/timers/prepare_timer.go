package timers

import (
	"sync"

	roles "orchestra-paxos/internal/domain/roles"
)

type TimersOfCollectingPrepareFromAcceptors struct {
	timers map[roles.HighestID]bool
	lock   *sync.Mutex
}

func NewTimersOfCollectingPrepareFromAcceptors() *TimersOfCollectingPrepareFromAcceptors {
	return &TimersOfCollectingPrepareFromAcceptors{
		timers: make(map[roles.HighestID]bool),
		lock:   &sync.Mutex{},
	}
}

func (t *TimersOfCollectingPrepareFromAcceptors) CheckExpireTimer(operationID roles.HighestID) bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.timers[operationID]
}

func (t *TimersOfCollectingPrepareFromAcceptors) SetExpireTimer(operationID roles.HighestID) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.timers[operationID] = true
}

func (t *TimersOfCollectingPrepareFromAcceptors) InitExpireTimer(operationID roles.HighestID) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.timers[operationID]; !ok {
		t.timers[operationID] = false
	}
}
