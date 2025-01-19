package utils

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicatesStrings(t *testing.T) {
	input := []string{"a", "b", "c", "c", "b"}
	want := []string{"a", "b", "c"}

	got := RemoveDuplicates(input)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestRemoveDuplicatesInts(t *testing.T) {
	input := []int{1, 2, 3, 2, 1, 4}
	want := []int{1, 2, 3, 4}

	got := RemoveDuplicates(input)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestRemoveDuplicatesFloats(t *testing.T) {
	input := []float64{1.1, 2.2, 1.1, 3.3}
	want := []float64{1.1, 2.2, 3.3}

	got := RemoveDuplicates(input)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestRemoveDuplicatesEmptySlice(t *testing.T) {
	var input []int
	var want []int

	got := RemoveDuplicates(input)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestRemoveDuplicatesNoDuplicates(t *testing.T) {
	input := []string{"a", "b", "c"}
	want := []string{"a", "b", "c"}

	got := RemoveDuplicates(input)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestRemoveDuplicatesSingleElement(t *testing.T) {
	input := []int{42}
	want := []int{42}

	got := RemoveDuplicates(input)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}
