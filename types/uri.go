package types

import (
	"fmt"
	"strings"
)

type URI string

func (u URI) String() string {
	return string(u)
}

func (u URI) WithQuery(key string, value string) URI {
	kv := key + "=" + value
	if strings.Contains(string(u), "?") {
		return URI(string(u) + "&" + kv)
	} else {
		return URI(string(u) + "?" + kv)
	}
}

func (u URI) URL() URL {
	return URL{
		URI: u,
	}
}

type URL struct {
	Schema   string
	Endpoint string
	URI      URI
}

func (u URL) String() string {
	if u.Schema == "" {
		u.Schema = "http"
	}
	if u.Endpoint == "" {
		u.Endpoint = "localhost"
	}

	return fmt.Sprintf("%s://%s%s", u.Schema, u.Endpoint, u.URI.String())
}

func (u URL) WithSchema(schema string) URL {
	clone := u
	clone.Schema = schema
	return clone
}

func (u URL) WithEndpoint(endpoint string) URL {
	clone := u
	clone.Endpoint = endpoint
	return clone
}

func (u URL) WithQuery(key string, value string) URL {
	clone := u
	clone.URI = clone.URI.WithQuery(key, value)
	return clone
}
