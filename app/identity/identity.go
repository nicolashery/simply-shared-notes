package identity

import "github.com/nicolashery/simply-shared-notes/app/db"

type Identity struct {
	Member    *db.Member
	Anonymous bool
}

func Member(member *db.Member) Identity {
	return Identity{
		Member: member,
	}
}

func Anonymous() Identity {
	return Identity{
		Anonymous: true,
	}
}
