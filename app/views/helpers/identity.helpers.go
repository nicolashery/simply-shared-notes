package helpers

import "github.com/nicolashery/simply-shared-notes/app/identity"

func IdentityName(identity *identity.Identity) string {
	if identity.Anonymous {
		return "Anonymous"
	}

	return identity.Member.Name
}
