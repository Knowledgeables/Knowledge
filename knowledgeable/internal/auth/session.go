package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

var (
	sessions = map[string]int64{}
	mu       sync.RWMutex
)

func Create(userID int64) (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	sessionID := hex.EncodeToString(b)

	mu.Lock()
	sessions[sessionID] = userID
	mu.Unlock()

	return sessionID, nil
}

func Get(sessionID string) (int64, bool) {
	mu.Lock()
	defer mu.Unlock()

	id, ok := sessions[sessionID]
	return id, ok
}

func Delete(sessionID string) {
	mu.Lock()
	delete(sessions, sessionID)
	mu.Unlock()

}

func ResolveFromRequest(r *http.Request) (*int64, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, true
	}

	id, ok := Get(cookie.Value)
	if !ok {
		return nil, true
	}

	return &id, false
}
