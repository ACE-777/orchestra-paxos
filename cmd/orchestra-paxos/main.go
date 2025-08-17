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

		learnerNum  = [3]int{9006, 9007, 9008} // setup learners
		acceptorNum = [3]int{9003, 9004, 9005} // setup acceptors
		proposerNum = [3]int{9000, 9001, 9002} // setup proposers
		//learnerNum   = [1]int{9006} // setup learners
		//acceptorNum  = [1]int{9003} // setup acceptors
		//proposerNum  = [1]int{9000} // setup proposers
		learnersList = []int{}
		acceptorList = []int{}

		learners  = make([]interface_roles.InitRoles, len(learnerNum))
		acceptors = make([]interface_roles.InitRoles, len(acceptorNum))
		proposers = make([]interface_roles.InitRoles, len(proposerNum))
	)

	sequence_diagram.CreateNewFile("multi")

	for i, learnerPort := range learnerNum {
		wgRolesInit.Add(1)
		learners[i] = learner.NewLearner(0, roles.NodeID(i), network, learnerPort)
		learnersList = append(learnersList, learners[i].Port())
		go learners[i].Run(&wgRolesInit)
	}

	for i, acceptorPort := range acceptorNum {
		wgRolesInit.Add(1)
		acceptors[i] = acceptor.NewAcceptor(0, roles.NodeID(i), network, acceptorPort)
		acceptors[i].UpdateListOfParticipantsOfTheRequiredRoles(learnersList)
		go acceptors[i].Run(&wgRolesInit)
		acceptorList = append(acceptorList, acceptors[i].Port())
	}

	for i, proposerPort := range proposerNum {
		wgRolesInit.Add(1)
		proposers[i] = proposer.NewProposer(0, roles.NodeID(i), network, proposerPort)
		proposers[i].UpdateListOfParticipantsOfTheRequiredRoles(acceptorList)
		go proposers[i].Run(&wgRolesInit)
	}

	wgRolesInit.Wait()

	go func() {
		switch proposerFromList := proposers[0].(type) {
		case *proposer.Proposer:
			sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "FirstValue"))
			proposerFromList.Net.Send(proposerFromList.Port(), 0, &network_domain.NetworkMessage{
				Stage:      roles.REQUEST,
				SenderPort: 0,
				Data: network_domain.MessageRequest{
					Value: "FirstValue",
				},
			})
		default:
			fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		}
	}()

	//go func() {
	//	switch proposerFromList := proposers[0].(type) {
	//	case *proposer.Proposer:
	//		sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "SecondValue"))
	//		proposerFromList.Net.Send(proposerFromList.Port(), 0, &network_domain.NetworkMessage{
	//			Stage:      roles.REQUEST,
	//			SenderPort: 0,
	//			Data: network_domain.MessageRequest{
	//				Value: "SecondValue",
	//			},
	//		})
	//	default:
	//		fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
	//	}
	//}()

	//go func() {
	//	sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer 0: Request: %v", "FValue"))
	//	switch proposerFromList := proposers[0].(type) {
	//	case *proposer.Proposer:
	//		proposerFromList.Net.Send(proposerFromList.Name(), 0, network_domain.NetworkMessage{
	//			Stage:  roles.REQUEST,
	//			Sender: "client",
	//			Data: network_domain.MessageRequest{
	//				Value: "FValue",
	//			},
	//		})
	//	default:
	//		fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
	//	}
	//}()

	go func() {
		switch proposerFromList := proposers[1].(type) {
		case *proposer.Proposer:
			sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "SecondValue"))
			proposerFromList.Net.Send(proposerFromList.Port(), 0, &network_domain.NetworkMessage{
				Stage:      roles.REQUEST,
				SenderPort: 0,
				Data: network_domain.MessageRequest{
					Value: "SecondValue",
				},
			})
		default:
			fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		}
	}()

	go func() {
		switch proposerFromList := proposers[2].(type) {
		case *proposer.Proposer:
			sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "ThirdValue"))
			proposerFromList.Net.Send(proposerFromList.Port(), 0, &network_domain.NetworkMessage{
				Stage:      roles.REQUEST,
				SenderPort: 0,
				Data: network_domain.MessageRequest{
					Value: "ThirdValue",
				},
			})
		default:
			fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		}
	}()

	time.Sleep(60 * time.Second)
}
