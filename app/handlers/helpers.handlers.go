package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/nicolashery/simply-shared-notes/app/db"
)

func safeRedirect(s string) string {
	if s == "" {
		return ""
	}

	// accept only absolute-path, not protocol-relative
	if !strings.HasPrefix(s, "/") || strings.HasPrefix(s, "//") {
		return ""
	}

	return s
}

func redirectFromReferer(r *http.Request) string {
	ref := r.Header.Get("Referer")
	if ref == "" {
		return ""
	}

	u, err := url.Parse(ref)
	if err != nil {
		return ""
	}

	// same-origin check
	if u.Host != r.Host {
		return ""
	}

	// avoid redirecting back to this handler or its subpaths
	if strings.HasPrefix(u.Path, r.URL.Path) {
		return ""
	}

	return safeRedirect(u.RequestURI())
}

func memberListToMap(members []db.Member) map[int64]db.Member {
	memberMap := make(map[int64]db.Member, len(members))
	for _, member := range members {
		memberMap[member.ID] = member
	}
	return memberMap
}

func collectCreatedUpdatedByIDsFromNotes(notes []db.Note) []int64 {
	idSet := make(map[int64]bool)

	for _, note := range notes {
		if note.CreatedBy.Valid {
			idSet[note.CreatedBy.Int64] = true
		}
		if note.UpdatedBy.Valid {
			idSet[note.UpdatedBy.Int64] = true
		}
	}

	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	return ids
}

func collectCreatedUpdatedByIDsFromMembers(members []db.Member) []int64 {
	idSet := make(map[int64]bool)

	for _, member := range members {
		if member.CreatedBy.Valid {
			idSet[member.CreatedBy.Int64] = true
		}
		if member.UpdatedBy.Valid {
			idSet[member.UpdatedBy.Int64] = true
		}
	}

	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	return ids
}

func collectCreatedUpdatedByIDsFromSpace(space *db.Space) []int64 {
	idSet := make(map[int64]bool)

	if space.CreatedBy.Valid {
		idSet[space.CreatedBy.Int64] = true
	}
	if space.UpdatedBy.Valid {
		idSet[space.UpdatedBy.Int64] = true
	}

	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	return ids
}
