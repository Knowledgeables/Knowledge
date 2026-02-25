package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

var (
	sessions = map[string]int64{} // sessionID -> userID
	mu       sync.Mutex
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
