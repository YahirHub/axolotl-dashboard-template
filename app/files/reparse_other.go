//go:build !windows

package files

import "os"

func isReparsePoint(info os.FileInfo) bool {
	return false
}
