// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 11:09:02 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 22:22:11 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type GIFParser struct{}

func NewGIFParser() *GIFParser {
	return &GIFParser{}
}

func (GIFParser) Read(path string, stat application.StatMetadata, file application.FileReader) (metadata *domain.Metadata, err error) {
	data, err := file.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}

	fileStat, err := stat.FileStat(path)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}

	header, err := decodeGIFHeader(data)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}

	return &domain.Metadata{
		Path:             path,
		Format:           domain.FormatGIF,
		TagsSystem:       buildSystemTags(path, fileStat),
		TagsNoneEditable: header.noneEditableTags(),
		TagsEditable:     editableTags(header.XMPPacket),
	}, nil
}

func (GIFParser) SetTag() error {
	return nil
}

func (GIFParser) DeleteTag(path, tag string) error {
	fmt.Printf("path: %v\n", path)
	fmt.Printf("tag: %v\n", tag)
	return nil
}
