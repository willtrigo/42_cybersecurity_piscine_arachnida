// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata_ui.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 11:59:45 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 12:05:06 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

const (
	metadataBgColorR = 46
	metadataBgColorG = 46
	metadataBgColorB = 50
	metadataBgColorA = 255

	divisorWidth  = 200
	divisorHeight = 1
)

func newDivisor(i, len int, blockContent *fyne.Container) {
	if i < len-1 {
		divisor := newBg(metadataBgColorR, metadataBgColorG, metadataBgColorB, metadataBgColorA)
		divisor.StrokeWidth = 0

		divisorContainer := container.NewWithoutLayout(divisor)
		divisorContainer.Resize(fyne.NewSize(divisorWidth, divisorHeight))

		divisorWithPadding := newPadded(3, 3, 0, 0, divisor)
		divisorWithPadding.Resize(fyne.NewSize(0, divisorHeight))

		blockContent.Add(divisorWithPadding)
	}
}

func newPadded(padTop, padBottom, padLeft, padRight float32, canvasObject fyne.CanvasObject) *fyne.Container {
	return container.New(
		layout.NewCustomPaddedLayout(
			padTop,
			padBottom,
			padLeft,
			padRight,
		),
		canvasObject,
	)
}
