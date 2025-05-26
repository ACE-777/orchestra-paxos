package learner

import (
	roles "orchestra-paxos/internal/domain/roles"
)

func (l *Learner) Run() {
	l.log("Init Learner in group %d with ID %d", l.GroupID, l.NodeID)
	for {
		msg := l.Net.Receive(l.Name())

		switch msg.Stage {
		case roles.ACCEPTED:
			l.handleAccepted(msg)
		}
	}
}
