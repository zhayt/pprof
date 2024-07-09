package main

import (
	"crypto/md5"
	"fmt"
)

const maxPasswordLength = 5

var chars = []string{"a", "b", "c", "d", "e", "f", "g"}

func BruteForcePassword(hash []byte) string {
	return bruteForceLinearReusingSpace(hash)
}

func bruteForceLinearReusingSpace(hash []byte) string {
	buff := make([]byte, maxPasswordLength)
	state := make([]int, maxPasswordLength)

	var pos int
	for pos != -1 {
		if pos == maxPasswordLength {
			pos--
			continue
		}
		if state[pos] == len(chars)-1 {
			state[pos] = 0
			pos--
			continue
		}
		buff[pos] = chars[state[pos]][0]
		state[pos]++
		pos++
		if compareHash(hash, getMD5Hash(string(buff[:pos]))) {
			return string(buff[:pos])
		}
	}

	return ""
}

func getMD5Hash(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hasher.Sum(nil)
}

func compareHash(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	passwords := []string{"a", "b", "c"}
	for i := 0; i < len(passwords); i++ {
		hash := getMD5Hash(passwords[i])
		if passwd := BruteForcePassword(hash); passwd == passwords[i] {
			fmt.Printf("Find password: %s - %s\n", passwords[i], passwd)
		} else {
			fmt.Printf("Loss: want %s, got %s\n", passwords[i], passwd)
		}
	}
}
