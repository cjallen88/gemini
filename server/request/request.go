// Package request provides functionality to parse and validate Gemini protocol requests.
// Specification: https://geminiprotocol.net/docs/protocol-specification.gmi

package request

import (
	"errors"
	"net/url"
	"strings"
)

type Request struct {
	Url *url.URL
}

func ParseUrl(in string) (*url.URL, error) {
	url, err := url.Parse(in)
	if err != nil {
		return nil, err
	}

	fragmentIndex := strings.Index(in, "#")
	if fragmentIndex != -1 {
		return nil, errors.New("request contains a fragment")
	}

	if url.User != nil {
		return nil, errors.New("request contains user info")
	}

	if url.Scheme != "gemini" {
		return nil, errors.New("invalid scheme, expected 'gemini'")
	}

	if url.Path == "" {
		url.Path = "/"
	}

	return url, err
}

func ParseRequest(requestStr string) (Request, error) {
	url, err := ParseUrl(requestStr)
	if err != nil {
		return Request{}, err
	}

	return Request{Url: url}, nil
}
