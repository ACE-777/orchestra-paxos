package proposer

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	domain_network "orchestra-paxos/internal/domain/network"
	domain_roles "orchestra-paxos/internal/domain/roles"
	network "orchestra-paxos/internal/repository/network"
	sequence_diagram "orchestra-paxos/internal/repository/sequence_diagram"
	interface_uc "orchestra-paxos/internal/repository/usecases"
	operations_log_uc "orchestra-paxos/internal/repository/usecases/operations_log"
	storage_uc "orchestra-paxos/internal/repository/usecases/storage"
	timers_uc "orchestra-paxos/internal/repository/usecases/timers"
	values_uc "orchestra-paxos/internal/repository/usecases/values"
)

// Proposer identify proposer role in Paxos
type Proposer struct {
	GroupID   domain_roles.GroupID   // group's ID of nodes
	NodeID    domain_roles.NodeID    // Proposer's ID
	HighestID domain_roles.HighestID // highest ID of proposal in each round
	Acceptors map[string]struct{}    // all acceptors in system, map for randomize sending (emulate network latency)

	aliveAcceptors                         interface_uc.AcceptorsStorage      // alive acceptors in round
	acceptedAcceptors                      interface_uc.AcceptorsStorage      // acceptors that accept proposal in round
	logOfOperations                        *operations_log_uc.LogOfOperations // marking that some round may can be restored
	timersOfCollectingPrepareFromAcceptors interface_uc.Timers                // mark that timer expired in prepare stage
	timersOfCollectingAcceptFromAcceptors  interface_uc.Timers                // mark that timer expired in prepare stage
	valuesFromUser                         interface_uc.ValuesFromUsers       // values from user for each round

	Net   network.NetworkActions // network
	lock  *sync.Mutex            // mutex
	logMu *sync.Mutex            // mutex for log
}

func NewProposer(
	groupID domain_roles.GroupID,
	nodeID domain_roles.NodeID,
	network network.NetworkActions,
) *Proposer {
	return &Proposer{
		GroupID:                                groupID,
		NodeID:                                 nodeID,
		Acceptors:                              make(map[string]struct{}),
		aliveAcceptors:                         storage_uc.NewAliveAcceptors(),
		acceptedAcceptors:                      storage_uc.NewAcceptedAcceptors(),
		logOfOperations:                        operations_log_uc.NewLogOfOperations(),
		timersOfCollectingPrepareFromAcceptors: timers_uc.NewTimersOfCollectingPrepareFromAcceptors(),
		timersOfCollectingAcceptFromAcceptors:  timers_uc.NewTimersOfCollectingAcceptFromAcceptors(),
		valuesFromUser:                         values_uc.NewValuesFromUser(),
		Net:                                    network,
		lock:                                   &sync.Mutex{},
		logMu:                                  &sync.Mutex{},
	}
}

func (p *Proposer) handleRequest(message domain_network.NetworkMessage) {
	p.lock.Lock()
	p.HighestID++
	roundID := p.HighestID
	p.lock.Unlock()

	msg, ok := message.Data.(domain_network.MessageRequest)
	if !ok {
		p.log(0, "can not convert message to valid value")

		return
	}

	p.valuesFromUser.AddValue(msg.Value, roundID)

	for acceptor := range p.Acceptors {
		p.Net.Send(acceptor, 0, domain_network.NetworkMessage{
			Stage:  domain_roles.PREPARE,
			Sender: p.Name(),
			Data: domain_network.MessagePrepare{
				ProposalID: roundID,
			},
		})
		sequence_diagram.WriteToFile(fmt.Sprintf("Proposer %d-->> %s:(%d) Prepare", p.NodeID, acceptor, roundID))
	}

	for acceptor := range p.Acceptors {
		p.log(roundID, "proposer %s send Prepare to acceptor %s", p.Name(), acceptor)
	}

	p.timersOfCollectingPrepareFromAcceptors.InitExpireTimer(roundID)
	p.sendAccept(roundID, msg)
}

func (p *Proposer) handlePromise(message domain_network.NetworkMessage) {
	msg, ok := message.Data.(domain_network.MessagePromise)
	if !ok {
		p.log(0, "can not convert message to valid value")

		return
	}

	if p.logOfOperations.CheckOperationOnRestartState(msg.ProposalID) {
		return
	}

	if !p.timersOfCollectingPrepareFromAcceptors.CheckExpireTimer(msg.ProposalID) {
		p.aliveAcceptors.AddAcceptor(message.Sender, msg.ProposalID)
	}
}

func (p *Proposer) handleAccepted(message domain_network.NetworkMessage) {
	msg, ok := message.Data.(domain_network.MessageAccept)
	if !ok {
		p.log(0, "can not convert message to valid value")

		return
	}

	if p.logOfOperations.CheckOperationOnRestartState(msg.ProposalID) {
		return
	}

	if !p.timersOfCollectingAcceptFromAcceptors.CheckExpireTimer(msg.ProposalID) {
		p.acceptedAcceptors.AddAcceptor(message.Sender, msg.ProposalID)
	}
}

