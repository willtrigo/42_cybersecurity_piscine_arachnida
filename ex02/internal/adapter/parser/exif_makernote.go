// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   exif_makernote.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/22 21:28:44 by dande-je          #+#    #+#             //
//   Updated: 2026/09/22 22:35:21 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	appleMakerNoteSignature = "Apple iOS\x00"
	appleMakerNoteHeaderLen = 14
)

const (
	tagMakerNoteVersion      = 0x0001
	tagRunTimeFlags          = 0x0003
	tagRunTimeValue          = 0x0004
	tagRunTimeScale          = 0x0005
	tagRunTimeEpoch          = 0x0006
	tagAEStable              = 0x0008
	tagAETarget              = 0x0009
	tagAEAverage             = 0x000A
	tagAFStable              = 0x000B
	tagAccelerationVector    = 0x000C
	tagFocusDistanceRange    = 0x000E
	tagImageCaptureType      = 0x0011
	tagLivePhotoVideoIndex   = 0x0013
	tagPhotosAppFeatureFlags = 0x0015
	tagHDRHeadroom           = 0x0016
	tagSignalToNoiseRatio    = 0x0017
	tagPhotoIdentifier       = 0x0019
	tagColorTemperature      = 0x001B
	tagCameraType            = 0x001D
	tagFocusPosition         = 0x001F
	tagRunTimeSincePowerUp   = 0x0027
)

func decodeMakerNote(raw []byte, order binary.ByteOrder, make string) []domain.Tag {
	if !strings.HasPrefix(make, "Apple") {
		return nil
	}
	if len(raw) < appleMakerNoteHeaderLen {
		return nil
	}
	if string(raw[:len(appleMakerNoteSignature)]) != appleMakerNoteSignature {
		return nil
	}

	body := raw[appleMakerNoteHeaderLen:]
	entries, _, err := readIFD(body, 0, order)
	if err != nil {
		return nil
	}

	var tags []domain.Tag
	for tagID, name := range appleMakerNoteNames {
		if e, ok := entries[tagID]; ok {
			if formatted, ok := formatAppleTag(tagID, e); ok {
				tags = append(tags, newEXIFTag(name, formatted))
			}
		}
	}
	return tags
}

var appleMakerNoteNames = map[uint16]string{
	tagMakerNoteVersion:      "Maker Note Version",
	tagRunTimeFlags:          "Run Time Flags",
	tagRunTimeValue:          "Run Time Value",
	tagRunTimeScale:          "Run Time Scale",
	tagRunTimeEpoch:          "Run Time Epoch",
	tagAEStable:              "AE Stable",
	tagAETarget:              "AE Target",
	tagAEAverage:             "AE Average",
	tagAFStable:              "AF Stable",
	tagAccelerationVector:    "Acceleration Vector",
	tagFocusDistanceRange:    "Focus Distance Range",
	tagImageCaptureType:      "Image Capture Type",
	tagLivePhotoVideoIndex:   "Live Photo Video Index",
	tagPhotosAppFeatureFlags: "Photos App Feature Flags",
	tagHDRHeadroom:           "HDR Headroom",
	tagSignalToNoiseRatio:    "Signal To Noise Ratio",
	tagPhotoIdentifier:       "Photo Identifier",
	tagColorTemperature:      "Color Temperature",
	tagCameraType:            "Camera Type",
	tagFocusPosition:         "Focus Position",
	tagRunTimeSincePowerUp:   "Run Time Since Power Up",
}

func formatAppleTag(id uint16, e tiffEntry) (string, bool) {
	switch id {
	case tagMakerNoteVersion:
		v, ok := e.shortValue()
		return strconv.FormatUint(uint64(v), 10), ok
	case tagRunTimeFlags:
		v, ok := e.shortValue()
		return runTimeFlagsString(v), ok
	case tagRunTimeValue, tagRunTimeScale, tagRunTimeEpoch,
		tagAETarget, tagAEAverage, tagPhotosAppFeatureFlags,
		tagHDRHeadroom, tagSignalToNoiseRatio, tagColorTemperature,
		tagFocusPosition:
		if v, ok := e.numericValue(); ok {
			return strconv.FormatUint(uint64(v), 10), true
		}
	case tagAEStable, tagAFStable:
		v, ok := e.shortValue()
		return yesNoString(v), ok
	case tagAccelerationVector:
		return formatSignedRationalTriplet(e.sRationals()), true
	case tagFocusDistanceRange:
		return formatFocusDistanceRange(e.rationals()), true
	case tagImageCaptureType:
		v, ok := e.shortValue()
		return imageCaptureTypeString(v), ok
	case tagLivePhotoVideoIndex:
		v, ok := e.longValue()
		return strconv.FormatUint(uint64(v), 10), ok
	case tagPhotoIdentifier:
		return e.ascii(), true
	case tagCameraType:
		v, ok := e.shortValue()
		return cameraTypeString(v), ok
	case tagRunTimeSincePowerUp:
		v, ok := e.longValue()
		return formatDurationSeconds(uint64(v)), ok
	}
	return "", false
}

func runTimeFlagsString(v uint16) string {
	if v == 0x0001 {
		return "Valid"
	}
	return fmt.Sprintf("Unknown (%d)", v)
}

func yesNoString(v uint16) string {
	if v == 1 {
		return "Yes"
	}
	return "No"
}

func formatSignedRationalTriplet(r []sRational) string {
	parts := make([]string, len(r))
	for i, v := range r {
		parts[i] = formatSRational(v)
	}
	return strings.Join(parts, " ")
}

func formatFocusDistanceRange(r []rational) string {
	if len(r) < 2 {
		return ""
	}
	return fmt.Sprintf("%s - %s m", formatRational(r[0]), formatRational(r[1]))
}

func imageCaptureTypeString(v uint16) string {
	switch v {
	case 1:
		return "Standard"
	case 2:
		return "Portrait"
	case 3:
		return "Landscape"
	case 10:
		return "ProRAW"
	default:
		return fmt.Sprintf("Unknown (%d)", v)
	}
}

func cameraTypeString(v uint16) string {
	switch v {
	case 0:
		return "Back Wide Angle"
	case 1:
		return "Back Normal"
	case 2:
		return "Back Telephoto"
	case 6:
		return "Front"
	default:
		return fmt.Sprintf("Unknown (%d)", v)
	}
}

func formatDurationSeconds(total uint64) string {
	days := total / 86400
	hours := (total % 86400) / 3600
	mins := (total % 3600) / 60
	secs := total % 60
	return fmt.Sprintf("%d days %02d:%02d:%02d", days, hours, mins, secs)
}
