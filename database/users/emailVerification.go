package users

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

const (
	MaxEmailCodeAttempts = 3
	EmailCodeTTL         = 10 * time.Minute
	EmailVerifiedTTL     = 30 * time.Minute
	EmailCodeCooldown    = 60 * time.Second
)

var (
	ErrEmailTaken       = errors.New("email is already registered")
	ErrCodeCooldown     = errors.New("verification code was sent recently")
	ErrCodeNotFound     = errors.New("no pending verification code")
	ErrCodeExpired      = errors.New("verification code has expired")
	ErrTooManyAttempts  = errors.New("too many incorrect attempts")
	ErrCodeMismatch     = errors.New("incorrect verification code")
	ErrEmailNotVerified = errors.New("email is not verified")
)
/*
SaveEmailCode saves or updates the email verification code in the database.
It also enforces the code request cooldown and returns the remaining wait time
if the user requests a new code too soon.

Parameters:
	db *sql.DB, email, code string, now time.Time 

Returns:
	time.Duration
		-> cooldown
	error
		-> nil if success
*/
func SaveEmailCode(db *sql.DB, email, code string, now time.Time) (time.Duration, error) {
	taken, err := CheckUserEmail(db, email)
	if err != nil {
		return 0, err
	}
	
	if taken {
		return 0, ErrEmailTaken
	}

	if _, err := db.Exec(
		`DELETE FROM email_verifications WHERE expires_at < ?`,
		now.Add(-time.Hour).Unix(),
	); err != nil {
		return 0, err
	}

	result, err := db.Exec(`
		INSERT INTO email_verifications (email, code, attempts, verified, token, sent_at, expires_at)
		VALUES (?, ?, 0, 0, '', ?, ?)
		ON CONFLICT(email) DO UPDATE SET
			code = excluded.code,
			attempts = 0,
			verified = 0,
			token = '',
			sent_at = excluded.sent_at,
			expires_at = excluded.expires_at
		WHERE email_verifications.sent_at <= ?
	`,
		email,
		code,
		now.Unix(),
		now.Add(EmailCodeTTL).Unix(),
		now.Add(-EmailCodeCooldown).Unix(),
	)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected > 0 {
		return 0, nil
	}

	var sentAt int64
	if err := db.QueryRow(
		`SELECT sent_at FROM email_verifications WHERE email = ?`,
		email,
	).Scan(&sentAt); err != nil {
		return 0, err
	}

	wait := time.Unix(sentAt, 0).Add(EmailCodeCooldown).Sub(now)
	if wait < time.Second {
		wait = time.Second
	}

	return wait, ErrCodeCooldown
}

func DeleteEmailCode(db *sql.DB, email string) error {
	_, err := db.Exec(`DELETE FROM email_verifications WHERE email = ?`, email)
	return err
}

func VerifyEmailCode(db *sql.DB, email, code string, now time.Time) (string, int, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE email_verifications
		SET attempts = attempts + 1
		WHERE email = ? AND verified = 0 AND attempts < ? AND expires_at > ?
	`, email, MaxEmailCodeAttempts, now.Unix())
	if err != nil {
		return "", 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", 0, err
	}

	if affected == 0 {
		var attempts, verified int
		var expiresAt int64

		err := tx.QueryRow(`
			SELECT attempts, verified, expires_at
			FROM email_verifications
			WHERE email = ?
		`, email).Scan(&attempts, &verified, &expiresAt)
		if errors.Is(err, sql.ErrNoRows) {
			return "", 0, ErrCodeNotFound
		}
		if err != nil {
			return "", 0, err
		}
		if verified == 1 {
			return "", 0, ErrCodeNotFound
		}
		if expiresAt <= now.Unix() {
			return "", 0, ErrCodeExpired
		}

		return "", 0, ErrTooManyAttempts
	}

	var stored string
	var attempts int

	if err := tx.QueryRow(
		`SELECT code, attempts FROM email_verifications WHERE email = ?`,
		email,
	).Scan(&stored, &attempts); err != nil {
		return "", 0, err
	}

	attemptsLeft := MaxEmailCodeAttempts - attempts

	if subtle.ConstantTimeCompare([]byte(stored), []byte(code)) != 1 {
		if err := tx.Commit(); err != nil {
			return "", 0, err
		}
		return "", attemptsLeft, ErrCodeMismatch
	}

	token, err := newVerifyToken()
	if err != nil {
		return "", 0, err
	}

	if _, err := tx.Exec(`
		UPDATE email_verifications
		SET verified = 1, token = ?, expires_at = ?
		WHERE email = ?
	`, token, now.Add(EmailVerifiedTTL).Unix(), email); err != nil {
		return "", 0, err
	}

	if err := tx.Commit(); err != nil {
		return "", 0, err
	}

	return token, attemptsLeft, nil
}

func HasVerifiedEmail(db *sql.DB, email, token string, now time.Time) (bool, error) {
	if token == "" {
		return false, nil
	}

	var exists int

	err := db.QueryRow(`
		SELECT 1
		FROM email_verifications
		WHERE email = ? AND token = ? AND verified = 1 AND expires_at > ?
		LIMIT 1
	`, email, token, now.Unix()).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func newVerifyToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
