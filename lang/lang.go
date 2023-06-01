package lang

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var i18nBundle *i18n.Bundle
var localizer *i18n.Localizer

func InitTranslations() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("lang/zh-CN.json")
	bundle.MustLoadMessageFile("lang/en-US.json")
	i18nBundle = bundle
	return bundle
}

func SetLanguage(lang string) {
	localizer = i18n.NewLocalizer(i18nBundle, lang)
}

func FT(id string, args map[string]interface{}) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
		TemplateData: args,
	})
}

func T(id string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	})
}
