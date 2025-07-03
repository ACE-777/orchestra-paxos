package values_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	roles "orchestra-paxos/internal/domain/roles"
	usecases "orchestra-paxos/internal/repository/usecases"
	values_uc "orchestra-paxos/internal/repository/usecases/values"
)

var _ = Describe("ValuesFromUser", func() {
	var (
		values      usecases.ValuesFromUsers
		operationID roles.HighestID = 1
	)

	BeforeEach(func() {
		values = values_uc.NewValuesFromUser()
	})

	Describe("AddValue method", func() {
		It("should add a value for a specific operationID", func() {
			values.AddValue("test value", operationID)
			Expect(values.ValueFromRound(operationID)).To(Equal("test value"))
		})

		It("should overwrite a value for the same operationID", func() {
			values.AddValue("first value", operationID)
			values.AddValue("second value", operationID)
			Expect(values.ValueFromRound(operationID)).To(Equal("second value"))
		})
	})

	Describe("ValueFromRound method", func() {
		It("should return an empty string if no value is set for the operationID", func() {
			Expect(values.ValueFromRound(operationID)).To(Equal(""))
		})

		It("should return the correct value for the operationID", func() {
			values.AddValue("testing value", operationID)
			Expect(values.ValueFromRound(operationID)).To(Equal("testing value"))
		})
	})
})
