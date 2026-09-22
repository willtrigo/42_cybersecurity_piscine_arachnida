// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   xmp_tags.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/14 12:34:46 by dande-je          #+#    #+#             //
//   Updated: 2026/09/14 12:43:38 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"regexp"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type xmpField struct {
	pattern *regexp.Regexp
	name    string
}

var xmpFields = []xmpField{
	{regexp.MustCompile(`x:xmptk="([^"]*)"`), "XMP Toolkit"},
	{regexp.MustCompile(`xmp:CreatorTool="([^"]*)"`), "Creator Tool"},
	{regexp.MustCompile(`xmpMM:InstanceID="([^"]*)"`), "Instance ID"},
	{regexp.MustCompile(`xmpMM:DocumentID="([^"]*)"`), "Document ID"},
	{regexp.MustCompile(`stRef:instanceID="([^"]*)"`), "Derived From Instance ID"},
	{regexp.MustCompile(`stRef:documentID="([^"]*)"`), "Derived From Document ID"},
}

func buildXmpTags(packet []byte) []domain.Tag {
	if len(packet) == 0 {
		return []domain.Tag{}
	}

	xml := string(packet)

	tags := make([]domain.Tag, 0, len(xmpFields))
	for _, field := range xmpFields {
		if match := field.pattern.FindStringSubmatch(xml); match != nil {
			tags = append(tags, domain.NewTag(field.pattern.String(), field.name, match[1]))
		}
	}

	return tags
}
