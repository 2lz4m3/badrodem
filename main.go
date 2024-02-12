package main

import (
	"badrodem/cmd"
	"fmt"
	"runtime/debug"
)

var (
	version = "dev"
	commit  = "unset"
	date    = "unset"
	builtBy = "unset"
)

func main() {
	versionString := buildVersionString(version, commit, date, builtBy)
	cmd.Execute(versionString)
}

func buildVersionString(version, commit, date, builtBy string) string {
	result := version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}
	return result
}
