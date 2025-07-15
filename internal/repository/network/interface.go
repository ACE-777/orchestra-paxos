package network

import (
	"time"

	network "orchestra-paxos/internal/domain/network"
)

// NetworkActions identify actions in network for all network's acceptor
type NetworkActions interface {
	Unmarshal(in interface{}, out interface{}) error
	Send(receiverID int, delay time.Duration, msg *network.NetworkMessage)
}
