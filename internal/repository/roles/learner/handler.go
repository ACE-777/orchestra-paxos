package learner

import (
	"encoding/json"
	"fmt"
	"net"
	"orchestra-paxos/internal/domain/network"
	"sync"
	"time"

	roles "orchestra-paxos/internal/domain/roles"
)

func (l *Learner) Run(wg *sync.WaitGroup) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", l.Port()))
	if err != nil {
		l.log("Failed to connect to port: %d, error: %v", l.Port(), err)
		return
	}

	l.log("Init Learner in group %d with ID %d on 127.0.0.1:%d", l.GroupID, l.NodeID, l.Port())
	wg.Done()

	go func() {
		for {
			connection, err := ln.Accept()
			if err != nil {
				if _, ok := err.(net.Error); ok {
					l.log("Error received while listening 127.0.0.1:%d", l.Port())
				}
			}

			time.Sleep(network.OneWayDelay)

			msg := &network.NetworkMessage{}
			if err = json.NewDecoder(connection).Decode(msg); err != nil {
				l.log("Error decoding %v", err)
			}

			switch msg.Stage {
			case roles.ACCEPTED:
				l.handleAccepted(msg)
			}
		}
	}()
}
