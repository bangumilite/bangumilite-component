package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("string_utils unit tests", func() {
	similarity := 90.0
	Describe("CalculateSimilarity", func() {
		It("should return 100.0 for two empty strings", func() {
			a := ""
			b := ""

			got := IsMatch(a, b, similarity)

			Expect(got).To(BeTrue())
		})

		It("should return 100.0 for same string", func() {
			a := "薬屋のひとりごと 第2期"
			b := "薬屋のひとりごと 第2期"

			got := IsMatch(a, b, similarity)

			Expect(got).To(BeTrue())
		})

		It("should return 100.0 for same string", func() {
			a := "進撃の巨人 Season3"
			b := "進撃の巨人 SeasonCN 3"

			got := IsMatch(a, b, similarity)

			Expect(got).To(BeTrue())
		})
	})
})
