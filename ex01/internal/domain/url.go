// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   url.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/16 15:05:15 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 17:08:10 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import (
	"context"
	"fmt"
	"net/url"
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

type HTTPClient interface {
	Get(ctx context.Context, u *URL) (body []byte, contentType string, err error)
}

type ImageStorage interface {
	Save(ctx context.Context, img *Image) error
}
