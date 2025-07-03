package main

import (
	"fmt"
	"sync"
	"time"

	network_domain "orchestra-paxos/internal/domain/network"
	roles "orchestra-paxos/internal/domain/roles"
	network_repository "orchestra-paxos/internal/repository/network"
	interface_roles "orchestra-paxos/internal/repository/roles"
	acceptor "orchestra-paxos/internal/repository/roles/acceptor"
	learner "orchestra-paxos/internal/repository/roles/learner"
	proposer "orchestra-paxos/internal/repository/roles/proposer"
	sequence_diagram "orchestra-paxos/internal/repository/sequence_diagram"
)

func main() {
	var (
		network = network_repository.NewNetwork(0) // setup internal network

		wgRolesInit sync.WaitGroup
		//wg          sync.WaitGroup

		learnerNum   = 3 // setup learners
		acceptorNum  = 3 // setup acceptors
		proposerNum  = 3 // setup proposers
		learnersList = []string{}
		acceptorList = []string{}

		learners  = make([]interface_roles.InitRoles, learnerNum)
		acceptors = make([]interface_roles.InitRoles, acceptorNum)
		proposers = make([]interface_roles.InitRoles, proposerNum)
	)

	sequence_diagram.CreateNewFile("multi")

	for i := range learnerNum {
		wgRolesInit.Add(1)
		learners[i] = learner.NewLearner(0, roles.NodeID(i), network)
		learnersList = append(learnersList, learners[i].Name())
		go learners[i].Run(&wgRolesInit)
	}

	for i := range acceptorNum {
		wgRolesInit.Add(1)
		acceptors[i] = acceptor.NewAcceptor(0, roles.NodeID(i), network)
		acceptors[i].UpdateListOfParticipantsOfTheRequiredRoles(learnersList)
		go acceptors[i].Run(&wgRolesInit)
		acceptorList = append(acceptorList, acceptors[i].Name())
	}

	for i := range proposerNum {
		wgRolesInit.Add(1)
		proposers[i] = proposer.NewProposer(0, roles.NodeID(i), network)
		proposers[i].UpdateListOfParticipantsOfTheRequiredRoles(acceptorList)
		go proposers[i].Run(&wgRolesInit)
	}

	wgRolesInit.Wait()

	go func() {
		sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer 0: Request: %v", "FirstValue"))
		switch proposerFromList := proposers[0].(type) {
		case *proposer.Proposer:
			proposerFromList.Net.Send(proposerFromList.Name(), 0, network_domain.NetworkMessage{
				Stage:  roles.REQUEST,
				Sender: "client",
				Data: network_domain.MessageRequest{
					Value: "FirstValue",
				},
			})
		default:
			fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		}
	}()

	go func() {
		sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer 1: Request: %v", "SecondValue"))
		switch proposerFromList := proposers[1].(type) {
		case *proposer.Proposer:
			proposerFromList.Net.Send(proposerFromList.Name(), 0, network_domain.NetworkMessage{
				Stage:  roles.REQUEST,
				Sender: "client",
				Data: network_domain.MessageRequest{
					Value: "SecondValue",
				},
			})
		default:
			fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		}
	}()

	go func() {
		sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer 2: Request: %v", "ThirdValue"))
		switch proposerFromList := proposers[2].(type) {
		case *proposer.Proposer:
			proposerFromList.Net.Send(proposerFromList.Name(), 0, network_domain.NetworkMessage{
				Stage:  roles.REQUEST,
				Sender: "client",
				Data: network_domain.MessageRequest{
					Value: "ThirdValue",
				},
			})
		default:
			fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		}
	}()

	time.Sleep(20 * time.Second)
}
