package auth

import (
	"testing"
)

func TestAuthenticationFlow(t *testing.T) {
	mgr := NewAuthManager()

	user, err := mgr.RegisterUser("testuser", "user@example.com", "SecurePassword123!", false)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	// Successful auth
	sess, err := mgr.Authenticate("testuser", "SecurePassword123!", "127.0.0.1", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("unexpected auth error: %v", err)
	}
	if sess.UserID != user.ID {
		t.Errorf("user ID mismatch")
	}

	// Failed auth attempts triggering lockout
	for i := 0; i < 5; i++ {
		_, _ = mgr.Authenticate("testuser", "WrongPassword", "127.0.0.1", "Mozilla/5.0")
	}

	// Verify locked out
	_, err = mgr.Authenticate("testuser", "SecurePassword123!", "127.0.0.1", "Mozilla/5.0")
	if err == nil {
		t.Errorf("expected account to be locked out")
	}
}
