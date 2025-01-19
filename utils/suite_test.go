package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestUtils(c *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(c, "utils test suite")
}
