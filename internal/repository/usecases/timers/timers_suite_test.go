package timers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTimers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timers Suite")
}
