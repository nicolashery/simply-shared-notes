package helpers

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicolashery/simply-shared-notes/app/identity"
	"github.com/nicolashery/simply-shared-notes/app/intl"
)

func IdentityName(intl *intl.Intl, identity *identity.Identity) string {
	if identity.Anonymous {
		return intl.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Helpers.Identity.Anonymous",
				Other: "Anonymous",
			},
		})
	}

	return identity.Member.Name
}
