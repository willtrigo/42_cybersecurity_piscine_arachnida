// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   ui.go                                              :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 11:19:52 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 12:29:41 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	cornerRadiusDefault = 12
)

func containerWithMinSize(content fyne.CanvasObject, width, height float32) fyne.CanvasObject {
	sizer := canvas.NewRectangle(color.Transparent)
	sizer.SetMinSize(fyne.NewSize(width, height))

	return container.NewStack(sizer, content)
}

func newBg(r, g, b, a uint8) *canvas.Rectangle {
	return canvas.NewRectangle(
		color.NRGBA{
			R: r,
			G: g,
			B: b,
			A: a,
		})
}
