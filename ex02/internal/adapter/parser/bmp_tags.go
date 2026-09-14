// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp_tags.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/13 20:58:20 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 23:34:14 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const bmpIFDpath = "BMP"

const (
	biRGB            = 0
	biRLE8           = 1
	biRLE4           = 2
	biBitfields      = 3
	biJPEG           = 4
	biPNG            = 5
	biAlphaBitfields = 6
	biCMYK           = 11
	biCMYKRLE8       = 12
	biCMYKRLE4       = 13
)

const (
	lcsCalibratedRGB = 0x00000000
	lcsSRGB          = 0x73524742
	lcsWindows       = 0x57696e20
	lcsProfileLinked = 0x4c494e4b
	lcsProfileEmbed  = 0x4d424544
)

const (
	lcsGMBusiness        = 1
	lcsGMGraphics        = 2
	lcsGMImages          = 4
	lcsGMAbsColorimetric = 8
)

func (h bmpHeader) noneEditableTags() []domain.Tag {
	tags := []domain.Tag{
		newBMPTag("File Type", "BMP"),
		newBMPTag("File Type Extension", "bmp"),
		newBMPTag("MIME Type", "image/bmp"),
		newBMPTag("BMP Version", bmpVersion(h.DIBHeaderSize)),
		newBMPTag("Image Width", fmt.Sprintf("%d", h.Width)),
		newBMPTag("Image Height", fmt.Sprintf("%d", h.Height)),
		newBMPTag("Planes", fmt.Sprintf("%d", h.Planes)),
		newBMPTag("Bit Depth", fmt.Sprintf("%d", h.BitCount)),
		newBMPTag("Compression", bmpCompression(h.Compression)),
		newBMPTag("Image Length", fmt.Sprintf("%d", h.SizeImage)),
		newBMPTag("Pixels Per Meter X", fmt.Sprintf("%d", h.XPelsPerMeter)),
		newBMPTag("Pixels Per Meter Y", fmt.Sprintf("%d", h.YPelsPerMeter)),
		newBMPTag("Num Colors", bmpColorCount(h.ClrUsed, "Use Bitdepth")),
		newBMPTag("Num Important Colors", bmpColorCount(h.ClrImportant, "All")),
	}

	if h.hasBitfieldMasks() {
		tags = append(tags,
			newBMPTag("Red Mask", bmpMask(h.RedMask)),
			newBMPTag("Green Mask", bmpMask(h.GreenMask)),
			newBMPTag("Blue Mask", bmpMask(h.BlueMask)),
		)
	}
	if h.hasAlphaMask() {
		tags = append(tags, newBMPTag("Alpha Mask", bmpMask(h.AlphaMask)))
	}
	if h.hasColorSpace() {
		tags = append(tags, newBMPTag("Color Space", bmpColorSpace(h.ColorSpace)))
	}
	if h.hasRenderingIntent() {
		tags = append(tags, newBMPTag("Rendering Intent", bmlRenderingIntent(h.RenderingIntent)))
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

func bmpVersion(dibHeaderSize uint32) string {
	switch dibHeaderSize {
	case bitmapCoreHeaderSize:
		return "Windows V2"
	case bitmapInfoHeaderSize:
		return "Windows V3"
	case bitmapV2HeaderSize, bitmapV3HeaderSize:
		return "Windows V3 (Adobe)"
	case bitmapV4HeaderSize:
		return "Windows V4"
	case bitmapV5HeaderSize:
		return "Windows V5"
	default:
		return fmt.Sprintf("Unknown (%d)", dibHeaderSize)
	}
}

func bmpCompression(value uint32) string {
	switch value {
	case biRGB:
		return "None"
	case biRLE8:
		return "8-bit RLE"
	case biRLE4:
		return "4-bit RLE"
	case biBitfields:
		return "Bitfields"
	case biJPEG:
		return "JPEG"
	case biPNG:
		return "PNG"
	case biAlphaBitfields:
		return "Alpha Bitfields"
	case biCMYK:
		return "CMYK"
	case biCMYKRLE8:
		return "CMYK 8-bit RLE"
	case biCMYKRLE4:
		return "CMYK 4-bit RLE"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func bmpColorCount(value uint32, zeroLabel string) string {
	if value == 0 {
		return zeroLabel
	}
	return fmt.Sprintf("%d", value)
}

func bmpMask(value uint32) string {
	return fmt.Sprintf("0x%08x", value)
}

func bmpColorSpace(value uint32) string {
	switch value {
	case lcsCalibratedRGB:
		return "Calibrated RGB"
	case lcsSRGB:
		return "sRGB"
	case lcsWindows:
		return "Windows Color Space"
	case lcsProfileLinked:
		return "Linked"
	case lcsProfileEmbed:
		return "Embedded"
	default:
		return fmt.Sprintf("Unknown (0x%08x)", value)
	}
}

func bmlRenderingIntent(value uint32) string {
	switch value {
	case lcsGMBusiness:
		return "Saturation (LGS_GM_BUSINESS)"
	case lcsGMGraphics:
		return "Graphic (LCS_GM_GRAPHICS)"
	case lcsGMImages:
		return "Picture (LCS_GM_IMAGES)"
	case lcsGMAbsColorimetric:
		return "Absolute (LCS_GM_ABS_COLORIMETRIC)"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}
