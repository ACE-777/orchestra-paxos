package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	network "orchestra-paxos/internal/domain/network"
)

type Network struct {
	loss float32
}

func NewNetwork(loss float32) *Network {
	return &Network{
		loss: loss,
	}
}

func (n *Network) Unmarshal(in interface{}, out interface{}) error {
	if raw, err := json.Marshal(in); err != nil {
		return err
	} else {
		return json.Unmarshal(raw, &out)
	}
}

func (n *Network) Send(receiverID int, delay time.Duration, msg *network.NetworkMessage) {
	if rand.Float32() < n.loss {
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return
	}
	msgSize := len(buf.Bytes())

	time.Sleep(network.OneWayDelay)

	bandwidthBytesPerSec := 10_000_000 / 8 // = 1.25 MB/s = 1_250_000 bytes/sec
	transferTime := time.Duration(float64(msgSize) / float64(bandwidthBytesPerSec) * float64(time.Second))

	time.Sleep(transferTime)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", receiverID), time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return
		}
	}

	if err := json.NewEncoder(conn).Encode(msg); err != nil {
		return
	}
}
