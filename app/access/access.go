package access

import (
	"fmt"
	"regexp"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nicolashery/simply-shared-notes/app/db"
)

const (
	AccessTokenAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	AccessTokenLength   = 20
)

func GenerateAccessToken() (string, error) {
	return gonanoid.Generate(AccessTokenAlphabet, AccessTokenLength)
}

type AccessTokens struct {
	AdminToken string
	EditToken  string
	ViewToken  string
}

func GenerateAccessTokens() (AccessTokens, error) {
	adminToken, err := GenerateAccessToken()
	if err != nil {
		return AccessTokens{}, err
	}

	editToken, err := GenerateAccessToken()
	if err != nil {
		return AccessTokens{}, err
	}

	viewToken, err := GenerateAccessToken()
	if err != nil {
		return AccessTokens{}, err
	}

	return AccessTokens{
		AdminToken: adminToken,
		EditToken:  editToken,
		ViewToken:  viewToken,
	}, nil
}

var accessTokenPattern = regexp.MustCompile(fmt.Sprintf("^[A-Za-z0-9]{%d}$",
	AccessTokenLength))

func IsValidAccessToken(token string) bool {
	return accessTokenPattern.MatchString(token)
}

type Role string

const (
	Role_Admin = "admin"
	Role_Edit  = "edit"
	Role_View  = "view"
)

func GetTokenRole(space *db.Space, token string) (Role, bool) {
	switch token {
	case space.AdminToken:
		return Role_Admin, true
	case space.EditToken:
		return Role_Edit, true
	case space.ViewToken:
		return Role_View, true
	}

	return "", false
}

type Access struct {
	Token string
	Role  Role
}
