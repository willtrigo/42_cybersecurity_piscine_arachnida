// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata_panel.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 10:59:25 by dande-je          #+#    #+#             //
//   Updated: 2026/08/30 18:21:05 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/adapter/output/format"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

func newMetadataPanel(result application.InspectionResult) fyne.CanvasObject {
	label := widget.NewLabel(format.RenderGUI(result))
	label.TextStyle = fyne.TextStyle{Monospace: true}
	label.Wrapping = fyne.TextWrapOff

	return container.NewScroll(label)
}
