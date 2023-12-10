package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/haidongNg/lenslocked/rand"
)

const (
	// The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	SysUserID int
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
}

func (ss *SessionService) Create(sysUserId int) (*Session, error) {
	// Create the session token
	// Implement SessionService.Create
	bytesPerToken := ss.BytesPerToken

	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("Create Token: %w", err)
	}

	// Hash the session token
	session := Session{
		SysUserID: sysUserId,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	// Store the session in our DB
	row := ss.DB.QueryRow(`UPDATE sys_sessions SET token_hash = $1 WHERE sys_user_id = $2 RETURNING id`, session.TokenHash, session.SysUserID)

	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		// If no session exists, we will get ErrNoRows. That means we need to
		// create a session object for that user.
		row = ss.DB.QueryRow(`INSERT INTO sys_sessions (sys_user_id, token_hash) VALUES ($1,$2) RETURNING id`, session.SysUserID, session.TokenHash)
		// The error will be overwritten with either a new error, or nil
		err = row.Scan(&session.ID)
	}
	// If the err was not sql.ErrNoRows, we need to check to see if it was any
	// other error. If it was sql.ErrNoRows it will be overwritten inside the if
	// block, and we still need to check for any errors.
	if err != nil {
		return nil, fmt.Errorf("Create  Session %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// Implement SessionService.User
	// Hash the session token
	tokenHash := ss.hash(token)
	// Query for the session with that hash
	var user User
	row := ss.DB.QueryRow(`SELECT sys_user_id FROM sys_sessions WHERE token_hash = $1`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("Session: %w", err)
	}
	// Using the UserID from the session, we need to query for that user
	row = ss.DB.QueryRow(`SELECT email, password_hash FROM sys_users WHERE id = $1`, user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("User: %w", err)
	}
	// Return the user
	return &user, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))

	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)

	_, err := ss.DB.Exec(`DELETE FROM sys_sessions WHERE token_hash = $1`, tokenHash)
	if err != nil {
		return fmt.Errorf("Detele Session: %w", err)
	}
	return nil
}
