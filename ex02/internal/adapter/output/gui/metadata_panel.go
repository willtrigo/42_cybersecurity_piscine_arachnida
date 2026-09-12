// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata_panel.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/30 10:59:25 by dande-je          #+#    #+#             //
//   Updated: 2026/09/12 02:21:42 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	containerPadTop    = 0
	containerPadBottom = 20
	containerPadLeft   = 0
	containerPadRight  = 0

	containerPadDefault = 20

	editableBlockIndex = 2
)

func newMetadataPanel(result application.InspectionResult, viewer metadataEditor) fyne.CanvasObject {
	bg := newBg(metadataBgColorR, metadataBgColorG, metadataBgColorB, metadataBgColorA)

	m := result.Metadata

	mainContainer := container.NewVBox()
	for i, block := range [][]domain.Tag{m.TagsSystem, m.TagsNoneEditable, m.TagsEditable} {
		if len(block) > 0 {
			edit := i == editableBlockIndex
			blockContainer := newBlockContainer(block, edit, viewer, m.Format)
			mainContainer.Add(newPadded(containerPadTop, containerPadBottom, containerPadLeft, containerPadRight, blockContainer))
		}
	}

	paddedContainer := newPadded(containerPadDefault, containerPadDefault, containerPadDefault, containerPadDefault, mainContainer)
	mainPanel := container.NewStack(bg, paddedContainer)

	return container.NewScroll(mainPanel)
}
