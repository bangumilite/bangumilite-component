package season

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("season seasons unit tests", func() {
	It("should return valid spring season", func() {
		got := spring(2025)

		Expect(got.ID()).To(Equal("202504"))
		Expect(got.Year()).To(Equal(2025))
		Expect(got.Name()).To(Equal(SpringName))
	})

	It("should return valid summer season", func() {
		got := summer(2025)

		Expect(got.ID()).To(Equal("202507"))
		Expect(got.Year()).To(Equal(2025))
		Expect(got.Name()).To(Equal(SummerName))
	})

	It("should return valid autumn season", func() {
		got := autumn(2025)

		Expect(got.ID()).To(Equal("202510"))
		Expect(got.Year()).To(Equal(2025))
		Expect(got.Name()).To(Equal(AutumnName))
	})

	It("should return valid winter season", func() {
		got := winter(2025)

		Expect(got.ID()).To(Equal("202501"))
		Expect(got.Year()).To(Equal(2025))
		Expect(got.Name()).To(Equal(WinterName))
	})
})
