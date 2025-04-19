package helpers

import "github.com/nicolashery/simply-shared-notes/app/access"

func RoleLabel(role access.Role) string {
	switch role {
	case access.Role_Admin:
		return "Admin"
	case access.Role_Edit:
		return "Editor"
	case access.Role_View:
		return "Viewer"
	default:
		return ""
	}
}
