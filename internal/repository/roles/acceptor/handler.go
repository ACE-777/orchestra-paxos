package acceptor

import (
	"encoding/json"
	"fmt"
	"net"
	"orchestra-paxos/internal/domain/network"
	"sync"
	"time"

	roles "orchestra-paxos/internal/domain/roles"
)

func (a *Acceptor) Run(wg *sync.WaitGroup) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", a.Port()))
	if err != nil {
		a.log(0, "Failed to connect to port: %d, error: %v", a.Port(), err)
		return
	}

	a.log(0, "Init Acceptor in group %d with ID %d on 127.0.0.1:%d", a.GroupID, a.NodeID, a.Port())
	wg.Done()

	go func() {
		for {
			connection, err := ln.Accept()
			if err != nil {
				if _, ok := err.(net.Error); ok {
					a.log(0, "Error received while listening 127.0.0.1:%d", a.Port())
				}
			}

			time.Sleep(network.OneWayDelay)

			msg := &network.NetworkMessage{}
			if err = json.NewDecoder(connection).Decode(msg); err != nil {
				a.log(0, "Error decoding %v", err)
			}

			switch msg.Stage {
			case roles.PREPARE:
				a.handlePrepare(msg)
			case roles.ACCEPT:
				a.handleAccept(msg)
			}
		}
	}()
}
