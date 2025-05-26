package acceptor

import (
	roles "orchestra-paxos/internal/domain/roles"
)

func (a *Acceptor) Run() {
	a.log("Init Acceptor in group %d with ID %d", a.GroupID, a.NodeID)
	for {
		msg := a.Net.Receive(a.Name())

		switch msg.Stage {
		case roles.PREPARE:
			a.handlePrepare(msg)
		case roles.ACCEPT:
			a.handleAccept(msg)
		}
	}
}
