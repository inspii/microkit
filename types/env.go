package types

import (
	"os"
)

func Env(name string, defaultValue string) StrValue {
	val, ok := os.LookupEnv(name)
	if !ok || val == "" {
		return StrValue(defaultValue)
	}

	return StrValue(val)
}
