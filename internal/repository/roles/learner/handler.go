package learner

import (
	"sync"

	roles "orchestra-paxos/internal/domain/roles"
)

func (l *Learner) Run(wg *sync.WaitGroup) {
	l.log("Init Learner in group %d with ID %d", l.GroupID, l.NodeID)
	wg.Done()

	go func() {
		for {
			msg := l.Net.Receive(l.Name())

			switch msg.Stage {
			case roles.ACCEPTED:
				l.handleAccepted(msg)
			}
		}
	}()
}
