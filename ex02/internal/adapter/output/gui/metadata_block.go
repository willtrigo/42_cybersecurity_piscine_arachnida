// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   metadata_block.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 11:28:07 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 22:28:00 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	fieldBgColorR = 63
	fieldBgColorG = 63
	fieldBgColorB = 66
	fieldBgColorA = 255

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

func newBlockContainer(block []domain.Tag, edit bool, viewer metadataEditor, format domain.Format, path string) *fyne.Container {
	fieldBg := newBg(fieldBgColorR, fieldBgColorG, fieldBgColorB, fieldBgColorA)
	fieldBg.CornerRadius = cornerRadiusDefault

	blockContent := buildBlockContent(block, edit, viewer, format, path)
	return container.NewStack(fieldBg, blockContent)
}

func buildBlockContent(block []domain.Tag, edit bool, viewer metadataEditor, format domain.Format, path string) *fyne.Container {
	blockContent := container.NewVBox()

	for i, field := range block {
		fieldContainer := container.NewVBox()

		fieldName := buildFieldName(field.Name, i)
		fieldContainer.Add(fieldName)

		fieldContent := buildFieldContent(field.Value, edit, format)
		fieldContainer.Add(fieldContent)

		if edit && format != domain.FormatBMP && field.Value != "none found" {
			deleteField := newDelete(path, field.IDFPath, viewer)
			contentRow := container.NewBorder(nil, nil, nil, deleteField.button, fieldContainer)
			blockContent.Add(contentRow)
		} else {
			blockContent.Add(fieldContainer)
		}
		newDivisor(i, len(block), blockContent)
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
