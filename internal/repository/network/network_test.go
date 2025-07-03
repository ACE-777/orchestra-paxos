package network_test

import (
	"orchestra-paxos/internal/domain/roles"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	network_domain "orchestra-paxos/internal/domain/network"
	network "orchestra-paxos/internal/repository/network"
)

const (
	delay = 100 * time.Nanosecond
)

var _ = Describe("Network", func() {
	var (
		net        *network.Network
		msg        network_domain.NetworkMessage
		receiverID string
	)

	BeforeEach(func() {
		net = network.NewNetwork(0)
		msg = network_domain.NetworkMessage{Stage: roles.REQUEST,
			Sender: "client",
			Data: network_domain.MessageRequest{
				Value: "SecondValue",
			},
		}
		receiverID = "testReceiver"
	})

	Describe("Send and receive methods", func() {
		It("should send a message to the correct receiver", func() {
			net.Send(receiverID, 0, msg)
			receivedMsg := net.Receive(receiverID)
			Expect(receivedMsg).To(Equal(msg))
		})

		It("should not send a message if lost", func() {
			netWithLoss := network.NewNetwork(18446744073709551615)
			netWithLoss.Send(receiverID, 0, msg)
			ch := make(chan network_domain.NetworkMessage, 1)
			timer := time.NewTimer(time.Second)
			Eventually(func(g Gomega) {
				select {
				case receivedMsg := <-ch:
					g.Expect(receivedMsg).Should(Equal(msg))
				case <-timer.C:
				}
			}).WithTimeout(2 * time.Second).WithPolling(100 * time.Microsecond).Should(Succeed())
		})

		It("should respect send delay", func() {
			start := time.Now()
			net.Send(receiverID, delay, msg)
			receivedMsg := net.Receive(receiverID)

			Expect(receivedMsg).To(Equal(msg))
			Expect(time.Since(start)).To(BeNumerically(">", delay))
		})
	})
})
