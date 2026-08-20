// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   url.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/16 15:05:15 by dande-je          #+#    #+#             //
//   Updated: 2026/08/19 23:33:44 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type URL struct {
	parsed *url.URL
	raw    string
}

func NewURL(raw string) (*URL, error) {
	if raw == "" {
		return nil, ErrEmptyURL
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %q: %v", ErrInvalidURL, raw, err)
	}

	if !parsed.IsAbs() || parsed.Host == "" {
		return nil, fmt.Errorf("%w: %q", ErrInvalidURL, raw)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedScheme, parsed.Scheme)
	}

	parsed.Fragment = ""

	return &URL{raw: parsed.String(), parsed: parsed}, nil
}

func (u *URL) String() string {
	return u.raw
}

func (u *URL) Host() string {
	return u.parsed.Host
}

func (u *URL) Path() string {
	return u.parsed.Path
}

func (u *URL) ResolveReference(ref string) (*URL, error) {
	if ref == "" {
		return nil, ErrEmptyURL
	}

	parsedRef, err := url.Parse(ref)
	if err != nil {
		return nil, fmt.Errorf("%w: %q: %v", ErrInvalidURL, ref, err)
	}

	resolved := u.parsed.ResolveReference(parsedRef)
	return NewURL(resolved.String())
}

func (u *URL) Normalize() (*URL, error) {
	s := u.raw

	if len(s) > 1 && strings.HasSuffix(s, "/") {
		s = s[:len(s)-1]
	}

	return NewURL(s)
}

type HTTPClient interface {
	Get(ctx context.Context, u *URL) (body []byte, contentType string, err error)
}

type ImageStorage interface {
	Save(ctx context.Context, img *Image) error
}

type HTMLParser interface {
	ExtractLinks(base *URL, body []byte) ([]*URL, error)

	ExtractImagesURLs(base *URL, body []byte) ([]*URL, error)
}
