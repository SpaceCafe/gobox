package utils

import (
	"runtime"
)

// SysConfigDir returns the default system configuration directory path based on the operating system.
func SysConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\ProgramData"
	case "darwin", "ios":
		return "/Library/Application Support"
	default:
		return "/etc"
	}
}
