// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image_viewer.go                                    :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/01 15:08:22 by dande-je          #+#    #+#             //
//   Updated: 2026/09/01 15:59:24 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/application"
)

const imagePanelRatio = 0.77

type ImageViewer struct {
	window fyne.Window

	imageSlot    *fyne.Container
	metadataSlot *fyne.Container
	navigation   *Navigation
	split        *container.Split

	animator *animatedImage

	results []application.InspectionResult
}

func newImageViewer(window fyne.Window, results []application.InspectionResult) (*ImageViewer, error) {
	imageViewer := &ImageViewer{
		window:       window,
		results:      results,
		imageSlot:    container.NewStack(),
		metadataSlot: container.NewStack(),
	}

	imageViewer.navigation = newNavigation(imageViewer, len(imageViewer.results)-1)

	imageArea := container.NewStack(imageViewer.imageSlot, imageViewer.navigation.buttons)

	imageViewer.split = container.NewHSplit(imageArea, imageViewer.metadataSlot)
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

	imageViewer.imageSlot.Objects = []fyne.CanvasObject{imagePanel}
	imageViewer.imageSlot.Refresh()

	imageViewer.metadataSlot.Objects = []fyne.CanvasObject{newMetadataPanel(result)}
	imageViewer.metadataSlot.Refresh()

	imageViewer.navigation.updateNavigationState()
	imageViewer.window.SetTitle(windowTitle + " - " + result.Metadata.Path)

	if imageViewer.animator != nil {
		imageViewer.animator.Start()
	}

	return nil
}

func (imageViewer *ImageViewer) Close() {
	if imageViewer.animator != nil {
		imageViewer.animator.Stop()
	}
}
