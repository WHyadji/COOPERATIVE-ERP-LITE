package utils

import (
	"testing"
	"unicode"
)

func TestGenerateRandomPassword(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"Standard length 12", 12},
		{"Longer length 16", 16},
		{"Short length (should default to 12)", 6},
		{"Very long length", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GenerateRandomPassword(tt.length)
			if err != nil {
				t.Errorf("GenerateRandomPassword() error = %v", err)
				return
			}

			// Check minimum length (should be at least 8, defaults to 12 if less)
			expectedLength := tt.length
			if tt.length < 8 {
				expectedLength = 12
			}

			if len(password) != expectedLength {
				t.Errorf("GenerateRandomPassword() length = %d, want %d", len(password), expectedLength)
			}

			// Verify password contains at least one lowercase letter
			hasLower := false
			for _, c := range password {
				if unicode.IsLower(c) {
					hasLower = true
					break
				}
			}
			if !hasLower {
				t.Errorf("GenerateRandomPassword() = %s, should contain at least one lowercase letter", password)
			}

			// Verify password contains at least one uppercase letter
			hasUpper := false
			for _, c := range password {
				if unicode.IsUpper(c) {
					hasUpper = true
					break
				}
			}
			if !hasUpper {
				t.Errorf("GenerateRandomPassword() = %s, should contain at least one uppercase letter", password)
			}

			// Verify password contains at least one digit
			hasDigit := false
			for _, c := range password {
				if unicode.IsDigit(c) {
					hasDigit = true
					break
				}
			}
			if !hasDigit {
				t.Errorf("GenerateRandomPassword() = %s, should contain at least one digit", password)
			}

			// Verify password contains only valid characters
			for _, c := range password {
				if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
					t.Errorf("GenerateRandomPassword() = %s, contains invalid character: %c", password, c)
				}
			}
		})
	}
}

func TestGenerateRandomPassword_Uniqueness(t *testing.T) {
	// Generate 100 passwords and ensure they are all unique
	passwords := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		password, err := GenerateRandomPassword(12)
		if err != nil {
			t.Errorf("GenerateRandomPassword() error = %v on iteration %d", err, i)
			return
		}

		if passwords[password] {
			t.Errorf("GenerateRandomPassword() generated duplicate password: %s", password)
		}
		passwords[password] = true
	}

	if len(passwords) != iterations {
		t.Errorf("Expected %d unique passwords, got %d", iterations, len(passwords))
	}
}

func TestGenerateRandomPassword_NoPredictablePattern(t *testing.T) {
	// Generate passwords and ensure they don't follow predictable patterns
	iterations := 50
	patterns := []string{
		"123", "456", "789", "abc", "xyz", "password", "admin",
		"111", "222", "333", "aaa", "bbb",
	}

	for i := 0; i < iterations; i++ {
		password, err := GenerateRandomPassword(12)
		if err != nil {
			t.Errorf("GenerateRandomPassword() error = %v", err)
			return
		}

		// Check that password doesn't contain common predictable patterns
		for _, pattern := range patterns {
			// We just check it doesn't have too many occurrences of the pattern
			// (some random overlap is expected and OK)
			count := 0
			for j := 0; j <= len(password)-len(pattern); j++ {
				if password[j:j+len(pattern)] == pattern {
					count++
				}
			}
			if count > 1 {
				t.Errorf("GenerateRandomPassword() = %s, contains predictable pattern '%s' multiple times", password, pattern)
			}
		}
	}
}

func TestGenerateRandomPassword_NoUsernamePattern(t *testing.T) {
	// Ensure the new password generator doesn't create passwords
	// that follow the old vulnerable pattern: {username}123
	testUsernames := []string{"admin", "user", "test", "bendahara", "kasir"}

	for _, username := range testUsernames {
		password, err := GenerateRandomPassword(12)
		if err != nil {
			t.Errorf("GenerateRandomPassword() error = %v", err)
			return
		}

		// Old vulnerable pattern
		vulnerablePattern := username + "123"

		// Ensure the generated password is NOT the vulnerable pattern
		if password == vulnerablePattern {
			t.Errorf("GenerateRandomPassword() generated vulnerable password pattern: %s", password)
		}
	}
}

func TestRandomChar(t *testing.T) {
	charset := "abc123"

	// Generate 100 random characters
	for i := 0; i < 100; i++ {
		char, err := randomChar(charset)
		if err != nil {
			t.Errorf("randomChar() error = %v", err)
			return
		}

		// Verify the character is from the charset
		found := false
		for _, c := range charset {
			if byte(c) == char {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("randomChar() = %c, not in charset %s", char, charset)
		}
	}
}

func TestShuffle(t *testing.T) {
	original := []byte("abcdefghijklmnop")
	data := make([]byte, len(original))
	copy(data, original)

	err := shuffle(data)
	if err != nil {
		t.Errorf("shuffle() error = %v", err)
		return
	}

	// Length should be the same
	if len(data) != len(original) {
		t.Errorf("shuffle() changed length from %d to %d", len(original), len(data))
	}

	// All characters should still be present (frequency check)
	originalMap := make(map[byte]int)
	shuffledMap := make(map[byte]int)

	for _, c := range original {
		originalMap[c]++
	}
	for _, c := range data {
		shuffledMap[c]++
	}

	for k, v := range originalMap {
		if shuffledMap[k] != v {
			t.Errorf("shuffle() changed character frequencies")
		}
	}

	// Note: There's a very small chance the shuffle results in the same order,
	// but with 16 characters, the probability is 1/16! which is negligible
}
