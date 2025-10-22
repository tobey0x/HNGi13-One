package main

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type StringProperties struct {
	Length				int			`json:"length" gorm:"-"`
	IsPalindrome		bool		`json:"is_palindrome" gorm:"-"`
	UniqueCharacters	int			`json:"unique_characters" gorm:"-"`
	WordCount			int			`json:"word_count" gorm:"-"`
	Sha256Hash			string		`json:"sha256_hash" gorm:"-"`
	CharFreqMap			CharacterFreqMap `json:"character_frequency_map" gorm:"type:text"`
}


func analyzeString(s string) StringProperties {
	uniqueCharacters := len(CharFreqMap(s))
	wordCount := len(strings.Fields(s))



	return StringProperties{
		Length: len(s),
		IsPalindrome: isPalindrome(s),
		UniqueCharacters: uniqueCharacters,
		WordCount: wordCount,
		Sha256Hash: computeSHA256Hash(s),
		CharFreqMap: CharFreqMap(s),
	}
}

func isPalindrome(rawStr string) bool {
	str := strings.ToLower(strings.ReplaceAll(rawStr, " ", ""))

	start, end := 0, len(str)-1

	for start < end {
		if str[start] != str[end] {
			return false
		}
		start++
		end--
	}
	return true
}


func computeSHA256Hash(rawStr string) string {
	hash := sha256.New()

	hash.Write([]byte(rawStr))

	hashSum := hash.Sum(nil)
	hexHash := hex.EncodeToString(hashSum)
	
	return hexHash
}

func CharFreqMap(rawStr string) map[string]int {
	charCount := make(map[string]int)

	for _, char := range rawStr {
		charCount[string(char)]++
	}

	return charCount
}

type CharacterFreqMap map[string]int

func (c CharacterFreqMap) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}

	data, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CharacterFreqMap: %w", err)
	}

	return data, nil
}

func (c *CharacterFreqMap) Scan(value interface{}) error {
	if value == nil {
		*c = make(CharacterFreqMap)
		return  nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("unsupported type for scanning CharacterFreqMap")
	}

	return json.Unmarshal(data, c)
}