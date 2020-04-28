package main

import (
	"testing"
	"fmt"
	//"net/http"
	//"strings"
	//"net/http/httptest"
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

//func TestCreateShortUrl(t *testing.T) {
//	req, _ := http.NewRequest(
//		"POST",
//		"http://localhost:8080/create/",
//		strings.NewReader("url=http://coditas.com"),
//	)
//
//	recorder := httptest.NewRecorder()
//	Create(recorder, req, nil)
//	if recorder.Code != http.StatusOK {
//		t.Errorf("Expected status 200 Received %d", recorder.Code)
//	}
//
//	if !strings.Contains(recorder.Body.String(), "SUCCESS") {
//		t.Error("Unexpected body in response %q", recorder.Body.String())
//	}
//}