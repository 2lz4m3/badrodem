package localize

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	I18n MyI18n
)

type MyI18n struct {
	Tag       language.Tag
	Bundle    *i18n.Bundle
	Localizer *i18n.Localizer
}

func T(messageId string) string {
	config := &i18n.LocalizeConfig{MessageID: messageId}
	message := I18n.Localizer.MustLocalize(config)
	return message
}
