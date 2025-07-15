package learner

import (
	"fmt"
	"log"
	"sync"

	domain_network "orchestra-paxos/internal/domain/network"
	domain_roles "orchestra-paxos/internal/domain/roles"
	network "orchestra-paxos/internal/repository/network"
)

type Learner struct {
	GroupID domain_roles.GroupID // group's ID of nodes
	NodeID  domain_roles.NodeID  // Learner's ID

	Net   network.NetworkActions // network
	port  int                    // port of learner for tcp network
	logMu *sync.Mutex            // mutex for log
}

func NewLearner(
	groupID domain_roles.GroupID,
	nodeID domain_roles.NodeID,
	network network.NetworkActions,
	port int,
) *Learner {
	return &Learner{
		GroupID: groupID,
		NodeID:  nodeID,
		Net:     network,
		port:    port,
		logMu:   &sync.Mutex{},
	}
}

func (l *Learner) handleAccepted(message *domain_network.NetworkMessage) {
	msg := &domain_network.MessageAccept{}
	if err := l.Net.Unmarshal(message.Data, msg); err != nil {
		l.log("learner %d can not unmarshal message", l.Port())
		return
	}

	l.log("learner %d receive value %s from acceptor 127.0.0.1:%d", l.Port(), msg.Value, message.SenderPort)
}

func (l *Learner) UpdateListOfParticipantsOfTheRequiredRoles(roles []int) {}

func (l *Learner) Port() int {
	return l.port
}

func (l *Learner) log(format string, v ...any) {
	l.logMu.Lock()
	log.Printf("Learner [%d] %s", l.Port(), fmt.Sprintf(format, v...))
	l.logMu.Unlock()
}
