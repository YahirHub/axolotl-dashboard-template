package tests

import (
	"github.com/goravel/framework/testing"

	"github.com/YahirHub/axolotl-dashboard-template/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
