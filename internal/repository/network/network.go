package network

import (
	"math/rand"
	"sync"
	"time"

	network "orchestra-paxos/internal/domain/network"
)

const MessageBuffer = 128

type Network struct {
	lock  sync.Mutex
	chans map[string]chan network.NetworkMessage
	loss  uint64
}

func NewNetwork(loss uint64) *Network {
	return &Network{
		chans: make(map[string]chan network.NetworkMessage),
		loss:  loss,
	}
}

func (n *Network) Send(receiverID string, delay int64, message network.NetworkMessage) {
	if delay != 0 {
		time.Sleep(time.Duration(delay) * time.Nanosecond)
	}

	if rand.Uint64() < n.loss {
		//time.Sleep(time.Duration(n.loss) * time.Nanosecond)

		return
	}

	n.lock.Lock()
	if _, ok := n.chans[receiverID]; !ok {
		n.chans[receiverID] = make(chan network.NetworkMessage, MessageBuffer)
	}
	n.lock.Unlock()

	n.chans[receiverID] <- message
}

func (n *Network) Receive(receiverID string) network.NetworkMessage {
	n.lock.Lock()
	if _, ok := n.chans[receiverID]; !ok {
		n.chans[receiverID] = make(chan network.NetworkMessage, MessageBuffer)
	}
	n.lock.Unlock()

	// add random delay for visualization
	//time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return <-n.chans[receiverID]
}
