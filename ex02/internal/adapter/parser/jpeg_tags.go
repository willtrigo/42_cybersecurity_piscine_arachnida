// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   jpeg_tags.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/22 21:10:38 by dande-je          #+#    #+#             //
//   Updated: 2026/09/22 22:55:46 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const jpegIFDPath = "JFIF"

func (h jpegHeader) noneEditableTags() []domain.Tag {
	return []domain.Tag{
		newJPEGTag("File Type", "JPEG"),
		newJPEGTag("File Type Extension", "jpg"),
		newJPEGTag("MIME Type", "image/jpeg"),
		newJPEGTag("Image Width", fmt.Sprintf("%d", h.Width)),
		newJPEGTag("Image Height", fmt.Sprintf("%d", h.Height)),
		newJPEGTag("Encoding Process", h.EncodingProcess),
		newJPEGTag("Bits Per Sample", fmt.Sprintf("%d", h.BitsPerSample)),
		newJPEGTag("Color Components", fmt.Sprintf("%d", h.NumComponents)),
		newJPEGTag("Y Cb Cr Sub Sampling", jpegSubSampling(h)),
		newJPEGTag("Image Size", fmt.Sprintf("%d x %d", h.Width, h.Height)),
		newJPEGTag("Megapixels", formatMegapixels(uint32(h.Width), uint32(h.Height))),
	}
}

func (h jpegHeader) editableTags() []domain.Tag {
	tags := make([]domain.Tag, 0)

	if len(h.ICCProfile) > 0 {
		if icc, err := decodeICCProfile(h.ICCProfile); err == nil {
			tags = append(tags, icc.tags()...)
		}
	}

	if len(h.EXIFData) > 0 {
		if exif, err := decodeEXIFData(h.EXIFData); err == nil {
			tags = append(tags, exif.tags()...)
			tags = append(tags, buildCompositeTags(exif)...)
		}
	}

	if len(h.XMPPacket) > 0 {
		tags = append(tags, buildXmpTags(h.XMPPacket)...)
	}

	if h.Comment != "" {
		tags = append(tags, newJPEGTag("Comment", h.Comment))
	}

	return tags
}

func newJPEGTag(name, value string) domain.Tag {
	return domain.NewTag(jpegIFDPath, name, value)
}

func jpegSubSampling(h jpegHeader) string {
	if len(h.Components) < 3 {
		return ""
	}
	y := h.Components[0]
	switch {
	case y.HSampling == 2 && y.VSampling == 2:
		return "YCbCr4:2:0 (2 2)"
	case y.HSampling == 2 && y.VSampling == 1:
		return "YCbCr4:2:2 (2 1)"
	case y.HSampling == 4 && y.VSampling == 1:
		return "YCbCr4:1:1 (4 1)"
	case y.HSampling == 1 && y.VSampling == 1:
		return "YCbCr4:4:4 (1 1)"
	default:
		return fmt.Sprintf("YCbCr (%d %d)", y.HSampling, y.VSampling)
	}
}
