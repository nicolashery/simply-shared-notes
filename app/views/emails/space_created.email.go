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

func SpaceCreatedText(ctx context.Context, memberName string, space *db.Space, baseURL string, tokens access.AccessTokens) string {
	intl := rctx.GetIntl(ctx)

	greeting := intl.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Emails.SpaceCreated.Greeting",
			Other: "Hi {{.Name}}!",
		},
		TemplateData: map[string]any{
			"Name": memberName,
		},
	})

	introduction := intl.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Emails.SpaceCreated.Introduction",
			Other: "You created a new space for your notes: {{.Space}}.",
		},
		TemplateData: map[string]any{
			"Space": space.Name,
		},
	})

	editorLink := intl.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Emails.SpaceCreated.EditorLink",
			Other: "Share this link with anyone you want to collaborate with in this space (editor access):\n{{.URL}}",
		},
		TemplateData: map[string]any{
			"URL": fmt.Sprintf("%s/s/%s", baseURL, tokens.EditToken),
		},
	})

	viewLink := intl.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Emails.SpaceCreated.ViewLink",
			Other: "You can instead share this link if you don't want to allow edits (view-only access):\n{{.URL}}",
		},
		TemplateData: map[string]any{
			"URL": fmt.Sprintf("%s/s/%s", baseURL, tokens.ViewToken),
		},
	})

	adminLink := intl.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Emails.SpaceCreated.AdminLink",
			Other: "Finally, keep this link for yourself as it will allow you to do anything in this space (admin access):\n{{.URL}}",
		},
		TemplateData: map[string]any{
			"URL": fmt.Sprintf("%s/s/%s", baseURL, tokens.AdminToken),
		},
	})

	signature := intl.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Emails.SpaceCreated.Signature",
			Other: "Enjoy,\n- Simply Shared Notes",
		},
	})

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
		greeting,
		introduction,
		editorLink,
		viewLink,
		adminLink,
		signature,
	)
}
