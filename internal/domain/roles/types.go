package roles

// GroupID define group ID of nodes
type GroupID uint

// NodeID define node's ID
type NodeID uint

// HighestID define the highest ID of proposal in each round
type HighestID uint64

// Proposal identify proposal in each round of Paxos
type Proposal struct {
	ProposalID int64
	Value      string
}
