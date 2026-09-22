// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   png_tags.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/14 12:23:11 by dande-je          #+#    #+#             //
//   Updated: 2026/09/22 18:11:37 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	pngIFDPath = "PNG"
	iccIFDPath = "ICC_Profile"
)

const (
	colorTypeGrayscale      = 0
	colorTypeRGB            = 2
	colorTypePalette        = 3
	colorTypeGrayscaleAlpha = 4
	colorTypeRGBA           = 6
)

const (
	interlaceNone  = 0
	interlaceAdam7 = 1
)

var pngTextTagNames = map[string]string{
	"Title":         "Title",
	"Author":        "Author",
	"Description":   "Description",
	"Copyright":     "Copyright",
	"Creation Time": "Creation Time",
	"Software":      "Software",
	"Disclaimer":    "Disclaimer",
	"Warning":       "Warning",
	"Source":        "Source",
	"Comment":       "Comment",
}

func (h pngHeader) noneEditableTags() []domain.Tag {
	return []domain.Tag{
		newPNGTag("File Type", "PNG"),
		newPNGTag("File Type Extension", "png"),
		newPNGTag("MIME Type", "image/png"),
		newPNGTag("Image Width", fmt.Sprintf("%d", h.Width)),
		newPNGTag("Image Height", fmt.Sprintf("%d", h.Height)),
		newPNGTag("Bit Depth", fmt.Sprintf("%d", h.BitDepth)),
		newPNGTag("Color Type", pngColorType(h.ColorType)),
		newPNGTag("Compression", pngCompression(h.CompressionMethod)),
		newPNGTag("Filter", pngFilter(h.FilterMethod)),
		newPNGTag("Interlace", pngInterlace(h.InterlaceMethod)),
		newPNGTag("Image Size", fmt.Sprintf("%d x %d", h.Width, h.Height)),
		newPNGTag("Megapixels", formatMegapixels(h.Width, h.Height)),
	}
}

func newPNGTag(name, value string) domain.Tag {
	return domain.NewTag(pngIFDPath, name, value)
}

func pngColorType(value uint8) string {
	switch value {
	case colorTypeGrayscale:
		return "Grayscale"
	case colorTypeRGB:
		return "RGB"
	case colorTypePalette:
		return "Palette"
	case colorTypeGrayscaleAlpha:
		return "Grayscale with Alpha"
	case colorTypeRGBA:
		return "RGB with Alpha"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func pngCompression(value uint8) string {
	if value == 0 {
		return "Deflate/Inflate"
	}
	return fmt.Sprintf("Unknown (%d)", value)
}

func pngFilter(value uint8) string {
	if value == 0 {
		return "Adaptive"
	}
	return fmt.Sprintf("Unknown (%d)", value)
}

func pngInterlace(value uint8) string {
	switch value {
	case interlaceNone:
		return "Noninterlaced"
	case interlaceAdam7:
		return "Adam7 Interlace"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func (h pngHeader) editableTags() []domain.Tag {
	tags := make([]domain.Tag, 0, len(h.TextEntries)+2)

	if h.ICCProfileName != "" {
		tags = append(tags, newICCTag("Profile Name", h.ICCProfileName))
	}

	for _, entry := range h.TextEntries {
		tags = append(tags, newPNGTag(pngTextTagName(entry.Keyword), entry.Value))
	}

	if len(h.XMPPacket) > 0 {
		tags = append(tags, buildXmpTags(h.XMPPacket)...)
	}

	if len(h.EXIFData) > 0 {
		if exif, err := decodeEXIFData(h.EXIFData); err == nil {
			tags = append(tags, exif.tags()...)
		}
	}

	return tags
}

func newICCTag(name, value string) domain.Tag {
	return domain.NewTag(iccIFDPath, name, value)
}

func pngTextTagName(keyword string) string {
	if name, ok := pngTextTagNames[keyword]; ok {
		return name
	}
	return keyword
}
