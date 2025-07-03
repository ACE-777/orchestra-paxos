package network

import (
	"math/rand"
	"sync"
	"time"

	network "orchestra-paxos/internal/domain/network"
)

const MessageBuffer = 1024

type Network struct {
	chans map[string]chan network.NetworkMessage
	loss  uint64
	lock  *sync.RWMutex
}

func NewNetwork(loss uint64) *Network {
	return &Network{
		chans: make(map[string]chan network.NetworkMessage),
		loss:  loss,
		lock:  &sync.RWMutex{},
	}
}

func (n *Network) getChannel(receiverID string) chan network.NetworkMessage {
	n.lock.RLock()
	channel, exists := n.chans[receiverID]
	n.lock.RUnlock()

	if exists {
		return channel
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	if channel, exists = n.chans[receiverID]; exists {
		return channel
	}

	newChannel := make(chan network.NetworkMessage, MessageBuffer)
	n.chans[receiverID] = newChannel

	return newChannel
}

func (n *Network) Send(receiverID string, delay time.Duration, message network.NetworkMessage) {
	if delay != 0 {
		time.Sleep(delay)
	}

	if rand.Uint64() <= n.loss {
		//time.Sleep(time.Duration(n.loss) * time.Nanosecond)

		return
	}
	channel := n.getChannel(receiverID)
	channel <- message
}

func (n *Network) Receive(receiverID string) network.NetworkMessage {
	channel := n.getChannel(receiverID)
	msg := <-channel
	//close(channel)

	return msg
}
