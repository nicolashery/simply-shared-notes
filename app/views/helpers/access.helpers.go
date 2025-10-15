package helpers

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicolashery/simply-shared-notes/app/access"
	"github.com/nicolashery/simply-shared-notes/app/intl"
)

func RoleLabel(intl *intl.Intl, access *access.Access) string {
	if access.IsAdmin() {
		return intl.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Helpers.Role.Admin",
				Other: "Admin",
			},
		})
	}

	if access.IsEdit() {
		return intl.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Helpers.Role.Editor",
				Other: "Editor",
			},
		})
	}

	if access.IsView() {
		return intl.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Helpers.Role.ViewOnly",
				Other: "View-only",
			},
		})
	}

	return "<unknown role>"
}
