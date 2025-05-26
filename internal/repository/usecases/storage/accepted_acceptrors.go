package storage

import (
	"sync"

	"orchestra-paxos/internal/domain/roles"
)

type AcceptedAcceptors struct {
	acceptors map[roles.HighestID][]string
	lock      sync.Mutex
}

func NewAcceptedAcceptors() *AcceptedAcceptors {
	return &AcceptedAcceptors{
		acceptors: make(map[roles.HighestID][]string),
	}
}

func (p *AcceptedAcceptors) AddAcceptor(acceptor string, roundID roles.HighestID) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, ok := p.acceptors[roundID]; !ok {
		p.acceptors[roundID] = []string{acceptor}

		return
	}

	p.acceptors[roundID] = append(p.acceptors[roundID], acceptor)
}

func (p *AcceptedAcceptors) NumberOfAcceptorsAtRound(roundID roles.HighestID) int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return len(p.acceptors[roundID])
}

func (p *AcceptedAcceptors) AllAcceptorsAtRound(roundID roles.HighestID) []string {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.acceptors[roundID]
}
