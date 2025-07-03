package acceptor

import (
	"fmt"
	"log"
	"sync"

	domain_network "orchestra-paxos/internal/domain/network"
	domain_roles "orchestra-paxos/internal/domain/roles"
	network "orchestra-paxos/internal/repository/network"
	sequence_diagram "orchestra-paxos/internal/repository/sequence_diagram"
)

type Acceptor struct {
	GroupID domain_roles.GroupID // group's ID of nodes
	NodeID  domain_roles.NodeID  // Acceptor's ID

	HighestID domain_roles.HighestID // highest ID of proposal in each round
	Learners  map[string]struct{}    // group of learners, map for randomize sending (emulate network latency)

	Net   network.NetworkActions // network
	lock  *sync.Mutex            // mutex
	logMu *sync.Mutex            // mutex for log
}

func NewAcceptor(groupID domain_roles.GroupID, nodeID domain_roles.NodeID, net network.NetworkActions) *Acceptor {
	return &Acceptor{
		GroupID:  groupID,
		NodeID:   nodeID,
		Learners: make(map[string]struct{}),
		Net:      net,
		lock:     &sync.Mutex{},
		logMu:    &sync.Mutex{},
	}
}

func (a *Acceptor) handlePrepare(message domain_network.NetworkMessage) {
	msg, ok := message.Data.(domain_network.MessagePrepare)
	if !ok {
		a.log(0, "can not convert message to valid value")

		return
	}

	if msg.ProposalID <= a.HighestID {
		a.Net.Send(message.Sender, 0, domain_network.NetworkMessage{
			Stage:  domain_roles.NACK,
			Sender: a.Name(),
			Data: domain_network.MessageNack{
				ProposalID: msg.ProposalID,
				AcceptorID: a.HighestID,
			},
		})
		sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d--x %s:(%d) Nack", a.NodeID, message.Sender, a.HighestID))

		a.log(a.HighestID, "acceptor %s send NACK to proposer %s", a.Name(), message.Sender)

		return
	}

	a.HighestID = msg.ProposalID

	promise := domain_roles.Proposal{
		ProposalID: int64(msg.ProposalID),
		Value:      "", // mark deafult value, that wasn't decide some consensus value
	}

	a.Net.Send(message.Sender, 0, domain_network.NetworkMessage{
		Stage:  domain_roles.PROMISE,
		Sender: a.Name(),
		Data: domain_network.MessagePromise{
			ProposalID: domain_roles.HighestID(promise.ProposalID),
			Value:      promise.Value,
		},
	})
	sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d-->> %s:(%d) Promise", a.NodeID, message.Sender, a.HighestID))

	a.log(msg.ProposalID, "acceptor %s send PROMISE to proposer %s", a.Name(), message.Sender)
}

func (a *Acceptor) handleAccept(message domain_network.NetworkMessage) {
	msg, ok := message.Data.(domain_network.MessageAccept)
	if !ok {
		a.log(0, "can not convert message to valid value")

		return
	}

	if msg.ProposalID < a.HighestID {
		a.Net.Send(message.Sender, 0, domain_network.NetworkMessage{
			Stage:  domain_roles.NACK,
			Sender: a.Name(),
			Data: domain_network.MessageNack{
				ProposalID: msg.ProposalID,
				AcceptorID: a.HighestID,
			},
		})
		sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d--x %s:(%d) Nack", a.NodeID, message.Sender, a.HighestID))

		a.log(a.HighestID, "acceptor %s send NACK to proposer %s", a.Name(), message.Sender)

		return
	}

	for learner := range a.Learners {
		a.Net.Send(learner, 0, domain_network.NetworkMessage{
			Stage:  domain_roles.ACCEPTED,
			Sender: a.Name(),
			Data: domain_network.MessageAccept{
				ProposalID:     a.HighestID,
				Value:          msg.Value,
				ALiveAcceptors: msg.ALiveAcceptors,
			},
		})
		sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d-->> %s:(%d) Accepted: %s", a.NodeID, learner, a.HighestID, msg.Value))
	}

	for _, learner := range a.Learners {
		a.log(a.HighestID, "acceptor %s send ACCEPTED to learner %s", a.Name(), learner)
	}

	sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d-->> %s:(%d) Accepted: %s", a.NodeID, message.Sender, a.HighestID, msg.Value))
	a.Net.Send(message.Sender, 0, domain_network.NetworkMessage{
		Stage:  domain_roles.ACCEPTED,
		Sender: a.Name(),
		Data: domain_network.MessageAccept{
			ProposalID:     a.HighestID,
			Value:          msg.Value,
			ALiveAcceptors: msg.ALiveAcceptors,
		},
	})

	a.log(a.HighestID, "acceptor %s send ACCEPTED to proposer %s", a.Name(), message.Sender)
}

func (a *Acceptor) UpdateListOfParticipantsOfTheRequiredRoles(learners []string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, acceptor := range learners {
		a.Learners[acceptor] = struct{}{}
	}
}

func (a *Acceptor) Name() string {
	return fmt.Sprintf("Acceptor %d", a.NodeID)
}

func (a *Acceptor) log(proposerHighestID domain_roles.HighestID, format string, v ...any) {
	a.logMu.Lock()
	log.Printf("[%d] %s", proposerHighestID, fmt.Sprintf(format, v...))
	a.logMu.Unlock()
}
