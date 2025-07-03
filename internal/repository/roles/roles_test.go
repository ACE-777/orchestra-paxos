package roles_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	network "orchestra-paxos/internal/repository/network"
	roles "orchestra-paxos/internal/repository/roles"
	acceptor "orchestra-paxos/internal/repository/roles/acceptor"
	learner "orchestra-paxos/internal/repository/roles/learner"
	proposer "orchestra-paxos/internal/repository/roles/proposer"
)

var _ = Describe("Roles", func() {
	var (
		mockNet network.NetworkActions
		role    roles.InitRoles
	)

	Context("Proposer", func() {
		BeforeEach(func() {
			mockNet = network.NewNetwork(0)
			role = proposer.NewProposer(1, 1, mockNet)
		})

		Describe("UpdateListOfParticipantsOfTheRequiredRoles method", func() {
			It("should update the Learners list", func() {
				participants := []string{"Acceptor1", "Acceptor2"}
				role.UpdateListOfParticipantsOfTheRequiredRoles(participants)

				proposerType, ok := role.(*proposer.Proposer)
				Expect(ok).To(BeTrue())
				Expect(proposerType.Acceptors).To(HaveLen(2))
				Expect(proposerType.Acceptors).To(HaveKey("Acceptor1"))
				Expect(proposerType.Acceptors).To(HaveKey("Acceptor2"))
			})
		})

		Describe("Name method", func() {
			It("should return the correct name", func() {
				expectedName := fmt.Sprintf("Proposer %d", 1)
				Expect(role.Name()).To(Equal(expectedName))
			})
		})
	})

	Context("Acceptor", func() {
		BeforeEach(func() {
			mockNet = network.NewNetwork(0)
			role = acceptor.NewAcceptor(1, 1, mockNet)
		})

		Describe("UpdateListOfParticipantsOfTheRequiredRoles method", func() {
			It("should update the Learners list", func() {
				participants := []string{"Learner1", "Learner2"}
				role.UpdateListOfParticipantsOfTheRequiredRoles(participants)

				acceptorType, ok := role.(*acceptor.Acceptor)
				Expect(ok).To(BeTrue())
				Expect(acceptorType.Learners).To(HaveLen(2))
				Expect(acceptorType.Learners).To(HaveKey("Learner1"))
				Expect(acceptorType.Learners).To(HaveKey("Learner1"))
			})
		})

		Describe("Name method", func() {
			It("should return the correct name", func() {
				expectedName := fmt.Sprintf("Acceptor %d", 1)
				Expect(role.Name()).To(Equal(expectedName))
			})
		})
	})

	Context("Learner", func() {
		BeforeEach(func() {
			mockNet = network.NewNetwork(0)
			role = learner.NewLearner(1, 1, mockNet)
		})
		Describe("Name method", func() {
			It("should return the correct name", func() {
				expectedName := fmt.Sprintf("Learner %d", 1)
				Expect(role.Name()).To(Equal(expectedName))
			})
		})
	})
})
