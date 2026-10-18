package gomega

import (
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	roles "orchestra-paxos/internal/domain/roles"
	"orchestra-paxos/internal/repository/network"
	interface_roles "orchestra-paxos/internal/repository/roles"
	"orchestra-paxos/internal/repository/roles/acceptor"
	"orchestra-paxos/internal/repository/roles/learner"
	"orchestra-paxos/internal/repository/roles/proposer"
	"orchestra-paxos/internal/repository/usecases"
	storage2 "orchestra-paxos/internal/repository/usecases/storage"
)

var _ = Describe("Orchestra Paxos", func() {
	var (
		mockNet     network.NetworkActions
		role        interface_roles.InitRoles
		wgRolesInit sync.WaitGroup
		storage     usecases.Storage

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

	Context("3 proposer, 3 acceptors, 3 learners", func() {
		BeforeEach(func() {
			mockNet = network.NewNetwork(0)
			storage = storage2.NewStorage()

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

			Eventually(func() int {
				return storage.CheckResult()
			}, 60*time.Second, 100*time.Millisecond).Should(Equal(3))

		})

		Describe("UpdateListOfParticipantsOfTheRequiredRoles method", func() {
			It("should update the Learners list", func() {
				// participants := []string{"Acceptor1", "Acceptor2"}
				// role.UpdateListOfParticipantsOfTheRequiredRoles(participants)

				// proposerType, ok := role.(*proposer.Proposer)
				// Expect(ok).To(BeTrue())
				// Expect(proposerType.Acceptors).To(HaveLen(2))
				// Expect(proposerType.Acceptors).To(HaveKey("Acceptor1"))
				// Expect(proposerType.Acceptors).To(HaveKey("Acceptor2"))
			})
		})

	})
})
