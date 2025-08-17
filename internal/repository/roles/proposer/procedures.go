package proposer

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
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
	Acceptors map[int]struct{}       // all acceptors in system, map for randomize sending (emulate network latency)

	aliveAcceptors                         interface_uc.AcceptorsStorage      // alive acceptors in round
	acceptedAcceptors                      interface_uc.AcceptorsStorage      // acceptors that accept proposal in round
	logOfOperations                        *operations_log_uc.LogOfOperations // marking that some round may can be restored
	timersOfCollectingPrepareFromAcceptors interface_uc.Timers                // mark that timer expired in prepare stage
	timersOfCollectingAcceptFromAcceptors  interface_uc.Timers                // mark that timer expired in prepare stage
	valuesFromUser                         interface_uc.ValuesFromUsers       // values from user for each round

	Net   network.NetworkActions // network
	port  int                    // port of producer for tcp network
	lock  *sync.Mutex            // mutex
	logMu *sync.Mutex            // mutex for log
}

func NewProposer(
	groupID domain_roles.GroupID,
	nodeID domain_roles.NodeID,
	network network.NetworkActions,
	port int,
) *Proposer {
	return &Proposer{
		GroupID:                                groupID,
		NodeID:                                 nodeID,
		Acceptors:                              make(map[int]struct{}),
		aliveAcceptors:                         storage_uc.NewAliveAcceptors(),
		acceptedAcceptors:                      storage_uc.NewAcceptedAcceptors(),
		logOfOperations:                        operations_log_uc.NewLogOfOperations(),
		timersOfCollectingPrepareFromAcceptors: timers_uc.NewTimersOfCollectingPrepareFromAcceptors(),
		timersOfCollectingAcceptFromAcceptors:  timers_uc.NewTimersOfCollectingAcceptFromAcceptors(),
		valuesFromUser:                         values_uc.NewValuesFromUser(),
		Net:                                    network,
		port:                                   port,
		lock:                                   &sync.Mutex{},
		logMu:                                  &sync.Mutex{},
	}
}

func (p *Proposer) handleRequest(message *domain_network.NetworkMessage) {
	msg := &domain_network.MessageRequest{}
	if err := p.Net.Unmarshal(message.Data, msg); err != nil {
		p.log(0, "can not unmarshal message")
		return
	}

	p.lock.Lock()
	p.HighestID++ // maybe increase on delta or not needed
	roundID := p.HighestID
	p.lock.Unlock()

	p.valuesFromUser.AddValue(msg.Value, roundID)

	for acceptor := range p.Acceptors {
		p.Net.Send(acceptor, 0, &domain_network.NetworkMessage{
			Stage:      domain_roles.PREPARE,
			SenderPort: p.Port(),
			Data: domain_network.MessagePrepare{
				ProposalID: roundID,
			},
		})
		sequence_diagram.WriteToFile(fmt.Sprintf("Proposer %d-->> Acceptor %s:(%d) Prepare", p.Port(), strconv.Itoa(acceptor), roundID))
	}

	for acceptor := range p.Acceptors {
		p.log(roundID, "proposer %d send Prepare to acceptor %d", p.Port(), acceptor)
	}

	p.timersOfCollectingPrepareFromAcceptors.InitExpireTimer(roundID)
	p.sendAccept(roundID, msg)
}

func (p *Proposer) handlePromise(message *domain_network.NetworkMessage) {
	msg := &domain_network.MessagePromise{}
	if err := p.Net.Unmarshal(message.Data, msg); err != nil {
		p.log(0, "can not unmarshal message")
		return
	}

	if p.logOfOperations.CheckOperationOnRestartState(msg.ProposalID) {
		return
	}

	if !p.timersOfCollectingPrepareFromAcceptors.CheckExpireTimer(msg.ProposalID) {
		p.aliveAcceptors.AddAcceptor(message.SenderPort, msg.ProposalID)
	}
}

func (p *Proposer) handleAccepted(message *domain_network.NetworkMessage) {
	msg := &domain_network.MessageAccept{}
	if err := p.Net.Unmarshal(message.Data, msg); err != nil {
		p.log(0, "can not unmarshal message")
		return
	}

	if p.logOfOperations.CheckOperationOnRestartState(msg.ProposalID) {
		return
	}

	if !p.timersOfCollectingAcceptFromAcceptors.CheckExpireTimer(msg.ProposalID) {
		p.acceptedAcceptors.AddAcceptor(message.SenderPort, msg.ProposalID)
	}
}

