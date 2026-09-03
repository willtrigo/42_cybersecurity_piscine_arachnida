// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata_block.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 11:28:07 by dande-je          #+#    #+#             //
//   Updated: 2026/09/03 11:14:00 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	infoBgColorR = 63
	infoBgColorG = 63
	infoBgColorB = 66
	infoBgColorA = 255

	fieldNamePadTopBegin    = 20
	fieldNamePadTopEnd      = 0
	fieldNamePadBottomBegin = 0
	fieldNamePadBottomEnd   = -5

	textPadLeft  = 10
	textPadRight = 10

	fieldContentPadTop         = -5
	fieldContentPadBottomBegin = 25
	fieldContentPadBottomEnd   = 0

	BlockContentPadTop    = 0
	BlockContentPadBottom = 20
	BlockContentPadLeft   = 20
	BlockContentPadRight  = 20
)

func newBlockContainer(fileName string, format domain.Format, block string, idxBlock int, viewer metadataEditor) *fyne.Container {
	bgInfo := newBg(infoBgColorR, infoBgColorG, infoBgColorB, infoBgColorA)
	bgInfo.CornerRadius = cornerRadiusDefault

	blockContent := buildBlockContent(fileName, format, block, idxBlock, viewer)
	return container.NewStack(bgInfo, blockContent)
}

func buildBlockContent(fileName string, format domain.Format, block string, idxBlock int, viewer metadataEditor) *fyne.Container {
	blockContent := container.NewVBox()

	lines := strings.Split(block, "\n")
	for i, line := range lines {

		if line == "" {
			continue
		}
		fieldParts := strings.Split(line, ":|:")

		if len(fieldParts) == 1 || fieldParts[1] == "#title" {
			fieldName := buildFieldName(fieldParts[0], i)
			blockContent.Add(fieldName)
		} else {
			edit := idxBlock >= 1

			fieldContainer := container.NewVBox()

			fieldName := buildFieldName(fieldParts[0], i)
			fieldContainer.Add(fieldName)

			if fieldParts[1] == "none found" {
				newDivisor(i, len(lines), fieldContainer)
				viewer.setSaveVisibility(false)
			}

			fieldContent := buildFieldContent(fieldParts[1], edit, format)
			fieldContainer.Add(fieldContent)

			if edit && format != domain.FormatBMP && fieldParts[1] != "none found" {
				deleteField := newDelete(fileName, fieldParts[0], viewer)
				contentRow := container.NewBorder(nil, nil, nil, deleteField.button, fieldContainer)
				blockContent.Add(contentRow)
			} else {
				blockContent.Add(fieldContainer)
			}
		}
		newDivisor(i, len(lines)-1, blockContent)
	}

	return newPadded(BlockContentPadTop, BlockContentPadBottom, BlockContentPadLeft, BlockContentPadRight, blockContent)
}

func buildFieldName(name string, idx int) *fyne.Container {
	fieldName := widget.NewLabel(name)
	fieldName.TextStyle = fyne.TextStyle{Monospace: true}
	fieldName.Wrapping = fyne.TextWrapWord

	var fieldNameWithPadding *fyne.Container
	if idx == 0 {
		fieldNameWithPadding = newPadded(fieldNamePadTopBegin, fieldNamePadBottomBegin, textPadLeft, textPadRight, fieldName)
	} else {
		fieldNameWithPadding = newPadded(fieldNamePadTopEnd, fieldNamePadBottomEnd, textPadLeft, textPadRight, fieldName)
	}

	return fieldNameWithPadding
}

func buildFieldContent(content string, idxBlock bool, format domain.Format) *fyne.Container {
	var fieldContent fyne.CanvasObject

	if idxBlock && format != domain.FormatBMP && content != "none found" {
		entry := widget.NewEntry()
		entry.SetText(content)
		entry.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		fieldContent = entry
	} else {
		label := widget.NewLabel(content)
		label.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		label.Wrapping = fyne.TextWrapWord
		fieldContent = label
	}

	fieldContentWithPadding := newPadded(fieldContentPadTop, fieldContentPadBottomEnd, textPadLeft, textPadRight, fieldContent)

	return fieldContentWithPadding
}
