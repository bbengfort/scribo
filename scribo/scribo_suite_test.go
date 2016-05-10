package scribo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestScribo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scribo Suite")
}
