package response

import (
	"fmt"
	"io"
	"net/url"
)

type RedirectStatus Status

const (
	RedirectTemporary RedirectStatus = 31
	RedirectPermanent RedirectStatus = 32
)

type RedirectResponse struct {
	Status RedirectStatus
	URL    url.URL
}

func (r *RedirectResponse) WriteTo(w io.Writer) (int64, error) {
	n, error := fmt.Fprintf(w, "%d %s\r\n", r.Status, &r.URL)
	return int64(n), error
}

func NewRedirectResponse(status RedirectStatus, url url.URL) (*RedirectResponse, error) {
	if url.Scheme == "" || url.Host == "" {
		return nil, fmt.Errorf("invalid URL: %s", url.String())
	}

	url.RawFragment = ""
	url.User = nil

	return &RedirectResponse{
		Status: status,
		URL:    url,
	}, nil
}
