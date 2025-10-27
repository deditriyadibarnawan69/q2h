package main

import (
	"context"
	"net/http"
)

// Config holds the plugin configuration.
type Config struct {
	// QueryParameter is the name of the query param to read the token from, e.g. "token"
	QueryParameter string `yaml:"queryParameter,omitempty"`
	// Header is the header name to set, e.g. "Authorization"
	Header string `yaml:"header,omitempty"`
	// Prefix will be prepended to the token value if non-empty, e.g. "Bearer "
	Prefix string `yaml:"prefix,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		QueryParameter: "token",
		Header:         "Authorization",
		Prefix:         "Bearer ",
	}
}

type queryToHeader struct {
	next   http.Handler
	config *Config
	name   string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config == nil {
		config = CreateConfig()
	}
	return &queryToHeader{
		next:   next,
		config: config,
		name:   name,
	}, nil
}

func (m *queryToHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if m.config != nil && m.config.QueryParameter != "" && m.config.Header != "" {
		if token := req.URL.Query().Get(m.config.QueryParameter); token != "" {
			req.Header.Set(m.config.Header, m.config.Prefix+token)
			// Optionally: remove the query param from the URL before forwarding
			// q := req.URL.Query()
			// q.Del(m.config.QueryParameter)
			// req.URL.RawQuery = q.Encode()
		}
	}
	m.next.ServeHTTP(rw, req)
}
