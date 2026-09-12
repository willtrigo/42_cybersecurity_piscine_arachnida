// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 18:23:57 by dande-je          #+#    #+#             //
//   Updated: 2026/09/12 14:36:17 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type BMPParser struct{}

func NewBMPParser() *BMPParser {
	return &BMPParser{}
}

func (BMPParser) Read(path string, stat application.StatMetadata, file application.FileReader) (metadata *domain.Metadata, err error) {
	data, err := file.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("bmp: %w", err)
	}

	fileStat, err := stat.FileStat(path)
	if err != nil {
		return nil, fmt.Errorf("bmp: %w", err)
	}

	header, err := decodeBMPHeader(data)
	if err != nil {
		return nil, fmt.Errorf("bmp: %w", err)
	}

	return &domain.Metadata{
		Path:             path,
		Format:           domain.FormatBMP,
		TagsSystem:       buildFileInfoTags(path, fileStat),
		TagsNoneEditable: header.tags(),
	}, nil
}
