// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp_tags.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/13 20:58:20 by dande-je          #+#    #+#             //
//   Updated: 2026/09/24 16:43:38 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const bmpIFDpath = "BMP"

func (h bmpHeader) noneEditableTags() []domain.Tag {
	tags := []domain.Tag{
		newBMPTag("File Type", "BMP"),
		newBMPTag("File Type Extension", "bmp"),
		newBMPTag("MIME Type", "image/bmp"),
		newBMPTag("BMP Version", formatBmpVersion(h.DIBHeaderSize)),
		newBMPTag("Image Width", fmt.Sprintf("%d", h.Width)),
		newBMPTag("Image Height", fmt.Sprintf("%d", h.Height)),
		newBMPTag("Planes", fmt.Sprintf("%d", h.Planes)),
		newBMPTag("Bit Depth", fmt.Sprintf("%d", h.BitCount)),
		newBMPTag("Compression", formatBmpCompression(h.Compression)),
		newBMPTag("Image Length", fmt.Sprintf("%d", h.SizeImage)),
		newBMPTag("Pixels Per Meter X", fmt.Sprintf("%d", h.XPelsPerMeter)),
		newBMPTag("Pixels Per Meter Y", fmt.Sprintf("%d", h.YPelsPerMeter)),
		newBMPTag("Num Colors", formatBmpColorCount(h.ClrUsed, "Use Bitdepth")),
		newBMPTag("Num Important Colors", formatBmpColorCount(h.ClrImportant, "All")),
	}

	if h.hasBitfieldMasks() {
		tags = append(tags,
			newBMPTag("Red Mask", formatBmpMask(h.RedMask)),
			newBMPTag("Green Mask", formatBmpMask(h.GreenMask)),
			newBMPTag("Blue Mask", formatBmpMask(h.BlueMask)),
		)
	}
	if h.hasAlphaMask() {
		tags = append(tags, newBMPTag("Alpha Mask", formatBmpMask(h.AlphaMask)))
	}
	if h.hasColorSpace() {
		tags = append(tags, newBMPTag("Color Space", formatBmpColorSpace(h.ColorSpace)))
	}
	if h.hasRenderingIntent() {
		tags = append(tags, newBMPTag("Rendering Intent", formatBmpRenderingIntent(h.RenderingIntent)))
	}

	tags = append(tags,
		newBMPTag("Image Size", fmt.Sprintf("%d x %d", h.Width, h.Height)),
		newBMPTag("Megapixels", formatMegapixels(h.Width, h.Height)),
	)

	return tags
}

func newBMPTag(name, value string) domain.Tag {
	return domain.NewTag(bmpIFDpath, name, value)
}
