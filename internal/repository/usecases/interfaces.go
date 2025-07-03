package usecases

import (
	roles "orchestra-paxos/internal/domain/roles"
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

type ValuesFromUsers interface {
	AddValue(value string, operationID roles.HighestID)
	ValueFromRound(operationID roles.HighestID) string
}
