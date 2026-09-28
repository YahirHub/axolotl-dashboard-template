package security

import (
	"testing"
	"time"
)

func TestAuthLimiterProgressiveLockAndRecovery(t *testing.T) {
	limiter := NewAuthLimiter()
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	limiter.now = func() time.Time { return now }

	ip := "127.0.0.1"
	email := "user@example.test"

	for i := 0; i < 4; i++ {
		if retry := limiter.RecordFailure(ip, email); retry != 0 {
			t.Fatalf("no debía bloquear en intento %d: %s", i+1, retry)
		}
	}

	retry := limiter.RecordFailure(ip, email)
	if retry < 15*time.Second {
		t.Fatalf("se esperaba bloqueo progresivo, recibido %s", retry)
	}

	if wait, blocked := limiter.Check(ip, email); !blocked || wait <= 0 {
		t.Fatalf("se esperaba login bloqueado")
	}

	now = now.Add(16 * time.Second)
	if _, blocked := limiter.Check(ip, email); blocked {
		t.Fatalf("el primer bloqueo debía expirar")
	}

	limiter.RecordSuccess(ip, email)
	if _, blocked := limiter.Check(ip, email); blocked {
		t.Fatalf("login exitoso debía limpiar el bloqueo de identidad")
	}
}

func TestAuthLimiterGlobalIPProtection(t *testing.T) {
	limiter := NewAuthLimiter()
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	limiter.now = func() time.Time { return now }

	for i := 0; i < 20; i++ {
		limiter.RecordFailure("10.0.0.5", "user"+string(rune('a'+i))+"@example.test")
	}

	if wait, blocked := limiter.Check("10.0.0.5", "another@example.test"); !blocked || wait < time.Minute {
		t.Fatalf("se esperaba bloqueo global de IP, wait=%s blocked=%v", wait, blocked)
	}
}

func TestIdentityFingerprintDoesNotExposeEmail(t *testing.T) {
	fingerprint := IdentityFingerprint("Secret.User@Example.Test")
	if fingerprint == "" || fingerprint == "Secret.User@Example.Test" {
		t.Fatalf("fingerprint inválido")
	}
	if len(fingerprint) != 16 {
		t.Fatalf("longitud inesperada: %d", len(fingerprint))
	}
}
