package publicid

import (
	"fmt"
	"regexp"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	PublicIDAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	PublicIDLength   = 10
)

func Generate() (string, error) {
	return gonanoid.Generate(PublicIDAlphabet, PublicIDLength)
}

var publicIDPattern = regexp.MustCompile(fmt.Sprintf("^[A-Za-z0-9]{%d}$",
	PublicIDLength))

func IsValidPublicID(token string) bool {
	return publicIDPattern.MatchString(token)
}
