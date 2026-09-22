// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   png.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 10:25:58 by dande-je          #+#    #+#             //
//   Updated: 2026/09/14 11:40:08 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type PNGParser struct{}

func NewPNGParser() *PNGParser {
	return &PNGParser{}
}

func (PNGParser) Read(path string, stat application.StatMetadata, file application.FileReader) (*domain.Metadata, error) {
	data, err := file.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}

	fileStat, err := stat.FileStat(path)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}

	header, err := decodePNGHeader(data)
	if err != nil {
		return nil, fmt.Errorf("png: %w", err)
	}

	return &domain.Metadata{
		Path:             path,
		Format:           domain.FormatPNG,
		TagsSystem:       buildSystemTags(path, fileStat),
		TagsNoneEditable: header.noneEditableTags(),
		TagsEditable:     header.editableTags(),
	}, nil
}

func (PNGParser) SetTag() error {
	return nil
}

func (PNGParser) DeleteTag(path, tag string) error {
	return nil
}
