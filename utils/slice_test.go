package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("slice utils unit tests", func() {
	Describe("RemoveDuplicates", func() {
		It("should remove duplicate strings", func() {
			input := []string{"a", "b", "c", "c", "b"}
			want := []string{"a", "b", "c"}
			got := RemoveDuplicates(input)

			Expect(got).To(Equal(want))
		})

		It("should remove duplicate ints", func() {
			input := []int{1, 2, 3, 2, 1, 4}
			want := []int{1, 2, 3, 4}
			got := RemoveDuplicates(input)

			Expect(got).To(Equal(want))
		})

		It("should remove duplicate floats", func() {
			input := []float64{1.1, 2.2, 1.1, 3.3}
			want := []float64{1.1, 2.2, 3.3}
			got := RemoveDuplicates(input)

			Expect(got).To(Equal(want))
		})

		It("should not remove any elements if there are no duplicates", func() {
			input := []string{"a", "b", "c"}
			want := []string{"a", "b", "c"}
			got := RemoveDuplicates(input)

			Expect(got).To(Equal(want))
		})
	})
})
