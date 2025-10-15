package helpers

import "github.com/nicolashery/simply-shared-notes/app/access"

func RoleLabel(access *access.Access) string {
	if access.IsAdmin() {
		return "Admin"
	}

	if access.IsEdit() {
		return "Editor"
	}

	if access.IsView() {
		return "View-only"
	}

	return "<unknown role>"
}
