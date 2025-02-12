package season

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestUtils(c *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(c, "season test suite")
}
