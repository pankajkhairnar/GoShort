package main

import (
	"errors"
	"bytes"
	"fmt"
)

var seedChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var seedCharsLen = len(seedChars)

func main() {
	var myStr = "a99"
	newString, err := getNextCode(myStr)
	fmt.Println(">>", myStr)
	fmt.Println(">>", newString, err)
}

func getNextCode(code string) (string, error) {
	if code == "" {
		return "a", nil
	}
	codeBytes := []byte(code)
	codeByteLen := len(codeBytes)

	codeCharIndex := -1
	for i := (codeByteLen - 1); i >= 0; i-- {
		codeCharIndex = bytes.IndexByte(seedChars, codeBytes[i])
		if codeCharIndex == -1 || codeCharIndex >= seedCharsLen {
			return "", errors.New("Invalid code")
		} else if codeCharIndex == (seedCharsLen - 1) {
			codeBytes[i] = 97
		} else {
			codeBytes[i] = seedChars[(codeCharIndex + 1)]
			return string(codeBytes), nil
		}
	}
	for _, byteVal := range codeBytes {
		if byteVal != 97 {
			return string(codeBytes), nil
		}
	}
	return "a" + string(codeBytes), nil
}
