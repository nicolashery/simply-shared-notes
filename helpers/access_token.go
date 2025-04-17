package helpers

import gonanoid "github.com/matoous/go-nanoid/v2"

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
