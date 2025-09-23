package helpers

import (
	"database/sql"

	"github.com/nicolashery/simply-shared-notes/app/db"
)

func DisplayMemberName(id sql.NullInt64, membersByID map[int64]db.Member) string {
	if !id.Valid {
		return "Deleted Member"
	}

	member, ok := membersByID[id.Int64]
	if !ok {
		return "Unknown Member"
	}

	return member.Name
}
