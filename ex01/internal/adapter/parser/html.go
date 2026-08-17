// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   html.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 10:54:25 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 18:50:42 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"regexp"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

var (
	anchorTagPattern = regexp.MustCompile(`(?is)<a\b([^>]*)>`)
	imgTagPattern    = regexp.MustCompile(`(?is)<img\b([^>]*)>`)
	attrPattern      = regexp.MustCompile(`(?is)\b(\w[\w-]*)\s*=\s*"([^"]*)"|\b(\w[\w-]*)\s*=\s*'([^']*)'`)
)

type HTMLParser struct{}

func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

func (p *HTMLParser) ExtractLinks(base *domain.URL, body []byte) ([]*domain.URL, error) {
	return resolveAll(base, extractAttr(body, anchorTagPattern, "href")), nil
}

func (p *HTMLParser) ExtractImagesURLs(base *domain.URL, body []byte) ([]*domain.URL, error) {
	return resolveAll(base, extractAttr(body, imgTagPattern, "src")), nil
}

func extractAttr(html []byte, tagPattern *regexp.Regexp, attrName string) []string {
	var values []string

	for _, tag := range tagPattern.FindAllSubmatch(html, -1) {
		attrs := tag[1]
		for _, m := range attrPattern.FindAllSubmatch(attrs, -1) {
			name, value := attrNameValue(m)
			if strings.EqualFold(name, attrName) && value != "" {
				values = append(values, value)
			}
		}
	}

	return values
}

func attrNameValue(m [][]byte) (name, value string) {
	if len(m[1]) > 0 {
		return string(m[1]), string(m[2])
	}

	return string(m[3]), string(m[4])
}

func resolveAll(base *domain.URL, refs []string) []*domain.URL {
	urls := make([]*domain.URL, 0, len(refs))

	for _, ref := range refs {
		resolved, err := base.ResolveReference(ref)
		if err != nil {
			continue
		}
		urls = append(urls, resolved)
	}

	return urls
}
