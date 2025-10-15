package helpers

import (
	"database/sql"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/intl"
)

func DisplayMemberName(intl *intl.Intl, id sql.NullInt64, membersByID map[int64]db.Member) string {
	if !id.Valid {
		return intl.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Helpers.Member.Deleted",
				Other: "Deleted Member",
			},
		})
	}

	member, ok := membersByID[id.Int64]
	if !ok {
		return intl.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Helpers.Member.Unknown",
				Other: "Unknown Member",
			},
		})
	}

	return member.Name
}
