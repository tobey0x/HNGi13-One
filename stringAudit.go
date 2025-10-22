package main

import (
	"log"
	"strings"
	"crypto/sha256"
	"encoding/hex"
)

func checkStringLength(text string) int {
	return len(text)
}

func is_palindrome(rawStr string) bool {
	str := strings.ToLower(rawStr)

	start, end := 0, len(str)-1

	for start < end {
		if str[start] != str[end] {
			return false
		}
		start++
		end--
	}
	log.Printf("%s is a palindrome", str)
	return true
}

func uniqueCharaters(rawStr string) int {
	str := strings.ToLower(rawStr)
	charSet := make(map[rune]bool)

	for _, char := range str {
		charSet[char] = true
	}

	return len(charSet)	
}


func createSHA256Hash(rawStr string) string {
	hash := sha256.New()

	hash.Write([]byte(rawStr))

	hashSum := hash.Sum(nil)
	hexHash := hex.EncodeToString(hashSum)
	
	return hexHash
}

func charFreqMap(rawStr string) map[rune]int {
	charCount := make(map[rune]int)

	for _, char := range rawStr {
		charCount[char]++
	}

	return charCount
}