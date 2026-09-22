// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif_tags.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/13 21:31:06 by dande-je          #+#    #+#             //
//   Updated: 2026/09/14 12:43:32 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	gifIFDPath = "GIF"
	xmpIFDPath = "XMP"
)

const centisecondsPerSecond = 100

func (h gifHeader) noneEditableTags() []domain.Tag {
	tags := []domain.Tag{
		newGIFTag("File Type", "GIF"),
		newGIFTag("File Type extension", "gif"),
		newGIFTag("MIME Type", "image/gif"),
		newGIFTag("GIF Version", h.Version),
		newGIFTag("Image Width", fmt.Sprintf("%d", h.Width)),
		newGIFTag("Image Height", fmt.Sprintf("%d", h.Height)),
		newGIFTag("Has Color Map", yesNo(h.HasGlobalColorTable)),
		newGIFTag("Color Resolution Depth", fmt.Sprintf("%d", h.ColorResolutionDepth)),
		newGIFTag("Bits Per Pixel", fmt.Sprintf("%d", h.BitsPerPixel)),
		newGIFTag("Background Color", fmt.Sprintf("%d", h.BackgroundColorIndex)),
	}

	if h.HasLoopCount {
		tags = append(tags, newGIFTag("Animation Iterations", animationIterations(h.LoopCount)))
	}
	if h.HasTransparantColor {
		tags = append(tags, newGIFTag("Transparent Color", fmt.Sprintf("%d", h.TransparentColorIndex)))
	}

	tags = append(tags,
		newGIFTag("Frame Count", fmt.Sprintf("%d", h.FrameCount)),
		newGIFTag("Duration", formatDuration(h.DurationCentiseconds)),
		newGIFTag("Image Size", fmt.Sprintf("%d x %d", h.Width, h.Height)),
		newGIFTag("Megapixels", formatMegapixels(h.Width, h.Height)),
	)

	return tags
}

func newGIFTag(name, value string) domain.Tag {
	return domain.NewTag(gifIFDPath, name, value)
}

func yesNo(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func animationIterations(loopCount uint16) string {
	if loopCount == 0 {
		return "Infinite"
	}
	return fmt.Sprintf("%d", loopCount)
}

func formatDuration(centiseconds int) string {
	return fmt.Sprintf("%.2f s", float64(centiseconds)/centisecondsPerSecond)
}

func editableTags(packet []byte) []domain.Tag {
	return buildXmpTags(packet)
}
