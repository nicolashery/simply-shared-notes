package helpers

import (
	"database/sql"
)

func DisplayMemberName(name sql.NullString) string {
	if name.Valid {
		return name.String
	} else {
		return "Deleted Member"
	}
}
