package tests

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goravel/framework/testing"

	"github.com/YahirHub/axolotl-dashboard-template/bootstrap"
)

func init() {
	testDatabase := filepath.Join(os.TempDir(), fmt.Sprintf("axolotl-dashboard-template-test-%d.sqlite", os.Getpid()))
	testFiles := filepath.Join(os.TempDir(), fmt.Sprintf("axolotl-dashboard-template-files-%d", os.Getpid()))
	_ = os.Remove(testDatabase)
	_ = os.RemoveAll(testFiles)

	_ = os.Setenv("APP_ENV", "testing")
	_ = os.Setenv("DB_CONNECTION", "sqlite")
	_ = os.Setenv("DB_DATABASE", testDatabase)
	_ = os.Setenv("FILES_ROOT", testFiles)

	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
