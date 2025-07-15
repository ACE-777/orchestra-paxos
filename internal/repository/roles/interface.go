package roles

import (
	"sync"
)

type InitRoles interface {
	Run(wg *sync.WaitGroup)
	Port() int
	UpdateListOfParticipantsOfTheRequiredRoles(participantsOfTheRequiredRoles []int)
}
