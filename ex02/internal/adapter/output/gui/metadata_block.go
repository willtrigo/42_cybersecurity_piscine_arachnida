// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata_block.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 11:28:07 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 16:31:52 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
)

func newBlockContainer(block string, idxBlock int) *fyne.Container {
	bgInfo := newBg(infoBgColorR, infoBgColorG, infoBgColorB, infoBgColorA)
	bgInfo.CornerRadius = cornerRadiusDefault

	blockContent := buildBlockContent(block, idxBlock)
	return container.NewStack(bgInfo, blockContent)
}

func buildBlockContent(block string, idxBlock int) *fyne.Container {
	blockContent := container.NewVBox()

	lines := strings.Split(block, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		fieldParts := strings.Split(line, ":|:")

		fieldName := buildFieldName(fieldParts[0], i)
		blockContent.Add(fieldName)

		if len(fieldParts) != 1 {
			if fieldParts[1] != "#title" {
				edit := false
				if idxBlock >= 1 {
					edit = true
				}
				fieldContent := buildFieldContent(fieldParts[1], i+1, len(lines)-1, edit)
				blockContent.Add(fieldContent)
			}
		}
		newDivisor(i, len(lines)-1, blockContent)
	}

	return blockContent
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

func buildFieldContent(content string, idx, lenLines int, idxBlock bool) *fyne.Container {
	var fieldContent fyne.CanvasObject

	if idxBlock {
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

	var fieldContentWithPadding *fyne.Container
	if idx == lenLines {
		fieldContentWithPadding = newPadded(fieldContentPadTop, fieldContentPadBottomBegin, textPadLeft, textPadRight, fieldContent)
	} else {
		fieldContentWithPadding = newPadded(fieldContentPadTop, fieldContentPadBottomEnd, textPadLeft, textPadRight, fieldContent)
	}
	return fieldContentWithPadding
}
