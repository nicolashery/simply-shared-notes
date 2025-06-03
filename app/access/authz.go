package access

type Action string

const (
	Action_UpdateSpace  Action = "update_space"
	Action_DeleteSpace  Action = "delete_space"
	Action_ViewTokens   Action = "view_tokens"
	Action_UpdateTokens Action = "updated_tokens"
	Action_CreateMember Action = "create_member"
	Action_UpdateMember Action = "update_member"
	Action_DeleteMember Action = "delete_member"
	Action_CreateNote   Action = "create_note"
	Action_UpdateNote   Action = "update_note"
	Action_DeleteNote   Action = "delete_note"
)

var rolePermissions = map[Role]map[Action]bool{
	Role_Admin: {
		Action_UpdateSpace:  true,
		Action_DeleteSpace:  true,
		Action_ViewTokens:   true,
		Action_UpdateTokens: true,
		Action_CreateMember: true,
		Action_UpdateMember: true,
		Action_DeleteMember: true,
		Action_CreateNote:   true,
		Action_UpdateNote:   true,
		Action_DeleteNote:   true,
	},
	Role_Edit: {
		Action_CreateNote:   true,
		Action_UpdateNote:   true,
		Action_DeleteNote:   true,
		Action_UpdateMember: true,
	},
	Role_View: {
		// This role can only view resources
	},
}

func (a *Access) Can(action Action) bool {
	permissions, exists := rolePermissions[a.Role]
	if !exists {
		return false
	}

	return permissions[action]
}
