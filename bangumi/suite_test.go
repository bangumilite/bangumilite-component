package bangumi

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBangumiClient(c *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(c, "Bangumi Client Test Suite")
}

var _ = BeforeEach(func() {
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})
