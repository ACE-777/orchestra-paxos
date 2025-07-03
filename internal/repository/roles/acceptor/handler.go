package acceptor

import (
	"sync"

	roles "orchestra-paxos/internal/domain/roles"
)

func (a *Acceptor) Run(wg *sync.WaitGroup) {
	a.log(0, "Init Acceptor in group %d with ID %d", a.GroupID, a.NodeID)
	wg.Done()

	go func() {
		for {
			msg := a.Net.Receive(a.Name())

			switch msg.Stage {
			case roles.PREPARE:
				a.handlePrepare(msg)
			case roles.ACCEPT:
				a.handleAccept(msg)
			}
		}
	}()
}
