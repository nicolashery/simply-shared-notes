package helpers

import (
	"golang.org/x/text/language"
)

func LanguageLabel(lang language.Tag) string {
	if lang == language.English {
		return "English (EN)"
	}

	if lang == language.French {
		return "Français (FR)"
	}

	return "<unknown language>"
}
