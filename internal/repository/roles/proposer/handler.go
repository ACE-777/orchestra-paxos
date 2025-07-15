package proposer

import (
	"encoding/json"
	"fmt"
	"net"
	"orchestra-paxos/internal/domain/network"
	"sync"
	"time"

	roles "orchestra-paxos/internal/domain/roles"
)

func (p *Proposer) Run(wg *sync.WaitGroup) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p.Port()))
	if err != nil {
		p.log(0, "Failed to connect to port: %d, error: %v", p.Port(), err)
		return
	}

	p.log(0, "Init Proposer in group %d with ID %d on 127.0.0.1:%d", p.GroupID, p.NodeID, p.Port())
	wg.Done()

	go func() {
		for {
			connection, err := ln.Accept()
			if err != nil {
				if _, ok := err.(net.Error); ok {
					p.log(0, "Error received while listening 127.0.0.1:%d", p.Port())
				}
			}

			time.Sleep(network.OneWayDelay)

			msg := &network.NetworkMessage{}
			if err = json.NewDecoder(connection).Decode(msg); err != nil {
				p.log(0, "Error decoding %v", err)
			}

			switch msg.Stage {
			case roles.REQUEST:
				p.handleRequest(msg)
			case roles.PROMISE:
				p.handlePromise(msg)
			case roles.ACCEPTED:
				p.handleAccepted(msg)
			case roles.NACK:
				p.handleNack(msg)
			}
		}
	}()
}
