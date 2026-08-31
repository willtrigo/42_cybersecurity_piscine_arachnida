// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata_panel.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 10:59:25 by dande-je          #+#    #+#             //
//   Updated: 2026/08/31 10:32:57 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/adapter/output/format"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

const (
	metadataBgColorR = 46
	metadataBgColorG = 46
	metadataBgColorB = 50
	metadataBgColorA = 255
)

func newMetadataPanel(result application.InspectionResult) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.NRGBA{R: metadataBgColorR, G: metadataBgColorG, B: metadataBgColorB, A: metadataBgColorA})

	label := widget.NewLabel(format.RenderGUI(result))
	label.TextStyle = fyne.TextStyle{Monospace: true}
	label.Wrapping = fyne.TextWrapOff

	metadata := container.NewStack(
		bg,
		label,
	)

	return container.NewScroll(metadata)
}
