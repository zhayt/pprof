package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const maxPasswordLength = 5

var chars = []string{"a", "b", "c", "d", "e", "f", "g"}

func BruteForcePassword(hash string) string {
	return bruteForceRecursively(hash, "")
}

func bruteForceRecursively(hash string, passwd string) string {
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

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func compareHash(a, b string) bool {
	return a == b
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
