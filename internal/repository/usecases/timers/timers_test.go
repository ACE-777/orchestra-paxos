package timers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"orchestra-paxos/internal/repository/usecases"

	roles "orchestra-paxos/internal/domain/roles"
	timers_uc "orchestra-paxos/internal/repository/usecases/timers"
)

var _ = Describe("Timers", func() {
	Context("Accept timer", func() {
		var (
			timersOfCollectingAcceptFromAcceptors usecases.Timers
			operationID                           roles.HighestID = 1
		)

		BeforeEach(func() {
			timersOfCollectingAcceptFromAcceptors = timers_uc.NewTimersOfCollectingAcceptFromAcceptors()
		})

		Describe("CheckExpireTimer method", func() {
			It("should return false if timer has not been set", func() {
				Expect(timersOfCollectingAcceptFromAcceptors.CheckExpireTimer(operationID)).To(BeFalse())
			})

			It("should return true if timer has been set", func() {
				timersOfCollectingAcceptFromAcceptors.SetExpireTimer(operationID)
				Expect(timersOfCollectingAcceptFromAcceptors.CheckExpireTimer(operationID)).To(BeTrue())
			})
		})

		Describe("InitExpireTimer method", func() {
			It("should initialize the timer with false if not already set", func() {
				timersOfCollectingAcceptFromAcceptors.InitExpireTimer(operationID)
				Expect(timersOfCollectingAcceptFromAcceptors.CheckExpireTimer(operationID)).To(BeFalse())
			})

			It("should not change the timer if it has already been initialized", func() {
				timersOfCollectingAcceptFromAcceptors.SetExpireTimer(operationID)
				timersOfCollectingAcceptFromAcceptors.InitExpireTimer(operationID)
				Expect(timersOfCollectingAcceptFromAcceptors.CheckExpireTimer(operationID)).To(BeTrue())
			})
		})
	})

	Context("Prepare timer", func() {
		var (
			timersOfCollectingPrepareFromAcceptors usecases.Timers
			operationID                            roles.HighestID = 1
		)

		BeforeEach(func() {
			timersOfCollectingPrepareFromAcceptors = timers_uc.NewTimersOfCollectingPrepareFromAcceptors()
		})

		Describe("CheckExpireTimer method", func() {
			It("should return false if timer has not been set", func() {
				Expect(timersOfCollectingPrepareFromAcceptors.CheckExpireTimer(operationID)).To(BeFalse())
			})

			It("should return true if timer has been set", func() {
				timersOfCollectingPrepareFromAcceptors.SetExpireTimer(operationID)
				Expect(timersOfCollectingPrepareFromAcceptors.CheckExpireTimer(operationID)).To(BeTrue())
			})
		})

		Describe("InitExpireTimer method", func() {
			It("should initialize the timer with false if not already set", func() {
				timersOfCollectingPrepareFromAcceptors.InitExpireTimer(operationID)
				Expect(timersOfCollectingPrepareFromAcceptors.CheckExpireTimer(operationID)).To(BeFalse())
			})

			It("should not change the timer if it has already been initialized", func() {
				timersOfCollectingPrepareFromAcceptors.SetExpireTimer(operationID)
				timersOfCollectingPrepareFromAcceptors.InitExpireTimer(operationID)
				Expect(timersOfCollectingPrepareFromAcceptors.CheckExpireTimer(operationID)).To(BeTrue())
			})
		})
	})
})
