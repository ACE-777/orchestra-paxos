package network

import (
	roles "orchestra-paxos/internal/domain/roles"
)

// NetworkMessage identify message in network between Paxos roles
type NetworkMessage struct {
	Stage  roles.Stage
	Sender string
	Data   interface{}
}

// MessageRequest identify Data in NetworkMessage for prepare state (client -> proposer) or in case of
// refreshing the round
type MessageRequest struct {
	Value   string
	Restart string
}

// MessagePrepare identify Data in NetworkMessage for prepare stage (proposer -> acceptor)
type MessagePrepare struct {
	ProposalID roles.HighestID
}

// MessageNack identify Data in NetworkMessage for all stages in case of impossibility of continuing the algorithm
type MessageNack struct {
	ProposalID roles.HighestID
	AcceptorID roles.HighestID
}

// MessagePromise identify Data in NetworkMessage for prepare stage (acceptor -> proposer)
type MessagePromise struct {
	ProposalID roles.HighestID
	Value      string
}

// MessageAccept identify Data in NetworkMessage for accept stage
// (between acceptor -> learner and then acceptor -> proposer)
type MessageAccept struct {
	ProposalID     roles.HighestID
	Value          string
	ALiveAcceptors []string
}
