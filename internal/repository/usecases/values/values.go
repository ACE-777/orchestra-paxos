package values

import (
	"sync"

	"orchestra-paxos/internal/domain/roles"
)

type ValuesFromUser struct {
	values map[roles.HighestID]string
	lock   sync.Mutex
}

func NewValuesFromUser() *ValuesFromUser {
	return &ValuesFromUser{
		values: make(map[roles.HighestID]string),
	}
}

func (v *ValuesFromUser) AddValue(value string, operationID roles.HighestID) {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.values[operationID] = value
}

func (v *ValuesFromUser) ValueFromRound(operationID roles.HighestID) string {
	v.lock.Lock()
	defer v.lock.Unlock()

	return v.values[operationID]
}
