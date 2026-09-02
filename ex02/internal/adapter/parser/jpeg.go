// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   jpeg.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 08:43:09 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 18:20:59 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"

	jis "github.com/dsoprea/go-jpeg-image-structure/v2"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type JPEGParser struct{}

func NewJPEGParser() *JPEGParser {
	return &JPEGParser{}
}

func (JPEGParser) Read(path string) (*domain.Metadata, error) {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("jpeg: invalid file path")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("jpeg: %w", err)
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("jpeg: decoding header: %w", err)
	}

	var tags []domain.Tag
	if sl, err := jis.NewJpegMediaParser().ParseBytes(data); err == nil {
		segments := sl.(*jis.SegmentList)
		if _, _, exifTags, err := segments.DumpExif(); err == nil {
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

	return &domain.Metadata{
		Format:     domain.FormatJPEG,
		Dimensions: domain.Dimensions{Width: cfg.Width, Height: cfg.Height},
		Tags:       tags,
	}, nil
}

func (JPEGParser) SetTag() error {
	return nil
}

func (JPEGParser) DeleteTag(path, tag string) error {
	return nil
}
