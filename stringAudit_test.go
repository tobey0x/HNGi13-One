package main

import "testing"

func TestStringLength(t *testing.T) {
	inputString := "string to analyze"

	actual := checkStringLength(inputString)
	expected := 17

	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}

func TestStringPalindromeFalse(t *testing.T) {
	inputString := "string to analyze"

	actual := is_palindrome(inputString)
	expected :=  false

	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}


func TestStringPalindromeTrue(t *testing.T) {
	inputString := "boob"

	actual := is_palindrome(inputString)
	expected :=  true

	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}


func TestUniqueCharacters(t *testing.T) {
	inputString := "string To analyze"

	actual := uniqueCharaters(inputString)
	expected := 13

	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}


func TestSHA256Hash(t *testing.T) {
	inputString := "boob"

	actual := createSHA256Hash(inputString)
	expected := "e99d55248f67be7623332be8dabcd143ce2495eb923860b4ed0d963621ece901"

	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}
