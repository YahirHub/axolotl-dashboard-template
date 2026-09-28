package files

import "github.com/YahirHub/axolotl-dashboard-template/app/facades"

func NewConfiguredWorkspace() (*Workspace, error) {
	return NewWorkspace(
		facades.Config().GetString("files.root"),
		int64(facades.Config().GetInt("files.max_edit_bytes", 2*1024*1024)),
	)
}
