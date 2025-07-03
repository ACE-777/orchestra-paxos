package proposer

import (
	"sync"

	roles "orchestra-paxos/internal/domain/roles"
)

func (p *Proposer) Run(wg *sync.WaitGroup) {
	p.log(0, "Init Proposer in group %d with ID %d", p.GroupID, p.NodeID)
	wg.Done()

	go func() {
		for {
			msg := p.Net.Receive(p.Name())

			switch msg.Stage {
			case roles.REQUEST:
				p.handleRequest(msg)
			case roles.PROMISE:
				p.handlePromise(msg)
			case roles.ACCEPTED:
				p.handleAccepted(msg)
			case roles.NACK:
				p.handleNack(msg)
			}
		}
	}()
}
