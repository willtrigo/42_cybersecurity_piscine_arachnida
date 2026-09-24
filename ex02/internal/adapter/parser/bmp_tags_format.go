// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp_tags_format.go                                 :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/24 16:31:00 by dande-je          #+#    #+#             //
//   Updated: 2026/09/24 16:43:42 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import "fmt"

const (
	bitmapCoreHeaderSize = 12
	bitmapInfoHeaderSize = 40
)

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

func formatBmpVersion(dibHeaderSize uint32) string {
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

func formatBmpCompression(value uint32) string {
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

func formatBmpColorCount(value uint32, zeroLabel string) string {
	if value == 0 {
		return zeroLabel
	}
	return fmt.Sprintf("%d", value)
}

func formatBmpMask(value uint32) string {
	return fmt.Sprintf("0x%08x", value)
}

func formatBmpColorSpace(value uint32) string {
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

func formatBmpRenderingIntent(value uint32) string {
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
