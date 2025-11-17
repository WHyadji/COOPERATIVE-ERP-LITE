package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	// Character sets for password generation
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	allChars       = lowercaseChars + uppercaseChars + digitChars
)

// GenerateRandomPassword generates a cryptographically secure random password
// with the specified length. The password will contain a mix of lowercase,
// uppercase, and digit characters.
//
// Parameters:
//   - length: The desired length of the password (minimum 12 recommended)
//
// Returns:
//   - A random password string
//   - An error if random generation fails
func GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 12 // Default to 12 if too short
	}

	password := make([]byte, length)

	// Ensure at least one character from each set
	// This guarantees the password has good complexity
	var err error
	password[0], err = randomChar(lowercaseChars)
	if err != nil {
		return "", err
	}
	password[1], err = randomChar(uppercaseChars)
	if err != nil {
		return "", err
	}
	password[2], err = randomChar(digitChars)
	if err != nil {
		return "", err
	}

	// Fill the rest with random characters from all sets
	for i := 3; i < length; i++ {
		password[i], err = randomChar(allChars)
		if err != nil {
			return "", err
		}
	}

	// Shuffle the password to avoid predictable patterns
	// (e.g., always starting with lowercase)
	if err := shuffle(password); err != nil {
		return "", err
	}

	return string(password), nil
}

// randomChar returns a random character from the given charset
// using crypto/rand for cryptographic security
func randomChar(charset string) (byte, error) {
	max := big.NewInt(int64(len(charset)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return charset[n.Int64()], nil
}

// shuffle randomly shuffles the bytes in the slice using Fisher-Yates algorithm
// with cryptographically secure random numbers
func shuffle(data []byte) error {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		max := big.NewInt(int64(i + 1))
		j, err := rand.Int(rand.Reader, max)
		if err != nil {
			return err
		}
		jInt := int(j.Int64())
		data[i], data[jInt] = data[jInt], data[i]
	}
	return nil
}
