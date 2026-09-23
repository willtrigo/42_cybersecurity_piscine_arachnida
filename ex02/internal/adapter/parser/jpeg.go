// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   jpeg.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 08:43:09 by dande-je          #+#    #+#             //
//   Updated: 2026/09/22 21:42:18 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type JPEGParser struct{}

func NewJPEGParser() *JPEGParser {
	return &JPEGParser{}
}

func (JPEGParser) Read(path string, stat application.StatMetadata, file application.FileReader) (*domain.Metadata, error) {
	data, err := file.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("jpg: %w", err)
	}

	fileStat, err := stat.FileStat(path)
	if err != nil {
		return nil, fmt.Errorf("jpg: %w", err)
	}

	header, err := decodeJPEGHeader(data)
	if err != nil {
		return nil, fmt.Errorf("jpg: %w", err)
	}

	return &domain.Metadata{
		Path:             path,
		Format:           domain.FormatJPEG,
		TagsSystem:       buildSystemTags(path, fileStat),
		TagsNoneEditable: header.noneEditableTags(),
		TagsEditable:     header.editableTags(),
	}, nil
}

func (JPEGParser) SetTag() error {
	return nil
}

func (JPEGParser) DeleteTag(path, tag string) error {
	return nil
}
