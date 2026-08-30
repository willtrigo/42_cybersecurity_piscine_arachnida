// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image_panel.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 10:53:48 by dande-je          #+#    #+#             //
//   Updated: 2026/08/30 18:21:03 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const imageFrameStrokeWidth = 1

func newImagePanel(path string, format domain.Format, tags []domain.Tag) (fyne.CanvasObject, *animatedImage, error) {
	if format == domain.FormatGIF {
		return newAnimatedImagePanel(path, tags)
	}

	oriented, err := decodeOrientedImage(path, tags)
	if err != nil {
		return nil, nil, fmt.Errorf("loding image: %w", err)
	}

	image := canvas.NewImageFromImage(oriented)
	image.FillMode = canvas.ImageFillContain

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = imageFrameStrokeWidth

	return container.NewStack(frame, container.NewPadded(image)), nil, nil
}

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
