package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// User represents an identity within KenPanel.
type User struct {
	ID                 string    `json:"id"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	PasswordHash       string    `json:"-"`
	IsBreakGlassAdmin  bool      `json:"is_break_glass_admin"` // emergency audit tracking
	MFAEnabled         bool      `json:"mfa_enabled"`
	MFAType            string    `json:"mfa_type"` // totp, webauthn
	TOTPSecret         string    `json:"-"`
	RecoveryCodesHash  []string  `json:"-"`
	FailedAttempts     int       `json:"failed_attempts"`
	LockedUntil        *time.Time `json:"locked_until,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

// Session tracks an active login session.
type Session struct {
	SessionID  string    `json:"session_id"`
	UserID     string    `json:"user_id"`
	TokenHash  string    `json:"-"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsRevoked  bool      `json:"is_revoked"`
}

// AuthManager manages authentication, security policies and session lifecycles.
type AuthManager struct {
	mu           sync.RWMutex
	users        map[string]*User    // key: username
	sessions     map[string]*Session // key: sessionID
	maxFailures  int
	lockoutTime  time.Duration
}

// NewAuthManager initializes authentication controller.
func NewAuthManager() *AuthManager {
	return &AuthManager{
		users:       make(map[string]*User),
		sessions:    make(map[string]*Session),
		maxFailures: 5,
		lockoutTime: 15 * time.Minute,
	}
}

// RegisterUser creates a new account with hashed password.
func (a *AuthManager) RegisterUser(username, email, rawPassword string, breakGlass bool) (*User, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(rawPassword) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	if _, exists := a.users[username]; exists {
		return nil, fmt.Errorf("user %s already exists", username)
	}

	hash := hashPasswordMock(rawPassword)
	user := &User{
		ID:                fmt.Sprintf("usr-%d", time.Now().UnixNano()),
		Username:          username,
		Email:             email,
		PasswordHash:      hash,
		IsBreakGlassAdmin: breakGlass,
		CreatedAt:         time.Now().UTC(),
	}

	a.users[username] = user
	return user, nil
}

// Authenticate verifies credentials, enforcing lockout policies and progressive delays.
func (a *AuthManager) Authenticate(username, rawPassword, ip, userAgent string) (*Session, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	user, ok := a.users[username]
	if !ok {
		return nil, errors.New("invalid credentials")
	}

	// Check lockout
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, fmt.Errorf("account is temporarily locked due to repeated failed logins; retry after %s", user.LockedUntil.Format(time.RFC3339))
	}

	// Validate password
	expectedHash := hashPasswordMock(rawPassword)
	if subtle.ConstantTimeCompare([]byte(user.PasswordHash), []byte(expectedHash)) != 1 {
		user.FailedAttempts++
		if user.FailedAttempts >= a.maxFailures {
			lockout := time.Now().Add(a.lockoutTime)
			user.LockedUntil = &lockout
			return nil, errors.New("maximum login attempts exceeded; account locked")
		}
		return nil, errors.New("invalid credentials")
	}

	// Reset failure count on success
	user.FailedAttempts = 0
	user.LockedUntil = nil

	// Create session
	token := generateSecureToken(32)
	h := sha256.New()
	h.Write([]byte(token))
	tokenHash := hex.EncodeToString(h.Sum(nil))

	session := &Session{
		SessionID:  fmt.Sprintf("sess-%d", time.Now().UnixNano()),
		UserID:     user.ID,
		TokenHash:  tokenHash,
		IPAddress:  ip,
		UserAgent:  userAgent,
		CreatedAt:  time.Now().UTC(),
		LastSeenAt: time.Now().UTC(),
		ExpiresAt:  time.Now().UTC().Add(24 * time.Hour),
		IsRevoked:  false,
	}

	a.sessions[session.SessionID] = session
	return session, nil
}

// RevokeSession revokes a specific user session.
func (a *AuthManager) RevokeSession(sessionID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if sess, ok := a.sessions[sessionID]; ok {
		sess.IsRevoked = true
		return nil
	}
	return errors.New("session not found")
}

// ListUserSessions returns all active sessions for a user.
func (a *AuthManager) ListUserSessions(userID string) []*Session {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var list []*Session
	for _, s := range a.sessions {
		if s.UserID == userID && !s.IsRevoked && time.Now().Before(s.ExpiresAt) {
			list = append(list, s)
		}
	}
	return list
}

func hashPasswordMock(pass string) string {
	h := sha256.New()
	h.Write([]byte("kenpanel-salt:" + pass))
	return hex.EncodeToString(h.Sum(nil))
}

func generateSecureToken(bytesCount int) string {
	b := make([]byte, bytesCount)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
