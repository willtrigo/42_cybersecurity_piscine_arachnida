// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   url.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/16 15:05:15 by dande-je          #+#    #+#             //
//   Updated: 2026/08/16 15:21:48 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import (
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
