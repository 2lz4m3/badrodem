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
	languageBaseString := language.English.String()
	systemLanguageTag, _ := locale.Detect()
	systemLanguageBase, _ := systemLanguageTag.Base()
	japaneseLanguageBase, _ := language.Japanese.Base()
	languageTag := language.English
	if systemLanguageBase.String() == japaneseLanguageBase.String() {
		languageTag = language.Japanese
		languageBaseString = japaneseLanguageBase.String()
	}
	localize.I18n.Tag = systemLanguageTag

	// set up localizer
	localize.I18n.Bundle = i18n.NewBundle(languageTag)
	localize.I18n.Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err := localize.I18n.Bundle.LoadMessageFileFS(embedFS, fmt.Sprintf("locale.%s.json", languageBaseString))
	if err != nil {
		log.Panic(err)
	}
	localize.I18n.Localizer = i18n.NewLocalizer(localize.I18n.Bundle, languageBaseString)

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
