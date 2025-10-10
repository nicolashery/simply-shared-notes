package emails

import (
	"context"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicolashery/simply-shared-notes/app/access"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
)

func SpaceCreatedSubject(ctx context.Context, space *db.Space) string {
	intl := rctx.GetIntl(ctx)
	msg := intl.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Emails.SpaceCreated",
			Other: "Your space was created: {{.Space}}",
		},
		TemplateData: map[string]any{
			"Space": space.Name,
		},
	})
	return msg
}

func SpaceCreatedText(memberName string, space *db.Space, baseURL string, tokens access.AccessTokens) string {
	return fmt.Sprintf(
		"Hi %s!\n\n"+
			"You created a new space for your notes: %s.\n\n"+
			"Share this link with anyone you want to collaborate with in this space (editor access):\n"+
			"%s/s/%s\n\n"+
			"You can instead share this link if you don't want to allow edits (view-only access):\n"+
			"%s/s/%s\n\n"+
			"Finally, keep this link for yourself as it will allow you to do anything in this space (admin access):\n"+
			"%s/s/%s\n\n"+
			"Enjoy,\n"+
			"- Simply Shared Notes",
		memberName, space.Name,
		baseURL, tokens.EditToken,
		baseURL, tokens.ViewToken,
		baseURL, tokens.AdminToken,
	)
}
