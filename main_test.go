package main

import (
	"fmt"
	"testing"
)

func TestBruteForcePassword(t *testing.T) {
	var table = []struct {
		input string
	}{
		{input: "a"},
		{input: "ba"},
		{input: "cf"},
	}

	for _, tab := range table {
		t.Run(fmt.Sprintf("input_%s", tab.input), func(t *testing.T) {
			if got := BruteForcePassword(getMD5HashAsSliceOfBytes(tab.input)); got != tab.input {
				t.Errorf("Error want %s, got %s", tab.input, got)
			}
		})
	}
}

func BenchmarkBruteForcePassword(b *testing.B) {
	var table = []struct {
		input string
	}{
		{input: "a"},
		{input: "ba"},
		{input: "cf"},
	}
	for _, tab := range table {
		b.Run(fmt.Sprintf("input_%s", tab.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				BruteForcePassword(getMD5HashAsSliceOfBytes(tab.input))
			}
		})
	}
}
