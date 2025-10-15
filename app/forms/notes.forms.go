package forms

import (
	"net/http"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zhttp"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
)

var noteTitleSchema = z.String().Trim().Required().Min(1).Max(255)

type CreateNote struct {
	Title   string `zog:"title"`
	Content string `zog:"content"`
}

var createNoteSchema = z.Struct(z.Shape{
	"title":   noteTitleSchema,
	"content": z.String().Trim().Required(),
})

func ParseCreateNote(r *http.Request) (CreateNote, map[string][]string) {
	var form CreateNote
	intl := rctx.GetIntl(r.Context())
	errs := createNoteSchema.Parse(zhttp.Request(r), &form, intl.ZogParseOpts())
	if errs == nil {
		return form, nil
	}

	return form, z.Issues.SanitizeMap(errs)
}

type UpdateNote struct {
	Title   string `zog:"title"`
	Content string `zog:"content"`
}

var updateNoteSchema = z.Struct(z.Shape{
	"title":   noteTitleSchema,
	"content": z.String().Trim(),
})

func ParseUpdateNote(r *http.Request) (UpdateNote, map[string][]string) {
	var form UpdateNote
	intl := rctx.GetIntl(r.Context())
	errs := updateNoteSchema.Parse(zhttp.Request(r), &form, intl.ZogParseOpts())
	if errs == nil {
		return form, nil
	}

	return form, z.Issues.SanitizeMap(errs)
}
