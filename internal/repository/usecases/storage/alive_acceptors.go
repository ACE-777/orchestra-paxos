package storage

import (
	"sync"

	roles "orchestra-paxos/internal/domain/roles"
)

type AliveAcceptors struct {
	acceptors map[roles.HighestID][]string
	lock      *sync.Mutex
}

func NewAliveAcceptors() *AliveAcceptors {
	return &AliveAcceptors{
		acceptors: make(map[roles.HighestID][]string),
		lock:      &sync.Mutex{},
	}
}

func (p *AliveAcceptors) AddAcceptor(acceptor string, roundID roles.HighestID) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.acceptors[roundID] = append(p.acceptors[roundID], acceptor)
}

func (p *AliveAcceptors) NumberOfAcceptorsAtRound(roundID roles.HighestID) int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return len(p.acceptors[roundID])
}

func (p *AliveAcceptors) AllAcceptorsAtRound(roundID roles.HighestID) []string {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.acceptors[roundID]
}
