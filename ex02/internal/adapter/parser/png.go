// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   png.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 10:25:58 by dande-je          #+#    #+#             //
//   Updated: 2026/08/27 11:02:00 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	exif "github.com/dsoprea/go-exif/v3"
	pis "github.com/dsoprea/go-png-image-structure/v2"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type PNGParser struct{}

func NewPNGParser() *PNGParser {
	return &PNGParser{}
}

func (PNGParser) Read(path string) (*domain.Metadata, error) {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("png: invalid file path")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("png: %w", err)
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("png: decoding header: %w", err)
	}

	var tags []domain.Tag
	if intfc, err := pis.NewPngMediaParser().ParseBytes(data); err == nil {
		chunks := intfc.(*pis.ChunkSlice)
		if chunk, err := chunks.FindExif(); err == nil {
			if exifTags, _, err := exif.GetFlatExifData(chunk.Data, nil); err == nil {
				tags = make([]domain.Tag, 0, len(exifTags))
				for _, t := range exifTags {
					tags = append(tags, domain.Tag{
						IDFPath: t.IfdPath,
						Name:    t.TagName,
						Value:   t.FormattedFirst,
					})
				}
			}
		}
	}

	return &domain.Metadata{
		Format:     domain.FormatPNG,
		Dimensions: domain.Dimensions{Width: cfg.Width, Height: cfg.Height},
		Tags:       tags,
	}, nil
}