func (p *Proposer) handleNack(message domain_network.NetworkMessage) {
	msg, ok := message.Data.(domain_network.MessageNack)
	if !ok {
		p.log(0, "can not convert message to valid value")

		return
	}

	if p.logOfOperations.CheckOperationOnRestartState(msg.ProposalID) {
		return
	}

	p.lock.Lock()
	if msg.AcceptorID > p.HighestID {
		p.HighestID = msg.AcceptorID
	}
	p.lock.Unlock()

	p.logOfOperations.SetRestartStateOperation(msg.ProposalID)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randomTimeout := time.Duration(r.Intn(500)+1000) * time.Microsecond
	fmt.Printf("Sleeping %d microseconds before restarting round\n", randomTimeout)
	time.Sleep(randomTimeout)
	p.handleRequest(domain_network.NetworkMessage{
		Stage:  domain_roles.REQUEST,
		Sender: p.Name(),
		Data: domain_network.MessageRequest{
			Value:   p.valuesFromUser.ValueFromRound(msg.ProposalID),
			Restart: "restart",
		},
	})

	p.log(msg.ProposalID, "restart round")
}

func (p *Proposer) sendAccept(roundID domain_roles.HighestID, msg domain_network.MessageRequest) {
	go func(roundID domain_roles.HighestID) {
		ticker := time.NewTicker(20 * time.Microsecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if p.logOfOperations.CheckOperationOnRestartState(roundID) {
					sequence_diagram.WriteToFile(
						fmt.Sprintf("Note left of Proposer %d: round %d restarted", p.NodeID, roundID),
					)

					return
				}

				p.timersOfCollectingPrepareFromAcceptors.SetExpireTimer(roundID)
				//sequence_diagram.WriteToFile(fmt.Sprintf("participant client\n==%s==", "timer Prepare expired"))
				sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: %s", p.NodeID, "timer Prepare expired"))
				p.log(roundID, "timer Prepare expired")

				msgAccept := domain_network.MessageAccept{
					ProposalID:     roundID,
					Value:          p.valuesFromUser.ValueFromRound(roundID),
					ALiveAcceptors: p.aliveAcceptors.AllAcceptorsAtRound(roundID),
				}

				for _, acceptor := range p.aliveAcceptors.AllAcceptorsAtRound(roundID) {
					p.Net.Send(acceptor, 0, domain_network.NetworkMessage{
						Stage:  domain_roles.ACCEPT,
						Sender: p.Name(),
						Data:   msgAccept,
					})
					sequence_diagram.WriteToFile(fmt.Sprintf("Proposer %d-->> %s:(%d) Accept: %s", p.NodeID, acceptor, roundID, p.valuesFromUser.ValueFromRound(roundID)))
				}

				for _, acceptor := range p.aliveAcceptors.AllAcceptorsAtRound(roundID) {
					p.log(roundID, "proposer %s send ACCEPT to acceptor %s", p.Name(), acceptor)
				}

				p.waitAccept(roundID, msg)

				return
			default:
				if p.logOfOperations.CheckOperationOnRestartState(roundID) {
					sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: round %d restarted", p.NodeID, roundID))

					return
				}

				continue
			}
		}
	}(roundID)
}

func (p *Proposer) waitAccept(roundID domain_roles.HighestID, msg domain_network.MessageRequest) {
	go func(roundID domain_roles.HighestID) {
		ticker1 := time.NewTicker(20 * time.Microsecond)
		defer ticker1.Stop()
		for {
			select {
			case <-ticker1.C:
				if p.logOfOperations.CheckOperationOnRestartState(roundID) {
					sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: round %d restarted", p.NodeID, roundID))

					return
				}

				p.timersOfCollectingAcceptFromAcceptors.SetExpireTimer(roundID)
				p.log(roundID, "timer Accept expired")
				//sequence_diagram.WriteToFile(fmt.Sprintf("participant client\n==%s==", "timer Accept expired"))
				sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: %s", p.NodeID, "timer Accept expired"))
				if p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID) > p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2 {
					sequence_diagram.WriteToFile(fmt.Sprintf("Proposer %d-->> client: %s was accepted as the value!", p.NodeID, p.valuesFromUser.ValueFromRound(roundID)))
					p.log(roundID, "!!!!!  value %s was accepted, proposer %s ::: %d > %d", msg.Value, p.Name(), p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID), p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2)

					return
				}

				p.logOfOperations.SetRestartStateOperation(roundID)
				source := rand.NewSource(time.Now().UnixNano())
				r := rand.New(source)
				randomTimeout := time.Duration(r.Intn(500)+1000) * time.Microsecond
				fmt.Printf("Sleeping %d microseconds before restarting round \n", randomTimeout)
				time.Sleep(randomTimeout)

				p.handleRequest(domain_network.NetworkMessage{
					Stage:  domain_roles.REQUEST,
					Sender: p.Name(),
					Data: domain_network.MessageRequest{
						Value:   p.valuesFromUser.ValueFromRound(roundID),
						Restart: "restart",
					},
				})

				p.log(roundID, "restart round")

				return
			default:
				if p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID) > p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2 {
					sequence_diagram.WriteToFile(fmt.Sprintf("Proposer %d-->> client: %s was accepted as the value!", p.NodeID, p.valuesFromUser.ValueFromRound(roundID)))
					p.log(roundID, "!!!!!  value %s was accepted, proposer %s ::: %d > %d", msg.Value, p.Name(), p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID), p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2)

					return
				}

				continue
			}
		}
	}(roundID)
}

func (p *Proposer) UpdateListOfParticipantsOfTheRequiredRoles(acceptors []string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, acceptor := range acceptors {
		p.Acceptors[acceptor] = struct{}{}
	}
}

func (p *Proposer) Name() string {
	return fmt.Sprintf("Proposer %d", p.NodeID)
}

func (p *Proposer) log(proposerHighestID domain_roles.HighestID, format string, v ...any) {
	p.logMu.Lock()
	log.Printf("[%d] %s", proposerHighestID, fmt.Sprintf(format, v...))
	p.logMu.Unlock()
}
