package acceptor

import (
	"fmt"
	"log"
	"sync"

	domain_network "orchestra-paxos/internal/domain/network"
	domain_roles "orchestra-paxos/internal/domain/roles"
	network "orchestra-paxos/internal/repository/network"
)

type Acceptor struct {
	GroupID domain_roles.GroupID // group's ID of nodes
	NodeID  domain_roles.NodeID  // Acceptor's ID

	HighestID domain_roles.HighestID // highest ID of proposal in each round
	Learners  map[int]struct{}       // group of learners, map for randomize sending (emulate network latency)

	Net   network.NetworkActions // network
	port  int                    // port of acceptor for tcp network
	lock  *sync.Mutex            // mutex
	logMu *sync.Mutex            // mutex for log
}

func NewAcceptor(
	groupID domain_roles.GroupID,
	nodeID domain_roles.NodeID,
	net network.NetworkActions,
	port int) *Acceptor {
	return &Acceptor{
		GroupID:  groupID,
		NodeID:   nodeID,
		Learners: make(map[int]struct{}),
		Net:      net,
		port:     port,
		lock:     &sync.Mutex{},
		logMu:    &sync.Mutex{},
	}
}

func (a *Acceptor) handlePrepare(message *domain_network.NetworkMessage) {
	msg := &domain_network.MessagePrepare{}
	if err := a.Net.Unmarshal(message.Data, msg); err != nil {
		a.log(0, "can not unmarshal message")

		return
	}

	if msg.ProposalID <= a.HighestID {
		a.Net.Send(message.SenderPort, 0, &domain_network.NetworkMessage{
			Stage:      domain_roles.NACK,
			SenderPort: a.Port(),
			Data: domain_network.MessageNack{
				ProposalID: msg.ProposalID,
				AcceptorID: a.HighestID,
			},
		})
		// sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d--x Proposer %s:(%d) Nack", a.Port(), strconv.Itoa(message.SenderPort), a.HighestID))

		a.log(a.HighestID, "acceptor %d send NACK to proposer %d", a.Port(), message.SenderPort)

		return
	}

	a.HighestID = msg.ProposalID

	promise := domain_roles.Proposal{
		ProposalID: int64(msg.ProposalID),
		Value:      "", // mark deafult value, that wasn't decide some consensus value
	}

	a.Net.Send(message.SenderPort, 0, &domain_network.NetworkMessage{
		Stage:      domain_roles.PROMISE,
		SenderPort: a.Port(),
		Data: domain_network.MessagePromise{
			ProposalID: domain_roles.HighestID(promise.ProposalID),
			Value:      promise.Value,
		},
	})
	// sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d-->> Proposer %s:(%d) Promise", a.Port(), strconv.Itoa(message.SenderPort), a.HighestID))

	a.log(msg.ProposalID, "acceptor %d send PROMISE to proposer %d", a.Port(), message.SenderPort)
}

func (a *Acceptor) handleAccept(message *domain_network.NetworkMessage) {
	msg := &domain_network.MessageAccept{}
	if err := a.Net.Unmarshal(message.Data, msg); err != nil {
		a.log(0, "can not unmarshal message")
		return
	}

	if msg.ProposalID < a.HighestID {
		a.Net.Send(message.SenderPort, 0, &domain_network.NetworkMessage{
			Stage:      domain_roles.NACK,
			SenderPort: a.Port(),
			Data: domain_network.MessageNack{
				ProposalID: msg.ProposalID,
				AcceptorID: a.HighestID,
			},
		})
		// sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d--x Proposer %s:(%d) Nack", a.Port(), strconv.Itoa(message.SenderPort), a.HighestID))

		a.log(a.HighestID, "acceptor %d send NACK to proposer %d", a.Port(), message.SenderPort)

		return
	}

	for learner := range a.Learners {
		a.Net.Send(learner, 0, &domain_network.NetworkMessage{
			Stage:      domain_roles.ACCEPTED,
			SenderPort: a.Port(),
			Data: domain_network.MessageAccept{
				ProposalID:     a.HighestID,
				Value:          msg.Value,
				ALiveAcceptors: msg.ALiveAcceptors,
			},
		})
		// sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d-->> Learner %s:(%d) Accepted: %s", a.Port(), strconv.Itoa(learner), a.HighestID, msg.Value))
	}

	for learner := range a.Learners {
		a.log(a.HighestID, "acceptor %d send ACCEPTED to learner %d", a.Port(), learner)
	}

	// sequence_diagram.WriteToFile(fmt.Sprintf("Acceptor %d-->> Proposer %s:(%d) Accepted: %s", a.Port(), strconv.Itoa(message.SenderPort), a.HighestID, msg.Value))
	a.Net.Send(message.SenderPort, 0, &domain_network.NetworkMessage{
		Stage:      domain_roles.ACCEPTED,
		SenderPort: a.Port(),
		Data: domain_network.MessageAccept{
			ProposalID:     a.HighestID,
			Value:          msg.Value,
			ALiveAcceptors: msg.ALiveAcceptors,
		},
	})

	a.log(a.HighestID, "acceptor %d send ACCEPTED to proposer %d", a.Port(), message.SenderPort)
}

func (a *Acceptor) UpdateListOfParticipantsOfTheRequiredRoles(learners []int) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, acceptor := range learners {
		a.Learners[acceptor] = struct{}{}
	}
}

func (a *Acceptor) Port() int {
	return a.port
}

func (a *Acceptor) log(proposerHighestID domain_roles.HighestID, format string, v ...any) {
	a.logMu.Lock()
	log.Printf("[%d] %s", proposerHighestID, fmt.Sprintf(format, v...))
	a.logMu.Unlock()
}
