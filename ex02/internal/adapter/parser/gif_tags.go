// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif_tags.go                                        :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/13 21:31:06 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 23:33:48 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"
	"regexp"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	gifIFDPath = "GIF"
	xmpIFDPath = "XMP"
)

const centisecondsPerSecond = 100

type xmpField struct {
	pattern *regexp.Regexp
	name    string
}

var xmpFields = []xmpField{
	{regexp.MustCompile(`x:xmptk="([^"]*)"`), "XMP Toolkit"},
	{regexp.MustCompile(`xmp:CreatorTool="([^"]*)"`), "Creator Tool"},
	{regexp.MustCompile(`xmpMM:InstanceID="([^"]*)"`), "Instance ID"},
	{regexp.MustCompile(`xmpMM:DocumentID="([^"]*)"`), "Document ID"},
	{regexp.MustCompile(`stRef:instanceID="([^"]*)"`), "Derived From Instance ID"},
	{regexp.MustCompile(`stRef:documentID="([^"]*)"`), "Derived From Document ID"},
}

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
	if len(packet) == 0 {
		return []domain.Tag{}
	}

	xml := string(packet)

	tags := make([]domain.Tag, 0, len(xmpFields))
	for _, field := range xmpFields {
		if match := field.pattern.FindStringSubmatch(xml); match != nil {
			tags = append(tags, domain.NewTag(field.pattern.String(), field.name, match[1]))
		}
	}

	return tags
}
