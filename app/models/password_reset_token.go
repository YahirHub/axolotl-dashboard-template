package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type PasswordResetToken struct {
	orm.Model
	UserID    uint
	TokenHash string
	ExpiresAt time.Time
	Used      bool
}
