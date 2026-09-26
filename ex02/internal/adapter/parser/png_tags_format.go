// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   png_tags_format.go                                 :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/24 22:15:23 by dande-je          #+#    #+#             //
//   Updated: 2026/09/24 22:15:47 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import "fmt"

func formatPngColorType(value uint8) string {
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

func formatPngCompression(value uint8) string {
	if value == 0 {
		return "Deflate/Inflate"
	}
	return fmt.Sprintf("Unknown (%d)", value)
}

func formatPngFilter(value uint8) string {
	if value == 0 {
		return "Adaptive"
	}
	return fmt.Sprintf("Unknown (%d)", value)
}

func formatPngInterlace(value uint8) string {
	switch value {
	case interlaceNone:
		return "Noninterlaced"
	case interlaceAdam7:
		return "Adam7 Interlace"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}
