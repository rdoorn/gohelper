package passwordhash

import "fmt"

var (
	DefaultPasswordSaltSize int = 32
)

type HasherInterface interface {
	Salt() string
	Hash(...string) string
}

func New(algorithm string) (HasherInterface, error) {
	switch algorithm {
	case "sha256":
		return NewSha256(), nil
	default:
		return nil, fmt.Errorf("Unknown hash algorithm: %s", algorithm)
	}
}
