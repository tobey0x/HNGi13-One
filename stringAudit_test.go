package main

import "testing"

// func TestStringLength(t *testing.T) {
// 	inputString := "string to analyze"

// 	actual := CheckStringLength(inputString)
// 	expected := 17

// 	if actual != expected {
// 		t.Errorf("Expected '%v', got '%v'", expected, actual)
// 	}
// }

func TestStringPalindromeFalse(t *testing.T) {
	inputString := "string to analyze"

	actual := isPalindrome(inputString)
	expected :=  false

	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}


func TestStringPalindromeTrue(t *testing.T) {
	inputString := "race car"

	actual := isPalindrome(inputString)
	expected :=  true

	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}


// func TestUniqueCharacters(t *testing.T) {
// 	inputString := "string To analyze"

// 	actual := UniqueCharaters(inputString)
// 	expected := 13

// 	if actual != expected {
// 		t.Errorf("Expected '%v', got '%v'", expected, actual)
// 	}
// }


func TestSHA256Hash(t *testing.T) {
	inputString := "boob"

	actual := computeSHA256Hash(inputString)
	expected := "e99d55248f67be7623332be8dabcd143ce2495eb923860b4ed0d963621ece901"

	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}
