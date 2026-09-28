package config

import (
	"fmt"
	"path/filepath"

	"github.com/goravel/framework/support/path"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
)

func init() {
	config := facades.Config()

	root := fmt.Sprint(config.Env("FILES_ROOT", "storage/demo-files"))
	if !filepath.IsAbs(root) {
		root = path.Base(root)
	}

	config.Add("files", map[string]any{
		"root":           root,
		"max_edit_bytes": config.Env("FILES_MAX_EDIT_BYTES", 2*1024*1024),
	})
}
