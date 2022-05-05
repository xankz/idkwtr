package idkwtr

import "testing"

func TestGeneratePassword(t *testing.T) {
	// Passwords must be > 3 characters at least.
	if _, err := GeneratePassword("sh"); err == nil {
		t.Errorf("GeneratePassword(\"sh\") err = nil, want err")
	}
}
