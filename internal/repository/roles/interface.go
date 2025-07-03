package roles

import (
	"sync"
)

type InitRoles interface {
	Run(wg *sync.WaitGroup)
	Name() string
	UpdateListOfParticipantsOfTheRequiredRoles(participantsOfTheRequiredRoles []string)
}
