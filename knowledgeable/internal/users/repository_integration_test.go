//go:build integration

package users

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		t.Skip("DATABASE_URL not set — skipping integration test")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("cannot reach database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func insertTestUser(t *testing.T, db *sql.DB, repo *Repository, username, email string) *User {
	t.Helper()
	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: "testhash",
	}
	if err := repo.Register(user); err != nil {
		t.Fatalf("insertTestUser: %v", err)
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = $1", user.ID)
	})
	return user
}

func uniqueName(base string) string {
	return fmt.Sprintf("%s_%d", base, time.Now().UnixNano())
}

func TestRepository_Register_FindByUsername(t *testing.T) {
	db := openTestDB(t)
	repo := NewRepository(db)

	name := uniqueName("reg")
	email := name + "@test.com"
	user := insertTestUser(t, db, repo, name, email)

	if user.ID == 0 {
		t.Fatal("expected non-zero ID after register")
	}

	found, err := repo.FindByUsername(name)
	if err != nil {
		t.Fatal(err)
	}
	if found == nil {
		t.Fatal("user not found by username")
	}
	if found.Username != name {
		t.Errorf("want username %q, got %q", name, found.Username)
	}
	if found.Email != email {
		t.Errorf("want email %q, got %q", email, found.Email)
	}
	if found.PasswordHash != "testhash" {
		t.Error("password hash not round-tripped correctly")
	}
}

func TestRepository_Register_DuplicateUsername(t *testing.T) {
	db := openTestDB(t)
	repo := NewRepository(db)

	name := uniqueName("dup")
	insertTestUser(t, db, repo, name, name+"@test.com")

	dup := &User{
		Username:     name,
		Email:        "other_" + name + "@test.com",
		PasswordHash: "hash",
	}
	err := repo.Register(dup)
	if err != ErrUsernameTaken {
		t.Fatalf("expected ErrUsernameTaken, got %v", err)
	}
}

func TestRepository_FindByEmail(t *testing.T) {
	db := openTestDB(t)
	repo := NewRepository(db)

	name := uniqueName("email")
	email := name + "@test.com"
	inserted := insertTestUser(t, db, repo, name, email)

	found, err := repo.FindByEmail(email)
	if err != nil {
		t.Fatal(err)
	}
	if found == nil {
		t.Fatal("user not found by email")
	}
	if found.ID != inserted.ID {
		t.Errorf("want ID %d, got %d", inserted.ID, found.ID)
	}
}

func TestRepository_FindById(t *testing.T) {
	db := openTestDB(t)
	repo := NewRepository(db)

	name := uniqueName("byid")
	inserted := insertTestUser(t, db, repo, name, name+"@test.com")

	found, err := repo.FindById(inserted.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found == nil {
		t.Fatal("user not found by ID")
	}
	if found.ID != inserted.ID {
		t.Errorf("want ID %d, got %d", inserted.ID, found.ID)
	}

	notFound, err := repo.FindById(-1)
	if err != nil {
		t.Fatal(err)
	}
	if notFound != nil {
		t.Fatal("expected nil for unknown ID")
	}
}

func TestRepository_UpdatePassword(t *testing.T) {
	db := openTestDB(t)
	repo := NewRepository(db)

	name := uniqueName("pw")
	user := insertTestUser(t, db, repo, name, name+"@test.com")

	newHash := "newhash_updated"
	if err := repo.UpdatePassword(user.ID, newHash); err != nil {
		t.Fatal(err)
	}

	found, err := repo.FindById(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.PasswordHash != newHash {
		t.Errorf("want hash %q, got %q", newHash, found.PasswordHash)
	}
	if found.ShouldChangePassword {
		t.Error("expected ShouldChangePassword to be cleared after UpdatePassword")
	}
}

func TestRepository_PasswordResetToken_CreateAndFind(t *testing.T) {
	db := openTestDB(t)
	repo := NewRepository(db)

	name := uniqueName("token")
	user := insertTestUser(t, db, repo, name, name+"@test.com")

	tokenHash := fmt.Sprintf("hash_%d", time.Now().UnixNano())
	expiresAt := time.Now().Add(time.Hour)

	if err := repo.CreatePasswordResetToken(user.ID, tokenHash, expiresAt); err != nil {
		t.Fatal(err)
	}

	token, err := repo.FindPasswordResetToken(tokenHash)
	if err != nil {
		t.Fatal(err)
	}
	if token == nil {
		t.Fatal("token not found")
	}
	if token.UserID != user.ID {
		t.Errorf("want userID %d, got %d", user.ID, token.UserID)
	}
	if token.UsedAt != nil {
		t.Error("expected UsedAt to be nil for a fresh token")
	}
}

func TestRepository_MarkTokenUsed(t *testing.T) {
	db := openTestDB(t)
	repo := NewRepository(db)

	name := uniqueName("markused")
	user := insertTestUser(t, db, repo, name, name+"@test.com")

	tokenHash := fmt.Sprintf("usehash_%d", time.Now().UnixNano())
	if err := repo.CreatePasswordResetToken(user.ID, tokenHash, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	token, err := repo.FindPasswordResetToken(tokenHash)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.MarkTokenUsed(token.ID); err != nil {
		t.Fatal(err)
	}

	updated, err := repo.FindPasswordResetToken(tokenHash)
	if err != nil {
		t.Fatal(err)
	}
	if updated.UsedAt == nil {
		t.Error("expected UsedAt to be set after MarkTokenUsed")
	}

	if err := repo.MarkTokenUsed(token.ID); err == nil {
		t.Error("expected error when marking an already-used token")
	}
}

func TestRepository_FindPasswordResetToken_NotFound(t *testing.T) {
	db := openTestDB(t)
	repo := NewRepository(db)

	token, err := repo.FindPasswordResetToken("nonexistent_hash_xyz")
	if err != nil {
		t.Fatal(err)
	}
	if token != nil {
		t.Fatal("expected nil for unknown token hash")
	}
}
