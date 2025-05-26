package acceptor

import (
	"fmt"
	"log"

	domain_network "orchestra-paxos/internal/domain/network"
	domain_roles "orchestra-paxos/internal/domain/roles"
	network "orchestra-paxos/internal/repository/network"
	sequence_diagram "orchestra-paxos/internal/repository/sequence_diagram"
)

type Acceptor struct {
	GroupID domain_roles.GroupID // group's ID of nodes
	NodeID  domain_roles.NodeID  // Acceptor's ID

	HighestID domain_roles.HighestID // highest ID of proposal in each round
	Learners  []string               // group of learners

	Net network.NetworkActions // network
}

func NewAcceptor(groupID domain_roles.GroupID, nodeID domain_roles.NodeID, net network.NetworkActions) *Acceptor {
	return &Acceptor{
		GroupID: groupID,
		NodeID:  nodeID,
		Net:     net,
	}
}

func (a *Acceptor) handlePrepare(message domain_network.NetworkMessage) {
	msg, ok := message.Data.(domain_network.MessagePrepare)
	if !ok {
		a.log("can not convert message to valid value")

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

		a.log("acceptor %s send NACK to proposer %s", a.Name(), message.Sender)

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

	a.log("acceptor %s send PROMISE to proposer %s", a.Name(), message.Sender)
}

func (a *Acceptor) handleAccept(message domain_network.NetworkMessage) {
	msg, ok := message.Data.(domain_network.MessageAccept)
	if !ok {
		a.log("can not convert message to valid value")

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

		a.log("acceptor %s send NACK to proposer %s", a.Name(), message.Sender)

		return
	}

	for _, learner := range a.Learners {
		a.Net.Send(learner, 0, domain_network.NetworkMessage{
			Stage:  domain_roles.ACCEPTED,
			Sender: a.Name(),
			Data: domain_network.MessageAccepted{
				ProposalID:     a.HighestID,
				Value:          msg.Value,
				ALiveAcceptors: msg.ALiveAcceptors,
			},
		})
		sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d-->> %s:(%d) Accepted: %s", a.NodeID, learner, a.HighestID, msg.Value))
	}

	for _, learner := range a.Learners {
		a.log("acceptor %s send ACCEPTED to learner %s", a.Name(), learner)
	}

	sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d-->> %s:(%d) Accepted: %s", a.NodeID, message.Sender, a.HighestID, msg.Value))
	a.Net.Send(message.Sender, 0, domain_network.NetworkMessage{
		Stage:  domain_roles.ACCEPTED,
		Sender: a.Name(),
		Data: domain_network.MessageAccepted{
			ProposalID:     a.HighestID,
			Value:          msg.Value,
			ALiveAcceptors: msg.ALiveAcceptors,
		},
	})

	a.log("acceptor %s send ACCEPTED to proposer %s", a.Name(), message.Sender)
}

func (a *Acceptor) UpdateLearners(learners []string) {
	a.Learners = learners
}

func (a *Acceptor) Name() string {
	return fmt.Sprintf("Acceptor %d", a.NodeID)
}

func (a *Acceptor) log(format string, v ...any) {
	log.Printf("[%d] %s", a.HighestID, fmt.Sprintf(format, v...))
}
