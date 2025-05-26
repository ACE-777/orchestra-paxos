package learner

import (
	"fmt"
	"log"

	domain_network "orchestra-paxos/internal/domain/network"
	domain_roles "orchestra-paxos/internal/domain/roles"
	network "orchestra-paxos/internal/repository/network"
)

type Learner struct {
	GroupID domain_roles.GroupID // group's ID of nodes
	NodeID  domain_roles.NodeID  // Learner's ID

	Net network.NetworkActions // network
}

func NewLearner(groupID domain_roles.GroupID, nodeID domain_roles.NodeID, network network.NetworkActions) *Learner {
	return &Learner{
		GroupID: groupID,
		NodeID:  nodeID,
		Net:     network,
	}
}

func (l *Learner) handleAccepted(message domain_network.NetworkMessage) {
	msg, ok := message.Data.(domain_network.MessageAccepted)
	if !ok {
		l.log("can not convert message to valid value")

		return
	}

	l.log("learner %s receive value %s from acceptor %s", l.Name(), msg.Value, message.Sender)
}

func (l *Learner) Name() string {
	return fmt.Sprintf("Learner %d", l.NodeID)
}

func (l *Learner) log(format string, v ...any) {
	log.Printf("Learner [%s] %s", l.Name(), fmt.Sprintf(format, v...))
}
