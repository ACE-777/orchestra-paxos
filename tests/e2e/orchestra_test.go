package e2e

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	network_domain "orchestra-paxos/internal/domain/network"
	roles "orchestra-paxos/internal/domain/roles"
	"orchestra-paxos/internal/repository/network"
	interface_roles "orchestra-paxos/internal/repository/roles"
	"orchestra-paxos/internal/repository/roles/acceptor"
	"orchestra-paxos/internal/repository/roles/learner"
	"orchestra-paxos/internal/repository/roles/proposer"
	"orchestra-paxos/internal/repository/usecases"
	storage2 "orchestra-paxos/internal/repository/usecases/storage"
)

func TestOrchestra(t *testing.T) {
	var (
		mockNet     network.NetworkActions
		wgRolesInit sync.WaitGroup
		storage     usecases.Storage

		learnerNum  = [5]int{9000, 9001, 9002, 9003, 9004} // setup learners
		acceptorNum = [5]int{8995, 8996, 8997, 8998, 8999}
		proposerNum = [5]int{8986, 8987, 8988, 8989, 8990}
		// learnerNum  = [7]int{9000, 9001, 9002, 9003, 9004, 9005, 9006} // setup learners
		// acceptorNum = [7]int{8993, 8994, 8995, 8996, 8997, 8998, 8999}
		// proposerNum = [7]int{8986, 8987, 8988, 8989, 8990, 8991, 8992}
		// learnerNum = [9]int{9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008} // setup learners
		// acceptorNum = [9]int{8991, 8992, 8993, 8994, 8995, 8996, 8997, 8998, 8999}
		// proposerNum = [9]int{8982, 8983, 8984, 8985, 8986, 8987, 8988, 8989, 8990}
		// learnerNum  = [3]int{9006, 9007, 9008} // setup learners
		// acceptorNum = [3]int{9003, 9004, 9005}
		// proposerNum = [3]int{8993, 8994, 8995}
		// acceptorNum = [5]int{9003, 9004, 9005, 9009, 9010} // setup acceptors
		// proposerNum = [10]int{8993, 8994, 8995, 8996, 8997, 8998, 8999, 9000, 9001, 9002} // setup proposers
		//learnerNum   = [1]int{9006} // setup learners
		//acceptorNum  = [1]int{9003} // setup acceptors
		//proposerNum  = [1]int{9000} // setup proposers
		learnersList = []int{}
		acceptorList = []int{}

		learners  = make([]interface_roles.InitRoles, len(learnerNum))
		acceptors = make([]interface_roles.InitRoles, len(acceptorNum))
		proposers = make([]interface_roles.InitRoles, len(proposerNum))
		counter   uint64
	)

	// mockNet = network.NewNetwork(0.31)
	// mockNet = network.NewNetwork(0.06)
	mockNet = network.NewNetwork(0.61)
	storage = storage2.NewStorage()
	// sequence_diagram.CreateNewFile("multi")
	for i, learnerPort := range learnerNum {
		wgRolesInit.Add(1)
		learners[i] = learner.NewLearner(0, roles.NodeID(i), mockNet, learnerPort)
		learnersList = append(learnersList, learners[i].Port())
		go learners[i].Run(&wgRolesInit)
	}

	for i, acceptorPort := range acceptorNum {
		wgRolesInit.Add(1)
		acceptors[i] = acceptor.NewAcceptor(0, roles.NodeID(i), mockNet, acceptorPort)
		acceptors[i].UpdateListOfParticipantsOfTheRequiredRoles(learnersList)
		go acceptors[i].Run(&wgRolesInit)
		acceptorList = append(acceptorList, acceptors[i].Port())
	}

	for i, proposerPort := range proposerNum {
		wgRolesInit.Add(1)
		proposers[i] = proposer.NewProposer(0, roles.NodeID(i), mockNet, proposerPort, storage)
		proposers[i].UpdateListOfParticipantsOfTheRequiredRoles(acceptorList)
		go proposers[i].Run(&wgRolesInit)
	}

	wgRolesInit.Wait()

	timeout := 120 * time.Second
	start := time.Now()
	for {
		if time.Since(start) >= timeout {
			break
		}

		go func() {
			switch proposerFromList := proposers[0].(type) {
			case *proposer.Proposer:
				// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "FirstValue"))
				proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
					Stage:      roles.REQUEST,
					SenderPort: 0,
					Data: network_domain.MessageRequest{
						Value: "FirstValue",
					},
				})

				atomic.AddUint64(&counter, 1)
			default:
				fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
			}
		}()

		// time.Sleep(100*time.Millisecond)
		go func() {
			switch proposerFromList := proposers[1].(type) {
			case *proposer.Proposer:
				// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "SecondValue"))
				proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
					Stage:      roles.REQUEST,
					SenderPort: 0,
					Data: network_domain.MessageRequest{
						Value: "SecondValue",
					},
				})

				atomic.AddUint64(&counter, 1)
			default:
				fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
			}
		}()

		// time.Sleep(100*time.Millisecond)
		go func() {
			switch proposerFromList := proposers[2].(type) {
			case *proposer.Proposer:
				// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "ThirdValue"))
				proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
					Stage:      roles.REQUEST,
					SenderPort: 0,
					Data: network_domain.MessageRequest{
						Value: "ThirdValue",
					},
				})

				atomic.AddUint64(&counter, 1)
			default:
				fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
			}
		}()

		// // time.Sleep(100*time.Millisecond)
		go func() {
			switch proposerFromList := proposers[3].(type) {
			case *proposer.Proposer:
				// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "FirstValue"))
				proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
					Stage:      roles.REQUEST,
					SenderPort: 0,
					Data: network_domain.MessageRequest{
						Value: "FourthValue",
					},
				})

				atomic.AddUint64(&counter, 1)
			default:
				fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
			}
		}()

		// // time.Sleep(100*time.Millisecond)
		go func() {
			switch proposerFromList := proposers[4].(type) {
			case *proposer.Proposer:
				// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "SecondValue"))
				proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
					Stage:      roles.REQUEST,
					SenderPort: 0,
					Data: network_domain.MessageRequest{
						Value: "FivthValue",
					},
				})

				atomic.AddUint64(&counter, 1)
			default:
				fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
			}
		}()

		// // // time.Sleep(100*time.Millisecond)
		// go func() {
		// 	switch proposerFromList := proposers[5].(type) {
		// 	case *proposer.Proposer:
		// 		// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "ThirdValue"))
		// 		proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
		// 			Stage:      roles.REQUEST,
		// 			SenderPort: 0,
		// 			Data: network_domain.MessageRequest{
		// 				Value: "SixthValue",
		// 			},
		// 		})

		// 		atomic.AddUint64(&counter, 1)
		// 	default:
		// 		fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		// 	}
		// }()

		// // // time.Sleep(100*time.Millisecond)
		// go func() {
		// 	switch proposerFromList := proposers[6].(type) {
		// 	case *proposer.Proposer:
		// 		// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "FirstValue"))
		// 		proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
		// 			Stage:      roles.REQUEST,
		// 			SenderPort: 0,
		// 			Data: network_domain.MessageRequest{
		// 				Value: "SeventhValue",
		// 			},
		// 		})

		// 		atomic.AddUint64(&counter, 1)
		// 	default:
		// 		fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		// 	}
		// }()

		// // time.Sleep(100*time.Millisecond)
		// go func() {
		// switch proposerFromList := proposers[7].(type) {
		// case *proposer.Proposer:
		// 	// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "SecondValue"))
		// 	proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
		// 		Stage:      roles.REQUEST,
		// 		SenderPort: 0,
		// 		Data: network_domain.MessageRequest{
		// 			Value: "EighthValue",
		// 		},
		// 	})

		// 	atomic.AddUint64(&counter, 1)
		// default:
		// 	fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		// }
		// }()

		// // // time.Sleep(100*time.Millisecond)
		// go func() {
		// switch proposerFromList := proposers[8].(type) {
		// case *proposer.Proposer:
		// 	// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "ThirdValue"))
		// 	proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
		// 		Stage:      roles.REQUEST,
		// 		SenderPort: 0,
		// 		Data: network_domain.MessageRequest{
		// 			Value: "NinethValue",
		// 		},
		// 	})

		// 	atomic.AddUint64(&counter, 1)
		// default:
		// 	fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		// }
		// }()

		// // time.Sleep(100*time.Millisecond)

		// go func() {
		// switch proposerFromList := proposers[9].(type) {
		// case *proposer.Proposer:
		// 	// sequence_diagram.WriteToFile(fmt.Sprintf("client ->> Proposer %d: Request: %v", proposerFromList.Port(), "ThirdValue"))
		// 	proposerFromList.Net.Send(proposerFromList.Port(), 1, &network_domain.NetworkMessage{
		// 		Stage:      roles.REQUEST,
		// 		SenderPort: 0,
		// 		Data: network_domain.MessageRequest{
		// 			Value: "TenthValue",
		// 		},
		// 	})

		// 	atomic.AddUint64(&counter, 1)
		// default:
		// 	fmt.Printf("Data is of unknown type: %T\n", proposerFromList)
		// }
		// }()

		time.Sleep(1 * time.Second)
	}
	startT := time.Now()
	// time.Sleep(450 * time.Second)

	for {
		currentAccept := storage.CheckResult()
		if currentAccept >= counter {
			break
		}
		fmt.Println("resssssssssssssssssssssssssss===================================================================", currentAccept, " ::: ", counter)
		time.Sleep(5 * time.Second)
	}
	fmt.Println("========time============")
	fmt.Println(time.Since(startT).Seconds())
	fmt.Println("=========================")
	fmt.Println(counter)
	if storage.CheckResult() != counter {
		t.Errorf("Result was incorrect, got: %v, want: %v.", storage.CheckResult(), counter)
	}
}
