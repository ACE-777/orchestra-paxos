package proposer

import (
	roles "orchestra-paxos/internal/domain/roles"
)

func (p *Proposer) Run() {
	p.log("Init Proposer in group %d with ID %d", p.GroupID, p.NodeID)
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
}
