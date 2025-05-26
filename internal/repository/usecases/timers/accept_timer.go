package timers

import (
	"sync"

	"orchestra-paxos/internal/domain/roles"
)

type TimersOfCollectingAcceptFromAcceptors struct {
	timers map[roles.HighestID]bool
	lock   sync.Mutex
}

func NewTimersOfCollectingAcceptFromAcceptors() *TimersOfCollectingAcceptFromAcceptors {
	return &TimersOfCollectingAcceptFromAcceptors{
		timers: make(map[roles.HighestID]bool),
	}
}

func (t *TimersOfCollectingAcceptFromAcceptors) CheckExpireTimer(operationID roles.HighestID) bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.timers[operationID]
}

func (t *TimersOfCollectingAcceptFromAcceptors) SetExpireTimer(operationID roles.HighestID) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.timers[operationID] = true
}

func (t *TimersOfCollectingAcceptFromAcceptors) InitExpireTimer(operationID roles.HighestID) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.timers[operationID]; !ok {
		t.timers[operationID] = false
	}
}
