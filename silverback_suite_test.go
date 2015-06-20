package silverback_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSilverback(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Silverback Suite")
}
