package forms

import (
	"net/http"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zhttp"
)

var memberNameSchema = z.String().Trim().Required().Min(1).Max(255)

type CreateMember struct {
	Name string `zog:"name"`
}

var createMemberSchema = z.Struct(z.Shape{
	"name": memberNameSchema,
})

func ParseCreateMember(r *http.Request) (CreateMember, map[string][]string) {
	var form CreateMember
	errs := createMemberSchema.Parse(zhttp.Request(r), &form)
	if errs == nil {
		return form, nil
	}

	return form, z.Issues.SanitizeMap(errs)
}

type UpdateMember struct {
	Name string `zog:"name"`
}

var updateMemberSchema = z.Struct(z.Shape{
	"name": memberNameSchema,
})

func ParseUpdateMember(r *http.Request) (UpdateMember, map[string][]string) {
	var form UpdateMember
	errs := updateMemberSchema.Parse(zhttp.Request(r), &form)
	if errs == nil {
		return form, nil
	}

	return form, z.Issues.SanitizeMap(errs)
}
