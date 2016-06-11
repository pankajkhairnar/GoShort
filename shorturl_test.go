package main

import (
	"testing"
	"fmt"
)

var testCases = []struct {
	in string
	out string
} {
	{"", "a"},
	{"aa", "ab"},
	{"a9", "ba"},
	{"99", "aaa"},
}

func TestCodeGeneration(t *testing.T) {

	for _, testCase := range testCases {
		expected := testCase.out
		response, _ := GenerateNextCode(testCase.in)
		if response != expected {
			fmt.Println("Expected :"+expected+ " Received :"+response)
			t.Error("Test failed")
		}
	}
}