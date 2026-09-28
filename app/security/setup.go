package security

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"math/big"
	"sync"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
	"github.com/YahirHub/axolotl-dashboard-template/app/models"
)

var setupState struct {
	sync.Mutex
	pin string
}

func AdminExists() bool {
	if !facades.Schema().HasTable("users") {
		return false
	}

	exists, err := facades.Orm().Query().
		Model(&models.User{}).
		Where("role", models.RoleAdmin).
		Where("active", true).
		Exists()
	return err == nil && exists
}

func InitializeAdminPIN() (string, error) {
	setupState.Lock()
	defer setupState.Unlock()

	if !facades.Schema().HasTable("users") {
		setupState.pin = ""
		return "", fmt.Errorf("tabla users no disponible; ejecuta migraciones")
	}

	if AdminExists() {
		setupState.pin = ""
		return "", nil
	}
	if setupState.pin != "" {
		return setupState.pin, nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}

	setupState.pin = fmt.Sprintf("%06d", n.Int64())
	logger := facades.Log().Channel("setup")
	logger.Warning("============================================================")
	logger.Warningf("AXOLOTL SETUP PIN: %s", setupState.pin)
	logger.Warning("No existe un administrador. Abre /setup para crear la cuenta inicial.")
	logger.Warning("El PIN es temporal y cambia al reiniciar hasta completar el setup.")
	logger.Warning("============================================================")
	return setupState.pin, nil
}

func CompleteWithPIN(pin string, createAdmin func() error) error {
	setupState.Lock()
	defer setupState.Unlock()

	if setupState.pin == "" {
		return fmt.Errorf("setup no disponible")
	}

	if subtle.ConstantTimeCompare([]byte(pin), []byte(setupState.pin)) != 1 {
		return fmt.Errorf("pin inválido")
	}

	if AdminExists() {
		setupState.pin = ""
		return fmt.Errorf("ya existe un administrador")
	}

	if err := createAdmin(); err != nil {
		return err
	}

	setupState.pin = ""
	return nil
}
