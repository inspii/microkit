package types

import "strings"

type Secret string

func (s Secret) MarshalJSON() ([]byte, error) {
	secret := strings.Repeat("*", len(s))
	return []byte(`"` + secret + `"`), nil
}
