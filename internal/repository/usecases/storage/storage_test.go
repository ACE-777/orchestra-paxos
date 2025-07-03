package storage_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"orchestra-paxos/internal/repository/usecases"

	roles "orchestra-paxos/internal/domain/roles"
	storage "orchestra-paxos/internal/repository/usecases/storage"
)

var _ = Describe("Acceptors storages", func() {
	Context("Alive Acceptors storage", func() {
		var (
			aliveAcceptors usecases.AcceptorsStorage
			roundID        roles.HighestID = 1
		)

		BeforeEach(func() {
			aliveAcceptors = storage.NewAliveAcceptors()
		})

		Describe("AddAcceptor method", func() {
			It("should add an acceptors for a given round", func() {
				acceptor := "acceptor1"
				aliveAcceptors.AddAcceptor(acceptor, roundID)

				Expect(aliveAcceptors.NumberOfAcceptorsAtRound(roundID)).To(Equal(1))
			})

			It("should allow adding multiple acceptors for the same round", func() {
				firstAcceptor := "acceptor1"
				secondAcceptor := "acceptor2"

				aliveAcceptors.AddAcceptor(firstAcceptor, roundID)
				aliveAcceptors.AddAcceptor(secondAcceptor, roundID)

				Expect(aliveAcceptors.NumberOfAcceptorsAtRound(roundID)).To(Equal(2))
			})
		})

		Describe("NumberOfAcceptorsAtRound method", func() {
			It("should return 0 if no acceptors have been added", func() {
				Expect(aliveAcceptors.NumberOfAcceptorsAtRound(roundID)).To(Equal(0))
			})

			It("should return the correct number of acceptors for a round", func() {
				firstAcceptor := "acceptor1"
				secondAcceptor := "acceptor2"

				aliveAcceptors.AddAcceptor(firstAcceptor, roundID)
				aliveAcceptors.AddAcceptor(secondAcceptor, roundID)

				Expect(aliveAcceptors.NumberOfAcceptorsAtRound(roundID)).To(Equal(2))
			})
		})

		Describe("AllAcceptorsAtRound method", func() {
			It("should return an empty slice if no acceptors are present", func() {
				Expect(aliveAcceptors.AllAcceptorsAtRound(roundID)).To(BeEmpty())
			})

			It("should return all acceptors added for the given round", func() {
				firstAcceptor := "acceptor1"
				secondAcceptor := "acceptor2"

				aliveAcceptors.AddAcceptor(firstAcceptor, roundID)
				aliveAcceptors.AddAcceptor(secondAcceptor, roundID)

				acceptors := aliveAcceptors.AllAcceptorsAtRound(roundID)

				Expect(acceptors).To(ConsistOf(firstAcceptor, secondAcceptor))
			})
		})
	})

	Context("Accepted Acceptors storage", func() {
		var (
			acceptedAcceptors usecases.AcceptorsStorage
			roundID           roles.HighestID = 1
		)

		BeforeEach(func() {
			acceptedAcceptors = storage.NewAcceptedAcceptors()
		})

		Describe("AddAcceptor method", func() {
			It("should add an acceptors for a given round", func() {
				acceptor := "acceptor1"
				acceptedAcceptors.AddAcceptor(acceptor, roundID)

				Expect(acceptedAcceptors.NumberOfAcceptorsAtRound(roundID)).To(Equal(1))
			})

			It("should allow adding multiple acceptors for the same round", func() {
				firstAcceptor := "acceptor1"
				secondAcceptor := "acceptor2"

				acceptedAcceptors.AddAcceptor(firstAcceptor, roundID)
				acceptedAcceptors.AddAcceptor(secondAcceptor, roundID)

				Expect(acceptedAcceptors.NumberOfAcceptorsAtRound(roundID)).To(Equal(2))
			})
		})

		Describe("NumberOfAcceptorsAtRound method", func() {
			It("should return 0 if no acceptors have been added", func() {
				Expect(acceptedAcceptors.NumberOfAcceptorsAtRound(roundID)).To(Equal(0))
			})

			It("should return the correct number of acceptors for a round", func() {
				firstAcceptor := "acceptor1"
				secondAcceptor := "acceptor2"

				acceptedAcceptors.AddAcceptor(firstAcceptor, roundID)
				acceptedAcceptors.AddAcceptor(secondAcceptor, roundID)

				Expect(acceptedAcceptors.NumberOfAcceptorsAtRound(roundID)).To(Equal(2))
			})
		})

		Describe("AllAcceptorsAtRound method", func() {
			It("should return an empty slice if no acceptors are present", func() {
				Expect(acceptedAcceptors.AllAcceptorsAtRound(roundID)).To(BeEmpty())
			})

			It("should return all acceptors added for the given round", func() {
				firstAcceptor := "acceptor1"
				secondAcceptor := "acceptor2"

				acceptedAcceptors.AddAcceptor(firstAcceptor, roundID)
				acceptedAcceptors.AddAcceptor(secondAcceptor, roundID)

				acceptors := acceptedAcceptors.AllAcceptorsAtRound(roundID)

				Expect(acceptors).To(ConsistOf(firstAcceptor, secondAcceptor))
			})
		})
	})
})
