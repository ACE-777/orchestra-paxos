package operations_log

import (
	"sync"

	roles "orchestra-paxos/internal/domain/roles"
)

type LogOfOperations struct {
	logOfOperations map[roles.HighestID]struct{}
	lock            *sync.Mutex
}

func NewLogOfOperations() *LogOfOperations {
	return &LogOfOperations{
		logOfOperations: make(map[roles.HighestID]struct{}),
		lock:            &sync.Mutex{},
	}
}

func (l *LogOfOperations) CheckOperationOnRestartState(operationID roles.HighestID) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	_, ok := l.logOfOperations[operationID]

	return ok
}

func (l *LogOfOperations) SetRestartStateOperation(operationID roles.HighestID) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.logOfOperations[operationID] = struct{}{}
}
