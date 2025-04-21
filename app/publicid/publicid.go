package publicid

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	PublicIdAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	PublicIdLength   = 10
)

func Generate() (string, error) {
	return gonanoid.Generate(PublicIdAlphabet, PublicIdLength)
}
