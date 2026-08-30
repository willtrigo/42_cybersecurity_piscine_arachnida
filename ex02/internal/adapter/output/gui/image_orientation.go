// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image_orientation.go                               :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 10:56:39 by dande-je          #+#    #+#             //
//   Updated: 2026/08/30 18:21:00 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"image"
	"strconv"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	exifOrientationTagName  = "Orientation"
	exifOrientationIdentity = 1
)

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return flipHorizontal(img)
	case 3:
		return rotate180(img)
	case 4:
		return flipVertical(img)
	case 5:
		return transpose(img)
	case 6:
		return rotate90CW(img)
	case 7:
		return transverse(img)
	case 8:
		return rotate90CCW(img)
	default:
		return img
	}
}

func flipHorizontal(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(width-1-x, y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return dst
}

func rotate180(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(width-1-x, height-1-y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return dst
}

func flipVertical(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(x, height-1-y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return dst
}

func transpose(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, height, width))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(y, x, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return dst
}

func rotate90CW(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, height, width))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(height-1-y, x, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return dst
}

func transverse(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, height, width))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(height-1-y, width-1-x, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return dst
}

func rotate90CCW(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, height, width))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(y, width-1-x, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
	return dst
}

func exifOrientation(tags []domain.Tag) int {
	for _, tag := range tags {
		if tag.Name != exifOrientationTagName {
			continue
		}
		if value, err := strconv.Atoi(strings.TrimSpace(tag.Value)); err == nil {
			return value
		}
	}
	return exifOrientationIdentity
}
