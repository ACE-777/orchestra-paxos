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

type MessageRequest struct {
	Value   string
	Restart string
}

// MessagePrepare identify Data in NetworkMessage for prepare stage (proposer -> acceptor)
type MessagePrepare struct {
	ProposalID roles.HighestID
}

type MessageNack struct {
	ProposalID roles.HighestID
	AcceptorID roles.HighestID
}

// MessageProposal identify Data in NetworkMessage for accept stage (proposer -> acceptor) or in case of already
// accepted value from another proposer with same round ID (acceptor -> proposer)
type MessagePromise struct {
	ProposalID roles.HighestID
	Value      string
}

type MessageAccept struct {
	ProposalID     roles.HighestID
	Value          string
	ALiveAcceptors []string
}

type MessageAccepted struct {
	ProposalID     roles.HighestID
	Value          string
	ALiveAcceptors []string
}

// MessageLearn identify Data in NetworkMessage for broadcast decision of completed round (acceptor -> learner)
type MessageLearn struct {
	Value string
}
