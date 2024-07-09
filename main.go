package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const maxPasswordLength = 5

var chars = []string{"a", "b", "c", "d", "e", "f", "g"}

func BruteForcePassword(hash []byte) string {
	return bruteForceRecursively(hash, "")
}

func bruteForceRecursively(hash []byte, passwd string) string {
	if compareSliceOfBitesHash(hash, getMD5HashAsSliceOfBytes(passwd)) {
		return passwd
	}
	for i := 0; i < len(chars); i++ {
		if len(passwd) == maxPasswordLength {
			return ""
		}
		if str := bruteForceRecursively(hash, passwd+chars[i]); str != "" {
			return str
		}
	}
	return ""
}

func bruteForceLinear(hash []byte, passwd string) string {
	stack := []string{""}

	for len(stack) > 0 {
		passwd := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if compareSliceOfBitesHash(hash, getMD5HashAsSliceOfBytes(passwd)) {
			return passwd
		}

		if len(passwd) < maxPasswordLength {
			for i := len(chars) - 1; i >= 0; i-- {
				stack = append(stack, passwd+chars[i])
			}
		}
	}
	return ""
}

func bruteForceLinearReusingSpace(hash []byte, passwd string) string {
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
		if compareSliceOfBitesHash(hash, getMD5HashAsSliceOfBytes(string(buff[:pos]))) {
			return string(buff[:pos])
		}
	}

	return ""
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getMD5HashAsSliceOfBytes(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hasher.Sum(nil)
}

func compareHash(a, b string) bool {
	return a == b
}

func compareSliceOfBitesHash(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	passwords := []string{"a", "b", "c"}
	for i := 0; i < len(passwords); i++ {
		hash := getMD5HashAsSliceOfBytes(passwords[i])
		if passwd := BruteForcePassword(hash); passwd == passwords[i] {
			fmt.Printf("Find password: %s - %s\n", passwords[i], passwd)
		} else {
			fmt.Printf("Loss: want %s, got %s\n", passwords[i], passwd)
		}
	}
}
