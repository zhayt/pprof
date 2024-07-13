package main

import (
	"crypto/md5"
	"fmt"
)

const maxPasswordLength = 5

var chars = []string{"a", "b", "c", "d", "e", "f", "g"}

func BruteForcePassword(hash []byte) string {
	return bruteForceLinear(hash)
}

func bruteForceRecursively(hash []byte, passwd string) string {
	if compareHash(hash, getMD5Hash(passwd)) {
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

func bruteForceLinear(hash []byte) string {
	stack := []string{""}

	for len(stack) > 0 {
		length := len(stack)
		for _, passwd := range stack {
			if compareHash(hash, getMD5Hash(passwd)) {
				return passwd
			}

			if len(passwd) < maxPasswordLength {
				for i := 0; i < len(chars); i++ {
					stack = append(stack, passwd+chars[i])
				}
			}
		}
		stack = stack[length:]
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
