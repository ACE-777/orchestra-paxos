package network

import (
	"time"

	network "orchestra-paxos/internal/domain/network"
)

// NetworkActions identify actions in network for all network's acceptor
type NetworkActions interface {
	Send(receiverID string, delay time.Duration, msg network.NetworkMessage)
	Receive(receiverID string) network.NetworkMessage
}
