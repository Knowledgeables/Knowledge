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

func Create(userID int64) string {
	b := make([]byte, 32)
	rand.Read(b)
	sessionID := hex.EncodeToString(b)

	mu.Lock()
	sessions[sessionID] = userID
	mu.Unlock()

	return sessionID
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
