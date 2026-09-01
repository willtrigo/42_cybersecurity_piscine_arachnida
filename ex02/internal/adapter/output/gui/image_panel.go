// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image_panel.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 10:53:48 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 15:51:25 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	imageBgColorR = 34
	imageBgColorG = 34
	imageBgColorB = 38
	imageBgColorA = 255
)

func newImagePanel(path string, format domain.Format, tags []domain.Tag) (fyne.CanvasObject, *animatedImage, error) {
	bg := newBg(imageBgColorR, imageBgColorG, imageBgColorB, imageBgColorA)

	if format == domain.FormatGIF {
		return newAnimatedImagePanel(path, tags)
	}

	oriented, err := decodeOrientedImage(path, tags)
	if err != nil {
		return nil, nil, fmt.Errorf("loading image: %w", err)
	}

	image := canvas.NewImageFromImage(oriented)
	image.FillMode = canvas.ImageFillContain

	return container.NewStack(bg, container.NewPadded(image)), nil, nil
}
