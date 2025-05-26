package usecases

import (
	"orchestra-paxos/internal/domain/roles"
)

type AcceptorsStorage interface {
	AddAcceptor(acceptor string, roundID roles.HighestID)
	NumberOfAcceptorsAtRound(roundID roles.HighestID) int
	AllAcceptorsAtRound(roundID roles.HighestID) []string
}
type Timers interface {
	CheckExpireTimer(operationID roles.HighestID) bool
	SetExpireTimer(operationID roles.HighestID)
	InitExpireTimer(operationID roles.HighestID)
}
