// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   tags_format.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/13 21:16:34 by dande-je          #+#    #+#             //
//   Updated: 2026/09/24 17:13:35 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import "fmt"

const pixelsPerMegapixel = 1_000_000
const megapixelPrecisionThreshold = 1

func formatMegapixels(width, height uint32) string {
	megapixels := float64(int64(width)*int64(height)) / pixelsPerMegapixel
	if megapixels < megapixelPrecisionThreshold {
		return fmt.Sprintf("%.3f", megapixels)
	}
	return fmt.Sprintf("%.1f", megapixels)
}
