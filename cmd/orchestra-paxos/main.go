package main

import (
	"fmt"
	"time"

	network_domain "orchestra-paxos/internal/domain/network"
	roles "orchestra-paxos/internal/domain/roles"
	network_repository "orchestra-paxos/internal/repository/network"
	acceptor "orchestra-paxos/internal/repository/roles/acceptor"
	learner "orchestra-paxos/internal/repository/roles/learner"
	proposer "orchestra-paxos/internal/repository/roles/proposer"
	sequence_diagram "orchestra-paxos/internal/repository/sequence_diagram"
)

func main() {
	sequence_diagram.CreateNewFile("multi")

	// setup internal network
	network := network_repository.NewNetwork(0)
	// setup learners
	learnerNum := 3
	learners := make([]*learner.Learner, learnerNum)
	learnersList := []string{}
	for i := range learnerNum {
		learners[i] = learner.NewLearner(0, roles.NodeID(i), network)
		learnersList = append(learnersList, learners[i].Name())
		go learners[i].Run()
	}
	// setup acceptors
	acceptorNum := 3
	acceptors := make([]*acceptor.Acceptor, acceptorNum)
	for i := range acceptorNum {
		acceptors[i] = acceptor.NewAcceptor(0, roles.NodeID(i), network)
		acceptors[i].UpdateLearners(learnersList)
		go acceptors[i].Run()
	}
	acceptorList := []string{}
	for _, a := range acceptors {
		acceptorList = append(acceptorList, a.Name())
	}
	// setup proposers
	proposerNum := 3
	proposers := make([]*proposer.Proposer, proposerNum)
	for i := range proposerNum {
		proposers[i] = proposer.NewProposer(0, roles.NodeID(i), network)
		proposers[i].UpdateAcceptor(acceptorList)
		go proposers[i].Run()
	}

	time.Sleep(2 * time.Second)

	go func() {
		sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer 0: Request: %v", "FirstValue"))

		proposers[0].Net.Send(proposers[0].Name(), 0, network_domain.NetworkMessage{
			Stage:  roles.REQUEST,
			Sender: "client",
			Data: network_domain.MessageRequest{
				Value: "ThirstValue",
			},
		})
	}()

	go func() {
		sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer 1: Request: %v", "SecondValue"))
		proposers[1].Net.Send(proposers[1].Name(), 0, network_domain.NetworkMessage{
			Stage:  roles.REQUEST,
			Sender: "client",
			Data: network_domain.MessageRequest{
				Value: "SecondValue",
			},
		})
	}()

	go func() {
		sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer 2: Request: %v", "ThirdValue"))
		proposers[2].Net.Send(proposers[2].Name(), 0, network_domain.NetworkMessage{
			Stage:  roles.REQUEST,
			Sender: proposers[2].Name(),
			Data: network_domain.MessageRequest{
				Value: "ThirdValue",
			},
		})
	}()

	time.Sleep(100 * time.Second)
}