func (p *Proposer) handleNack(message *domain_network.NetworkMessage) {
	msg := &domain_network.MessageNack{}
	if err := p.Net.Unmarshal(message.Data, msg); err != nil {
		p.log(0, "can not unmarshal message")
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

	if p.logOfOperations.SetRestartStateOperation(msg.ProposalID) {
		p.log(msg.ProposalID, "restart round")
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		randomTimeout := time.Duration(r.Intn(1000)+200) * time.Millisecond
		p.log(msg.ProposalID, "Sleeping %d milliseconds before restarting round", randomTimeout.Milliseconds())
		time.Sleep(randomTimeout)

		p.handleRequest(&domain_network.NetworkMessage{
			Stage:      domain_roles.REQUEST,
			SenderPort: p.Port(),
			Data: domain_network.MessageRequest{
				Value:   p.valuesFromUser.ValueFromRound(msg.ProposalID),
				Restart: "restart",
			},
		})
	}
}

func (p *Proposer) sendAccept(roundID domain_roles.HighestID, msg *domain_network.MessageRequest) {
	go func(roundID domain_roles.HighestID) {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if p.logOfOperations.CheckOperationOnRestartState(roundID) {
					sequence_diagram.WriteToFile(
						fmt.Sprintf("Note left of Proposer %d: round %d restarted", p.Port(), roundID),
					)

					p.log(roundID, "restart round")

					return
				}

				p.timersOfCollectingPrepareFromAcceptors.SetExpireTimer(roundID)
				//sequence_diagram.WriteToFile(fmt.Sprintf("participant client\n==%s==", "timer Prepare expired"))
				sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: (%v)%s", p.Port(), roundID, "timer Prepare expired"))
				p.log(roundID, "timer Prepare expired")

				msgAccept := domain_network.MessageAccept{
					ProposalID:     roundID,
					Value:          p.valuesFromUser.ValueFromRound(roundID),
					ALiveAcceptors: p.aliveAcceptors.AllAcceptorsAtRound(roundID),
				}

				for _, acceptor := range p.aliveAcceptors.AllAcceptorsAtRound(roundID) {
					p.Net.Send(acceptor, 0, &domain_network.NetworkMessage{
						Stage:      domain_roles.ACCEPT,
						SenderPort: p.Port(),
						Data:       msgAccept,
					})
					sequence_diagram.WriteToFile(fmt.Sprintf("Proposer %d-->> Acceptor %s:(%d) Accept: %s", p.Port(), strconv.Itoa(acceptor), roundID, p.valuesFromUser.ValueFromRound(roundID)))
				}

				for _, acceptor := range p.aliveAcceptors.AllAcceptorsAtRound(roundID) {
					p.log(roundID, "proposer %d send ACCEPT to acceptor %d", p.Port(), acceptor)
				}

				p.waitAccept(roundID, msg)

				return
			default:
				if p.logOfOperations.CheckOperationOnRestartState(roundID) {
					sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: round %d restarted", p.Port(), roundID))

					return
				}

				continue
			}
		}
	}(roundID)
}

func (p *Proposer) waitAccept(roundID domain_roles.HighestID, msg *domain_network.MessageRequest) {
	go func(roundID domain_roles.HighestID) {
		ticker1 := time.NewTicker(100 * time.Millisecond)
		defer ticker1.Stop()
		for {
			select {
			case <-ticker1.C:
				if p.logOfOperations.CheckOperationOnRestartState(roundID) {
					sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: round %d restarted", p.Port(), roundID))

					return
				}

				p.timersOfCollectingAcceptFromAcceptors.SetExpireTimer(roundID)
				p.log(roundID, "timer Accept expired")
				//sequence_diagram.WriteToFile(fmt.Sprintf("participant client\n==%s==", "timer Accept expired"))
				sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: (%v)%s", p.Port(), roundID, "timer Accept expired"))
				if p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID) > p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2 {
					sequence_diagram.WriteToFile(fmt.Sprintf("Proposer %d-->> client: %s was accepted as the value!", p.Port(), p.valuesFromUser.ValueFromRound(roundID)))
					p.log(roundID, "!!!!!  value %s was accepted, proposer %d ::: %d > %d", msg.Value, p.Port(), p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID), p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2)

					return
				} else {
					p.log(roundID, "number: %v, alive: %v", p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID), p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2)
				}

				if p.logOfOperations.SetRestartStateOperation(roundID) {
					p.log(roundID, "restart round")
					source := rand.NewSource(time.Now().UnixNano())
					r := rand.New(source)
					randomTimeout := time.Duration(r.Intn(1000)+200) * time.Millisecond
					p.log(roundID, "Sleeping %d milliseconds before restarting round", randomTimeout.Milliseconds())
					time.Sleep(randomTimeout)

					p.handleRequest(&domain_network.NetworkMessage{
						Stage:      domain_roles.REQUEST,
						SenderPort: p.Port(),
						Data: domain_network.MessageRequest{
							Value:   p.valuesFromUser.ValueFromRound(roundID),
							Restart: "restart",
						},
					})
				}

				return
			default:
				if p.logOfOperations.CheckOperationOnRestartState(roundID) {
					sequence_diagram.WriteToFile(fmt.Sprintf("Note left of Proposer %d: round %d restarted", p.Port(), roundID))

					p.log(roundID, "restart round")

					return
				}

				if p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID) > p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2 {
					sequence_diagram.WriteToFile(fmt.Sprintf("Proposer %d-->> client: %s was accepted as the value!", p.Port(), p.valuesFromUser.ValueFromRound(roundID)))
					p.log(roundID, "!!!!!  value %s was accepted, proposer %d ::: %d > %d", msg.Value, p.Port(), p.acceptedAcceptors.NumberOfAcceptorsAtRound(roundID), p.aliveAcceptors.NumberOfAcceptorsAtRound(roundID)/2)

					return
				}

				continue
			}
		}
	}(roundID)
}

func (p *Proposer) UpdateListOfParticipantsOfTheRequiredRoles(acceptors []int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, acceptor := range acceptors {
		p.Acceptors[acceptor] = struct{}{}
	}
}

func (p *Proposer) Port() int {
	return p.port
}

func (p *Proposer) log(proposerHighestID domain_roles.HighestID, format string, v ...any) {
	p.logMu.Lock()
	log.Printf("[%d] %s", proposerHighestID, fmt.Sprintf(format, v...))
	p.logMu.Unlock()
}
