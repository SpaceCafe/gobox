package config

import (
	"reflect"
	"runtime"
)

// IsStructPointer checks if the provided input is a non-nil pointer to a struct and returns a boolean value.
func IsStructPointer(in any) bool {
	value := reflect.ValueOf(in)
	return in != nil && value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct
}

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
