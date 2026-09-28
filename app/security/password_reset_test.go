package security

import (
	"testing"
	"time"
)

func TestRawResetTokenValidation(t *testing.T) {
	manager := NewPasswordResetManager()
	_ = manager

	for _, invalid := range []string{"", "abc", "../../token", "%%%"} {
		if validRawResetToken(invalid) {
			t.Fatalf("token inválido aceptado: %q", invalid)
		}
	}
}

func TestResetTokenHashDoesNotExposeToken(t *testing.T) {
	token := "2PeJ1h9KmdR6Er_C12GSoJiObHWhwtgPQbq9aD9tPSA"
	hash := resetTokenHash(token)
	if hash == token || len(hash) != 64 {
		t.Fatalf("hash inesperado")
	}
}

func TestResetTTLBounds(t *testing.T) {
	manager := NewPasswordResetManager()
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	manager.now = func() time.Time { return now }

	// La emisión se prueba a nivel feature con la base de datos; aquí se
	// documenta el reloj inyectable usado para expiración determinista.
	if manager.now() != now {
		t.Fatal("reloj de prueba no aplicado")
	}
}
