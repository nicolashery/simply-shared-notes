package intl

import (
	"embed"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/goodsign/monday"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var DefaultLang language.Tag = language.English

var SupportedLangs = []language.Tag{
	language.English,
	language.French,
}

const LocalizeErrorMessage string = "<failed to localize>"

type Intl struct {
	CurrentLang language.Tag
	localizer   *i18n.Localizer
	logger      *slog.Logger
}

func NewBundle(localesFS embed.FS) (*i18n.Bundle, error) {
	b := i18n.NewBundle(DefaultLang)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, tag := range SupportedLangs {
		path := fmt.Sprintf("locales/active.%s.toml", tag.String())
		// debug: check if path exists in localFS
		// if _, err := localeFS.Open(path); err != nil {
		// 	return nil, fmt.Errorf("loading i18n file %q: %w", path, err)
		// }
		if _, err := b.LoadMessageFileFS(localesFS, path); err != nil {
			return nil, fmt.Errorf("loading i18n file %q: %w", path, err)
		}
	}

	return b, nil
}

func New(logger *slog.Logger, i18nBundle *i18n.Bundle, lang language.Tag) *Intl {
	localizer := i18n.NewLocalizer(i18nBundle, lang.String())

	return &Intl{
		CurrentLang: lang,
		localizer:   localizer,
		logger:      logger,
	}
}

func (i *Intl) Localize(lc *i18n.LocalizeConfig) string {
	msg, err := i.localizer.Localize(lc)
	if err != nil {
		i.logger.Error(
			"failed to localize",
			slog.Any("error", err),
			slog.String("lang", i.CurrentLang.String()),
			slog.String("key", lc.MessageID),
		)
		return LocalizeErrorMessage
	}

	return msg
}

func (i *Intl) SplitOnSlot(s, slot string) (string, string) {
	j := strings.Index(s, slot)
	if j < 0 {
		return s, ""
	}
	return s[:j], s[j+len(slot):]
}

func mondayLocale(tag language.Tag) monday.Locale {
	base, _ := tag.Base()

	switch base.String() {
	case "fr":
		return monday.LocaleFrFR
	default:
		return monday.LocaleEnUS
	}
}

func (i *Intl) FormatDate(t time.Time) string {
	locale := mondayLocale(i.CurrentLang)
	format, ok := monday.MediumFormatsByLocale[locale]
	if !ok {
		format = "Jan 02, 2006"
	}
	return monday.Format(t, format, locale)
}

func (i *Intl) FormatTime(t time.Time) string {
	switch i.CurrentLang {
	case language.French:
		return t.Format("15:04 MST")
	default:
		return t.Format("3:04 PM MST")
	}
}
