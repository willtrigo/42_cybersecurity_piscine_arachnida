// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata_panel.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 10:59:25 by dande-je          #+#    #+#             //
//   Updated: 2026/09/02 20:17:03 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/adapter/output/format"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

const (
	containerPadTop    = 0
	containerPadBottom = 20
	containerPadLeft   = 0
	containerPadRight  = 0

	containerPadDefault = 20
)

func newMetadataPanel(result application.InspectionResult, viewer metadataEditor) fyne.CanvasObject {
	bg := newBg(metadataBgColorR, metadataBgColorG, metadataBgColorB, metadataBgColorA)

	blocks := format.RenderGUI(result)

	mainContainer := container.NewVBox()
	for i, block := range blocks {
		blockContainer := newBlockContainer(result.Metadata.Path, result.Metadata.Format, block, i, viewer)
		mainContainer.Add(newPadded(containerPadTop, containerPadBottom, containerPadLeft, containerPadRight, blockContainer))
	}

	paddedContainer := newPadded(containerPadDefault, containerPadDefault, containerPadDefault, containerPadDefault, mainContainer)
	mainPanel := container.NewStack(bg, paddedContainer)

	return container.NewScroll(mainPanel)
}
