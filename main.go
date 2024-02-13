package main

import (
	"badrodem/cmd"
	"badrodem/localize"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/Xuanwo/go-locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	version = "dev"
	commit  = "unset"
	date    = "unset"
	builtBy = "unset"

	//go:embed locale.*.json
	embedFS embed.FS
)

func main() {
	// detect language
	lang := language.English.String()
	tag, _ := locale.Detect()
	if tag == language.Japanese {
		lang = language.Japanese.String()
	} else {
		tag = language.English
	}

	// set up localizer
	localize.I18n.Bundle = i18n.NewBundle(tag)
	localize.I18n.Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err := localize.I18n.Bundle.LoadMessageFileFS(embedFS, fmt.Sprintf("locale.%s.json", lang))
	if err != nil {
		log.Panic(err)
	}
	localize.I18n.Localizer = i18n.NewLocalizer(localize.I18n.Bundle, lang)

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
