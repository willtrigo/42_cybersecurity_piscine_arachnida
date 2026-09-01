// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image_orientation.go                               :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 10:56:39 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 14:57:19 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	exifOrientationTagName  = "Orientation"
	exifOrientationIdentity = 1

	orientationFlipHorizontal = 2
	orientationRotate180      = 3
	orientationFlipVertical   = 4
	orientationTranspose      = 5
	orientationRotate90CW     = 6
	orientationTransverse     = 7
	orientationRotate90CCW    = 8
)

func decodeOrientedImage(path string, tags []domain.Tag) (img image.Image, err error) {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("decode: invalid file path")
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()

	decoded, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", cleanPath, err)
	}

	return applyOrientation(decoded, exifOrientation(tags)), nil
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case orientationFlipHorizontal:
		return flipHorizontal(img)
	case orientationRotate180:
		return rotate180(img)
	case orientationFlipVertical:
		return flipVertical(img)
	case orientationTranspose:
		return transpose(img)
	case orientationRotate90CW:
		return rotate90CW(img)
	case orientationTransverse:
		return transverse(img)
	case orientationRotate90CCW:
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
