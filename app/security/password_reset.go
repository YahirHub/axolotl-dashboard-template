package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	contractsorm "github.com/goravel/framework/contracts/database/orm"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
)

var ErrResetTokenInvalid = errors.New("token de recuperación inválido o expirado")

type PasswordResetManager struct {
	now func() time.Time
}

func NewPasswordResetManager() *PasswordResetManager {
	return &PasswordResetManager{now: time.Now}
}

func (m *PasswordResetManager) Issue(userID uint, ttl time.Duration) (string, error) {
	if ttl < 5*time.Minute {
		ttl = 5 * time.Minute
	}
	if ttl > 2*time.Hour {
		ttl = 2 * time.Hour
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)

	if _, err := facades.Orm().Query().
		Where("user_id", userID).
		Delete(&models.PasswordResetToken{}); err != nil {
		return "", err
	}

	record := models.PasswordResetToken{
		UserID:    userID,
		TokenHash: resetTokenHash(token),
		ExpiresAt: m.now().UTC().Add(ttl),
		Used:      false,
	}
	if err := facades.Orm().Query().Create(&record); err != nil {
		return "", err
	}

	return token, nil
}

func (m *PasswordResetManager) Validate(token string) (*models.PasswordResetToken, error) {
	if !validRawResetToken(token) {
		return nil, ErrResetTokenInvalid
	}

	var record models.PasswordResetToken
	if err := facades.Orm().Query().
		Where("token_hash", resetTokenHash(token)).
		Where("used", false).
		First(&record); err != nil {
		return nil, ErrResetTokenInvalid
	}

	if !record.ExpiresAt.After(m.now().UTC()) {
		return nil, ErrResetTokenInvalid
	}

	return &record, nil
}

func (m *PasswordResetManager) Consume(token, newPasswordHash string) (*models.User, error) {
	if !validRawResetToken(token) || newPasswordHash == "" {
		return nil, ErrResetTokenInvalid
	}

	var updated models.User
	err := facades.Orm().Transaction(func(tx contractsorm.Query) error {
		var record models.PasswordResetToken
		if err := tx.
			Where("token_hash", resetTokenHash(token)).
			Where("used", false).
			First(&record); err != nil {
			return ErrResetTokenInvalid
		}
		if !record.ExpiresAt.After(m.now().UTC()) {
			return ErrResetTokenInvalid
		}

		var user models.User
		if err := tx.Where("id", record.UserID).Where("active", true).First(&user); err != nil {
			return ErrResetTokenInvalid
		}

		user.Password = newPasswordHash
		if user.AuthVersion == 0 {
			user.AuthVersion = 1
		}
		user.AuthVersion++
		if err := tx.Save(&user); err != nil {
			return err
		}
		updated = user

		record.Used = true
		if err := tx.Save(&record); err != nil {
			return err
		}

		_, err := tx.
			Where("user_id", user.ID).
			Where("id <> ?", record.ID).
			Delete(&models.PasswordResetToken{})
		return err
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func resetTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func validRawResetToken(token string) bool {
	if len(token) < 40 || len(token) > 128 {
		return false
	}
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	return err == nil && len(decoded) == 32
}
