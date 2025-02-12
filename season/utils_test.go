package season

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("season seasons unit tests", func() {
	Describe("New", func() {
		It("should return winter season", func() {
			t := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
			got := New(t)

			want := winter(2024)
			Expect(got.ID()).To(Equal(want.ID()))
			Expect(got.Name()).To(Equal(want.Name()))
			Expect(got.Year()).To(Equal(want.Year()))
		})

		It("should return spring season", func() {
			t := time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC)
			got := New(t)

			want := spring(2024)
			Expect(got.ID()).To(Equal(want.ID()))
			Expect(got.Name()).To(Equal(want.Name()))
			Expect(got.Year()).To(Equal(want.Year()))
		})

		It("should return summer season", func() {
			t := time.Date(2024, 7, 13, 0, 0, 0, 0, time.UTC)
			got := New(t)

			want := summer(2024)
			Expect(got.ID()).To(Equal(want.ID()))
			Expect(got.Name()).To(Equal(want.Name()))
			Expect(got.Year()).To(Equal(want.Year()))
		})

		It("should return autumn season", func() {
			t := time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC)
			got := New(t)

			want := autumn(2024)
			Expect(got.ID()).To(Equal(want.ID()))
			Expect(got.Name()).To(Equal(want.Name()))
			Expect(got.Year()).To(Equal(want.Year()))
		})
	})

	Describe("Next", func() {
		It("should return winter's next season, spring", func() {
			t := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
			got := New(t).Next()

			want := spring(2024)
			Expect(got.ID()).To(Equal(want.ID()))
			Expect(got.Name()).To(Equal(want.Name()))
			Expect(got.Year()).To(Equal(want.Year()))
		})

		It("should return spring's next season, summer", func() {
			t := time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC)
			got := New(t).Next()

			want := summer(2024)
			Expect(got.ID()).To(Equal(want.ID()))
			Expect(got.Name()).To(Equal(want.Name()))
			Expect(got.Year()).To(Equal(want.Year()))
		})

		It("should return summer's next season, autumn", func() {
			t := time.Date(2024, 7, 13, 0, 0, 0, 0, time.UTC)
			got := New(t).Next()

			want := autumn(2024)
			Expect(got.ID()).To(Equal(want.ID()))
			Expect(got.Name()).To(Equal(want.Name()))
			Expect(got.Year()).To(Equal(want.Year()))
		})

		It("should return autumn's next season, winter", func() {
			t := time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC)
			got := New(t).Next()

			want := winter(2025)
			Expect(got.ID()).To(Equal(want.ID()))
			Expect(got.Name()).To(Equal(want.Name()))
			Expect(got.Year()).To(Equal(want.Year()))
		})
	})
})
