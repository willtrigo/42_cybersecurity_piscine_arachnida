// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   png.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 10:25:58 by dande-je          #+#    #+#             //
//   Updated: 2026/08/28 13:54:55 by dande-je         ###   ########.fr       //
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

	tags := extractPNGTags(data)

	return &domain.Metadata{
		Format:     domain.FormatPNG,
		Dimensions: domain.Dimensions{Width: cfg.Width, Height: cfg.Height},
		Tags:       tags,
	}, nil
}

func extractPNGTags(data []byte) []domain.Tag {
	var tags []domain.Tag

	intfc, err := pis.NewPngMediaParser().ParseBytes(data)
	if err != nil {
		return tags
	}
	chunks := intfc.(*pis.ChunkSlice)

	tags = append(tags, extractExifTags(chunks)...)
	tags = append(tags, extractTextTags(chunks)...)
	tags = append(tags, extractIHDRTags(chunks)...)

	return tags
}

func extractExifTags(chunks *pis.ChunkSlice) []domain.Tag {
	var tags []domain.Tag

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

	return tags
}

func extractTextTags(chunks *pis.ChunkSlice) []domain.Tag {
	var tags []domain.Tag

	for _, chunk := range chunks.Chunks() {
		switch chunk.Type {
		case "tEXt":
			parts := bytes.SplitN(chunk.Data, []byte{0}, 2)
			if len(parts) == 2 {
				key := string(parts[0])
				value := string(parts[1])
				if key != "" {
					tags = append(tags, domain.Tag{Name: key, Value: value})
				}
			}
		case "iTXt":
			parts := bytes.Split(chunk.Data, []byte{0})
			if len(parts) >= 2 {
				key := string(parts[0])
				value := string(parts[len(parts)-1])
				if key != "" {
					tags = append(tags, domain.Tag{Name: key, Value: value})
				}
			}
		}
	}

	return tags
}

func extractIHDRTags(chunks *pis.ChunkSlice) []domain.Tag {
	var tags []domain.Tag

	idx := chunks.Index()
	if ihdrChunks, ok := idx["IHDR"]; ok && len(ihdrChunks) > 0 {
		cd := pis.NewChunkDecoder()
		if decoded, err := cd.Decode(ihdrChunks[0]); err == nil {
			if ihdr, ok := decoded.(*pis.ChunkIHDR); ok {
				tags = append(tags,
					domain.Tag{Name: "BitDepth", Value: fmt.Sprintf("%d", ihdr.BitDepth)},
					domain.Tag{Name: "ColorType", Value: colorTypeName(ihdr.ColorType)},
					domain.Tag{Name: "Compression", Value: "Deflate/inflate"},
					domain.Tag{Name: "Filter", Value: "Adaptive"},
					domain.Tag{Name: "Interlace", Value: interlaceName(ihdr.InterlaceMethod)},
				)
			}
		}
	}

	return tags
}

func colorTypeName(ct uint8) string {
	switch ct {
	case 0:
		return "Grayscale"
	case 2:
		return "RGB"
	case 3:
		return "Palette"
	case 4:
		return "Grayscale+Alpha"
	case 6:
		return "RGBA"
	default:
		return fmt.Sprintf("Unknown(%d)", ct)
	}
}

func interlaceName(m uint8) string {
	if m == 0 {
		return "Noninterlaced"
	}
	return "Adam7"
}
