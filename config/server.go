package config

import (
	"fmt"
	"path/filepath"

	"github.com/goravel/framework/support/path"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
)

func init() {
	config := facades.Config()

	root := fmt.Sprint(config.Env("SERVER_TASKS_ROOT", "storage/server-tasks"))
	if !filepath.IsAbs(root) {
		root = path.Base(root)
	}

	config.Add("server", map[string]any{
		"tasks_root": root,
	})
}
