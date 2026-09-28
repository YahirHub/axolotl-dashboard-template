package tests

import (
	"os"
	"path/filepath"

	"github.com/goravel/framework/testing"

	"github.com/YahirHub/axolotl-dashboard-template/bootstrap"
)

func init() {
	testDatabase := filepath.Join(os.TempDir(), "axolotl-dashboard-template-test.sqlite")
	_ = os.Remove(testDatabase)

	_ = os.Setenv("APP_ENV", "testing")
	_ = os.Setenv("DB_CONNECTION", "sqlite")
	_ = os.Setenv("DB_DATABASE", testDatabase)

	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
