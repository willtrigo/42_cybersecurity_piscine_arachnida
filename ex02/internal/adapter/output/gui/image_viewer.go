// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image_viewer.go                                    :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 15:08:22 by dande-je          #+#    #+#             //
//   Updated: 2026/09/03 11:17:03 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const imagePanelRatio = 0.77

type ImageViewer struct {
	window fyne.Window

	imageSlot    *fyne.Container
	metadataSlot *fyne.Container
	navigation   *Navigation
	save         *Save
	split        *container.Split

	animator *animatedImage

	editor  *application.MetadataEditor
	inspect application.InspectionReload
	results []application.InspectionResult
}

type refresher interface {
	Refresh() error
}

type errorPresenter interface {
	ShowError(error)
}

type viewerContext interface {
	refresher
	errorPresenter
}

type metadataEditor interface {
	setSaveVisibility(visible bool)
	DeleteTag(path string, tag string) error
	viewerContext
}

type navigationContext interface {
	show(idx int) error
	errorPresenter
}

func newImageViewer(window fyne.Window, results []application.InspectionResult, editor *application.MetadataEditor, inspect application.InspectionReload) (*ImageViewer, error) {
	imageViewer := &ImageViewer{
		window:       window,
		results:      results,
		editor:       editor,
		imageSlot:    container.NewStack(),
		metadataSlot: container.NewStack(),
		inspect:      inspect,
	}

	imageViewer.navigation = newNavigation(imageViewer, len(imageViewer.results)-1)
	imageViewer.save = newSave(imageViewer, 0)

	imageArea := container.NewStack(imageViewer.imageSlot, imageViewer.navigation.buttons)
	metadataArea := container.NewStack(imageViewer.metadataSlot, imageViewer.save.button)

	imageViewer.split = container.NewHSplit(imageArea, metadataArea)
	imageViewer.split.Offset = imagePanelRatio

	if err := imageViewer.show(0); err != nil {
		return nil, err
	}

	return imageViewer, nil
}

func (imageViewer *ImageViewer) Content() fyne.CanvasObject {
	return imageViewer.split
}

func (imageViewer *ImageViewer) show(idx int) error {
	if imageViewer.animator != nil {
		imageViewer.animator.Stop()
		imageViewer.animator = nil
	}

	result := imageViewer.results[idx]

	imagePanel, animator, err := newImagePanel(result.Metadata.Path, result.Metadata.Format, result.Metadata.Tags)
	if err != nil {
		return fmt.Errorf("%s: %w", result.Path, err)
	}

	imageViewer.navigation.idx = idx
	imageViewer.animator = animator

	imageViewer.save.idx = idx
	if result.Metadata.Format != domain.FormatBMP {
		imageViewer.setSaveVisibility(true)
	} else {
		imageViewer.setSaveVisibility(false)
	}

	imageViewer.imageSlot.Objects = []fyne.CanvasObject{imagePanel}
	imageViewer.imageSlot.Refresh()

	imageViewer.metadataSlot.Objects = []fyne.CanvasObject{newMetadataPanel(result, imageViewer)}
	imageViewer.metadataSlot.Refresh()

	imageViewer.navigation.updateNavigationState()
	imageViewer.window.SetTitle(windowTitle + " - " + result.Metadata.Path)

	if imageViewer.animator != nil {
		imageViewer.animator.Start()
	}

	return nil
}

func (imageViewer *ImageViewer) setSaveVisibility(visible bool) {
	if visible {
		imageViewer.save.button.Show()
	} else {
		imageViewer.save.button.Hide()
	}
}

func (imageViewer *ImageViewer) Close() {
	if imageViewer.animator != nil {
		imageViewer.animator.Stop()
	}
}

func (imageViewer *ImageViewer) Refresh() error {
	return imageViewer.show(imageViewer.navigation.idx)
}

func (ImageViewer *ImageViewer) ShowError(err error) {
	dialog.ShowError(err, ImageViewer.window)
}

func (ImageViewer *ImageViewer) DeleteTag(path string, tag string) error {
	format := ImageViewer.results[ImageViewer.navigation.idx].Metadata.Format
	if err := ImageViewer.editor.DeleteTag(path, format, tag); err != nil {
		return err
	}

	newResults, err := ImageViewer.inspect.InspectionResultReload()
	if err != nil {
		return err
	}
	ImageViewer.results = newResults
	return ImageViewer.Refresh()
}
