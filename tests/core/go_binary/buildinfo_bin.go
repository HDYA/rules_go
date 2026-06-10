package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	info, ok := debug.ReadBuildInfo()
	fmt.Printf("ok=%t\n", ok)
	if !ok {
		return
	}
	fmt.Printf("path=%s\n", info.Path)
	for _, setting := range info.Settings {
		switch setting.Key {
		case "-buildmode", "-compiler", "CGO_ENABLED", "GOARCH", "GOOS":
			fmt.Printf("setting %s=%s\n", setting.Key, setting.Value)
		}
	}
}
