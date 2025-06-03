package forms

import (
	"net/http"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zhttp"
)

type CreateSpace struct {
	Name     string `zog:"name"`
	Identity string `zog:"identity"`
	Email    string `zog:"email"`
	Code     string `zog:"code"`
}

var spaceNameSchema = z.String().Trim().Required().Max(255)

func createSpaceSchema(requiresCode bool) *z.StructSchema {
	codeSchema := z.String().Trim()
	if requiresCode {
		codeSchema = codeSchema.Required()
	}

	return z.Struct(z.Shape{
		"name":     spaceNameSchema,
		"identity": z.String().Trim().Required().Max(255),
		"email":    z.String().Trim().Required().Email().Max(255),
		"code":     codeSchema,
	})
}

func ParseCreateSpace(r *http.Request, requiresCode bool) (CreateSpace, map[string][]string) {
	var form CreateSpace
	if !requiresCode {
		form.Code = "placeholder"
	}

	errs := createSpaceSchema(requiresCode).Parse(zhttp.Request(r), &form)
	if errs == nil {
		return form, nil
	}

	return form, z.Issues.SanitizeMap(errs)
}

type UpdateSpace struct {
	Name string `zog:"name"`
}

var updateSpaceSchema = z.Struct(z.Shape{
	"name": spaceNameSchema,
})

func ParseUpdateSpace(r *http.Request) (UpdateSpace, map[string][]string) {
	var form UpdateSpace
	errs := updateSpaceSchema.Parse(zhttp.Request(r), &form)
	if errs == nil {
		return form, nil
	}

	return form, z.Issues.SanitizeMap(errs)
}
