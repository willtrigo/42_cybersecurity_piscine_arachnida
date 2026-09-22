// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   exif_tags.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/22 18:02:52 by dande-je          #+#    #+#             //
//   Updated: 2026/09/22 20:00:45 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const exifIFDPath = "EXIF"

const (
	tagMake             = 0x010F
	tagModel            = 0x0110
	tagOrientation      = 0x0112
	tagXResolution      = 0x011A
	tagYResolution      = 0x011B
	tagResolutionUnit   = 0x0128
	tagSoftware         = 0x0131
	tagModifyDate       = 0x0132
	tagYCbCrPositioning = 0x0213
)

const (
	tagExposureTime            = 0x829A
	tagFNumber                 = 0x829D
	tagExposureProgram         = 0x8822
	tagISO                     = 0x8827
	tagExifVersion             = 0x9000
	tagDateTimeOriginal        = 0x9003
	tagCreateDate              = 0x9004
	tagOffsetTime              = 0x9010
	tagOffsetTimeOriginal      = 0x9011
	tagOffsetTimeDigitized     = 0x9012
	tagComponentsConfiguration = 0x9101
	tagShutterSpeedValue       = 0x9201
	tagApertureValue           = 0x9202
	tagBrightnessValue         = 0x9203
	tagExposureBiasValue       = 0x9204
	tagMaxApertureValue        = 0x9205
	tagSubjectDistance         = 0x9206
	tagMeteringMode            = 0x9207
	tagFlash                   = 0x9209
	tagFocalLength             = 0x920A
	tagSubSecTime              = 0x9291
	tagSubSecTimeOriginal      = 0x9292
	tagSubSecTimeDigitized     = 0x9293
	tagFlashpixVersion         = 0xA000
	tagPixelXDimension         = 0xA002
	tagPixelYDimension         = 0xA003
	tagSensingMethod           = 0xA217
	tagSceneType               = 0xA301
	tagCustomRendered          = 0xA401
	tagExposureMode            = 0xA402
	tagWhiteBalance            = 0xA403
	tagDigitalZoomRatio        = 0xA404
	tagFocalLengthIn35mm       = 0xA405
	tagSceneCaptureType        = 0xA406
	tagContrast                = 0xA408
	tagSaturation              = 0xA409
	tagSharpness               = 0xA40A
	tagSubjectDistanceRange    = 0xA40C
	tagLensMake                = 0xA433
	tagLensModel               = 0xA434
	tagCompositeImage          = 0xA460
)

const (
	tagGPSLatitudeRef     = 0x0001
	tagGPSLatitude        = 0x0002
	tagGPSLongitudeRef    = 0x0003
	tagGPSLongitude       = 0x0004
	tagGPSAltitudeRef     = 0x0005
	tagGPSAltitude        = 0x0006
	tagGPSTimeStamp       = 0x0007
	tagGPSImgDirectionRef = 0x0010
	tagGPSImgDirection    = 0x0011
	tagGPSDateStamp       = 0x001D
)

var flashDescriptions = map[uint16]string{
	0x00: "No Flash",
	0x01: "Fired",
	0x05: "Fired, Return not detected",
	0x07: "Fired, Return detected",
	0x08: "On, Did not fire",
	0x09: "On, Fired",
	0x0d: "On, Return not detected",
	0x0f: "On, Return detected",
	0x10: "Off, Did not fire",
	0x14: "Off, Did not fire, Return not detected",
	0x18: "Auto, Did not fire",
	0x19: "Auto, Fired",
	0x1d: "Auto, Fired, Return not detected",
	0x1f: "Auto, Fired, Return detected",
	0x20: "No flash function",
	0x30: "Off, No flash function",
	0x41: "Fired, Red-eye reduction",
	0x45: "Fired, Red-eye reduction, Return not detected",
	0x47: "Fired, Red-eye reduction, Return detected",
	0x49: "On, Red-eye reduction",
	0x4d: "On, Red-eye reduction, Return not detected",
	0x4f: "On, Red-eye reduction, Return detected",
	0x50: "Off, Red-eye reduction",
	0x58: "Auto, Did not fire, Red-eye reduction",
	0x59: "Auto, Fired, Red-eye reduction",
	0x5d: "Auto, Fired, Red-eye reduction, Return not detected",
	0x5f: "Auto, Fired, Red-eye reduction, Return detected",
}

func (h exifHeader) tags() []domain.Tag {
	tags := make([]domain.Tag, 0)
	tags = append(tags, h.ifd0Tags()...)
	tags = append(tags, h.exifSubIFDTags()...)
	tags = append(tags, h.gpsIFDTags()...)
	return tags
}

func (h exifHeader) ifd0Tags() []domain.Tag {
	var tags []domain.Tag
	m := h.IFD0

	if v, ok := asciiEntry(m, tagMake); ok {
		tags = append(tags, newEXIFTag("Make", v))
	}
	if v, ok := asciiEntry(m, tagModel); ok {
		tags = append(tags, newEXIFTag("Camera Model Name", v))
	}
	if e, ok := m[tagOrientation]; ok {
		if v, ok := e.shortValue(); ok {
			tags = append(tags, newEXIFTag("Orientation", orientationString(v)))
		}
	}
	if e, ok := m[tagXResolution]; ok {
		if v, ok := e.rationalValue(); ok {
			tags = append(tags, newEXIFTag("X Resolution", formatRational(v)))
		}
	}
	if e, ok := m[tagYResolution]; ok {
		if v, ok := e.rationalValue(); ok {
			tags = append(tags, newEXIFTag("Y Resolution", formatRational(v)))
		}
	}
	if e, ok := m[tagResolutionUnit]; ok {
		if v, ok := e.shortValue(); ok {
			tags = append(tags, newEXIFTag("Resolution Unit", resolutionUnitString(v)))
		}
	}
	if v, ok := asciiEntry(m, tagSoftware); ok {
		tags = append(tags, newEXIFTag("Software", v))
	}
	if v, ok := asciiEntry(m, tagModifyDate); ok {
		tags = append(tags, newEXIFTag("Modify Date", v))
	}
	if e, ok := m[tagYCbCrPositioning]; ok {
		if v, ok := e.shortValue(); ok {
			tags = append(tags, newEXIFTag("Y Cb Cr Positioning", yCbCrPositioningString(v)))
		}
	}

	return tags
}

func (h exifHeader) exifSubIFDTags() []domain.Tag {
	m := h.ExifIFD
	if m == nil {
		return nil
	}

	var tags []domain.Tag

	tags = appendRationalTag(tags, m, tagExposureTime, "Exposure Time", formatExposureTime)
	tags = appendRationalTag(tags, m, tagFNumber, "F Number", formatRational)
	tags = appendRationalTag(tags, m, tagApertureValue, "Aperture Value", formatApertureAPEX)
	tags = appendRationalTag(tags, m, tagMaxApertureValue, "Max Aperture Value", formatApertureAPEX)
	tags = appendRationalTag(tags, m, tagSubjectDistance, "Subject Distance", formatSubjectDistance)
	tags = appendRationalTag(tags, m, tagFocalLength, "Focal Length", func(r rational) string {
		return formatRational(r) + " mm"
	})
	tags = appendRationalTag(tags, m, tagDigitalZoomRatio, "Digital Zoom Ratio", formatRational)

	tags = appendSRationalTag(tags, m, tagShutterSpeedValue, "Shutter Speed Value", formatShutterSpeedAPEX)
	tags = appendSRationalTag(tags, m, tagBrightnessValue, "Brightness Value", formatSRational)
	tags = appendSRationalTag(tags, m, tagExposureBiasValue, "Exposure Compensation", formatSRational)

	tags = appendShortTag(tags, m, tagExposureProgram, "Exposure Program", exposureProgramString)
	tags = appendShortTag(tags, m, tagMeteringMode, "Metering Mode", meteringModeString)
	tags = appendShortTag(tags, m, tagFlash, "Flash", flashString)
	tags = appendShortTag(tags, m, tagSensingMethod, "Sensing Method", sensingMethodString)
	tags = appendShortTag(tags, m, tagCustomRendered, "Custom Rendered", customRenderedString)
	tags = appendShortTag(tags, m, tagExposureMode, "Exposure Mode", exposureModeString)
	tags = appendShortTag(tags, m, tagWhiteBalance, "White Balance", whiteBalanceString)
	tags = appendShortTag(tags, m, tagSceneCaptureType, "Scene Capture Type", sceneCaptureTypeString)
	tags = appendShortTag(tags, m, tagContrast, "Contrast", lowNormalHighString)
	tags = appendShortTag(tags, m, tagSaturation, "Saturation", lowNormalHighString)
	tags = appendShortTag(tags, m, tagSharpness, "Sharpness", sharpnessString)
	tags = appendShortTag(tags, m, tagSubjectDistanceRange, "Subject Distance Range", subjectDistanceRangeString)
	tags = appendShortTag(tags, m, tagCompositeImage, "Composite Image", compositeImageString)

	tags = appendShortSuffixTag(tags, m, tagFocalLengthIn35mm, "Focal Length In 35mm Format", " mm")

	tags = appendNumericTag(tags, m, tagISO, "ISO")
	tags = appendNumericTag(tags, m, tagPixelXDimension, "Exif Image Width")
	tags = appendNumericTag(tags, m, tagPixelYDimension, "Exif Image Height")

	tags = appendASCIITag(tags, m, tagDateTimeOriginal, "Date/Time Original")
	tags = appendASCIITag(tags, m, tagCreateDate, "Create Date")
	tags = appendASCIITag(tags, m, tagOffsetTime, "Offset Time")
	tags = appendASCIITag(tags, m, tagOffsetTimeOriginal, "Offset Time Original")
	tags = appendASCIITag(tags, m, tagOffsetTimeDigitized, "Offset Time Digitized")
	tags = appendASCIITag(tags, m, tagSubSecTime, "Sub Sec Time")
	tags = appendASCIITag(tags, m, tagSubSecTimeOriginal, "Sub Sec Time Original")
	tags = appendASCIITag(tags, m, tagSubSecTimeDigitized, "Sub Sec Time Digitized")
	tags = appendASCIITag(tags, m, tagLensMake, "Lens Make")
	tags = appendASCIITag(tags, m, tagLensModel, "Lens Model")

	tags = appendRawStringTag(tags, m, tagExifVersion, "Exif Version")
	tags = appendRawStringTag(tags, m, tagFlashpixVersion, "Flashpix Version")

	if e, ok := m[tagComponentsConfiguration]; ok {
		tags = append(tags, newEXIFTag("Components Configuration", componentsConfigurationString(e.raw)))
	}
	if e, ok := m[tagSceneType]; ok && len(e.raw) > 0 {
		tags = append(tags, newEXIFTag("Scene Type", sceneTypeString(e.raw[0])))
	}

	return tags
}

func appendASCIITag(tags []domain.Tag, m map[uint16]tiffEntry, tag uint16, name string) []domain.Tag {
	if v, ok := asciiEntry(m, tag); ok {
		tags = append(tags, newEXIFTag(name, v))
	}
	return tags
}

func appendRawStringTag(tags []domain.Tag, m map[uint16]tiffEntry, tag uint16, name string) []domain.Tag {
	if v, ok := versionEntry(m, tag); ok {
		tags = append(tags, newEXIFTag(name, v))
	}
	return tags
}

func appendShortTag(tags []domain.Tag, m map[uint16]tiffEntry, tag uint16, name string, fn func(uint16) string) []domain.Tag {
	if e, ok := m[tag]; ok {
		if v, ok := e.shortValue(); ok {
			tags = append(tags, newEXIFTag(name, fn(v)))
		}
	}
	return tags
}

func appendNumericTag(tags []domain.Tag, m map[uint16]tiffEntry, tag uint16, name string) []domain.Tag {
	if e, ok := m[tag]; ok {
		if v, ok := e.numericValue(); ok {
			tags = append(tags, newEXIFTag(name, strconv.FormatUint(uint64(v), 10)))
		}
	}
	return tags
}

func appendRationalTag(tags []domain.Tag, m map[uint16]tiffEntry, tag uint16, name string, fn func(rational) string) []domain.Tag {
	if e, ok := m[tag]; ok {
		if v, ok := e.rationalValue(); ok {
			tags = append(tags, newEXIFTag(name, fn(v)))
		}
	}
	return tags
}

func appendSRationalTag(tags []domain.Tag, m map[uint16]tiffEntry, tag uint16, name string, fn func(sRational) string) []domain.Tag {
	if e, ok := m[tag]; ok {
		if v, ok := e.sRationalValue(); ok {
			tags = append(tags, newEXIFTag(name, fn(v)))
		}
	}
	return tags
}

func appendShortSuffixTag(tags []domain.Tag, m map[uint16]tiffEntry, tag uint16, name, suffix string) []domain.Tag {
	return appendShortTag(tags, m, tag, name, func(v uint16) string {
		return fmt.Sprintf("%d%s", v, suffix)
	})
}

func (h exifHeader) gpsIFDTags() []domain.Tag {
	var tags []domain.Tag
	m := h.GPSIFD
	if m == nil {
		return tags
	}

	latRef, hasLatRef := asciiEntry(m, tagGPSLatitudeRef)
	if hasLatRef {
		tags = append(tags, newEXIFTag("GPS Latitude Ref", gpsRefName(latRef)))
	}
	if e, ok := m[tagGPSLatitude]; ok {
		tags = append(tags, newEXIFTag("GPS Latitude", formatGPSCoordinate(e.rationals(), latRef)))
	}

	lonRef, hasLonRef := asciiEntry(m, tagGPSLongitudeRef)
	if hasLonRef {
		tags = append(tags, newEXIFTag("GPS Longitude Ref", gpsRefName(lonRef)))
	}
	if e, ok := m[tagGPSLongitude]; ok {
		tags = append(tags, newEXIFTag("GPS Longitude", formatGPSCoordinate(e.rationals(), lonRef)))
	}

	altRef, hasAltRef := gpsByteEntry(m, tagGPSAltitudeRef)
	if hasAltRef {
		tags = append(tags, newEXIFTag("GPS Altitude Ref", gpsAltitudeRefName(altRef)))
	}
	if e, ok := m[tagGPSAltitude]; ok {
		if v, ok := e.rationalValue(); ok {
			tags = append(tags, newEXIFTag("GPS Altitude", formatGPSAltitude(v, altRef)))
		}
	}

	if e, ok := m[tagGPSTimeStamp]; ok {
		tags = append(tags, newEXIFTag("GPS Time Stamp", formatGPSTimeStamp(e.rationals())))
	}

	dirRef, hasDirRef := asciiEntry(m, tagGPSImgDirectionRef)
	if hasDirRef {
		tags = append(tags, newEXIFTag("GPS Img Direction Ref", gpsDirectionRefName(dirRef)))
	}
	if e, ok := m[tagGPSImgDirection]; ok {
		if v, ok := e.rationalValue(); ok {
			tags = append(tags, newEXIFTag("GPS Img Direction", formatRational(v)))
		}
	}

	if v, ok := asciiEntry(m, tagGPSDateStamp); ok {
		tags = append(tags, newEXIFTag("GPS Date Stamp", v))
	}

	return tags
}

func newEXIFTag(name, value string) domain.Tag {
	return domain.NewTag(exifIFDPath, name, value)
}

func asciiEntry(entries map[uint16]tiffEntry, tag uint16) (string, bool) {
	e, ok := entries[tag]
	if !ok {
		return "", false
	}
	return e.ascii(), true
}

func gpsByteEntry(entries map[uint16]tiffEntry, tag uint16) (byte, bool) {
	e, ok := entries[tag]
	if !ok || len(e.raw) == 0 {
		return 0, false
	}
	return e.raw[0], true
}

func versionEntry(entries map[uint16]tiffEntry, tag uint16) (string, bool) {
	e, ok := entries[tag]
	if !ok {
		return "", false
	}
	return string(e.raw), true
}

func formatRational(r rational) string {
	if r.Den == 0 {
		return "0"
	}
	return formatFloatValue(float64(r.Num) / float64(r.Den))
}

func formatSRational(r sRational) string {
	if r.Den == 0 {
		return "0"
	}
	return formatFloatValue(float64(r.Num) / float64(r.Den))
}

func formatFloatValue(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func roundTo(v float64, places int) float64 {
	p := math.Pow(10, float64(places))
	return math.Round(v*p) / p
}

func formatExposureTime(r rational) string {
	if r.Num == 0 || r.Den == 0 {
		return "0"
	}
	value := float64(r.Num) / float64(r.Den)
	if value < 1 {
		return fmt.Sprintf("1/%d", int64(math.Round(float64(r.Den)/float64(r.Num))))
	}
	return formatFloatValue(value)
}

func formatShutterSpeedAPEX(r sRational) string {
	if r.Den == 0 {
		return "0"
	}
	seconds := math.Pow(2, -float64(r.Num)/float64(r.Den))
	if seconds > 0 && seconds < 1 {
		return fmt.Sprintf("1/%d", int64(math.Round(1/seconds)))
	}
	return formatFloatValue(roundTo(seconds, 1))
}

func formatApertureAPEX(r rational) string {
	if r.Den == 0 {
		return "0"
	}
	fNumber := math.Pow(2, float64(r.Num)/float64(r.Den)/2)
	return formatFloatValue(roundTo(fNumber, 1))
}

func formatSubjectDistance(r rational) string {
	if r.Den == 0 {
		return "0 m"
	}
	return formatFloatValue(float64(r.Num)/float64(r.Den)) + " m"
}

func formatGPSTimeStamp(r []rational) string {
	if len(r) < 3 {
		return ""
	}
	return fmt.Sprintf("%02d:%02d:%02d", ratioToInt(r[0]), ratioToInt(r[1]), ratioToInt(r[2]))
}

func formatGPSCoordinate(r []rational, ref string) string {
	if len(r) < 3 {
		return ""
	}
	deg := float64(r[0].Num) / float64(r[0].Den)
	min := float64(r[1].Num) / float64(r[1].Den)
	sec := float64(r[2].Num) / float64(r[2].Den)
	return fmt.Sprintf(`%s deg %s' %.2f" %s`, formatFloatValue(deg), formatFloatValue(min), sec, ref)
}

func formatGPSAltitude(r rational, ref byte) string {
	value := 0.0
	if r.Den != 0 {
		value = float64(r.Num) / float64(r.Den)
	}
	return fmt.Sprintf("%s m %s", formatFloatValue(value), gpsAltitudeRefName(ref))
}

func ratioToInt(r rational) int64 {
	if r.Den == 0 {
		return 0
	}
	return int64(math.Round(float64(r.Num) / float64(r.Den)))
}

func componentsConfigurationString(raw []byte) string {
	names := []string{"-", "Y", "Cb", "Cr", "R", "G", "B"}
	parts := make([]string, 0, len(raw))
	for _, v := range raw {
		if int(v) < len(names) {
			parts = append(parts, names[v])
		} else {
			parts = append(parts, fmt.Sprintf("Unknown (%d)", v))
		}
	}
	return strings.Join(parts, ", ")
}

func orientationString(value uint16) string {
	switch value {
	case 1:
		return "Horizontal (normal)"
	case 2:
		return "Mirror horizontal"
	case 3:
		return "Rotate 180"
	case 4:
		return "Mirror vertical"
	case 5:
		return "Mirror horizontal and rotate 270 CW"
	case 6:
		return "Rotate 90 CW"
	case 7:
		return "Mirror horizontal and rotate 90 CW"
	case 8:
		return "Rotate 270 CW"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func resolutionUnitString(value uint16) string {
	switch value {
	case 1:
		return "None"
	case 2:
		return "inches"
	case 3:
		return "cm"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func yCbCrPositioningString(value uint16) string {
	switch value {
	case 1:
		return "Centered"
	case 2:
		return "Co-sited"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func exposureProgramString(value uint16) string {
	switch value {
	case 0:
		return "Not Defined"
	case 1:
		return "Manual"
	case 2:
		return "Program AE"
	case 3:
		return "Aperture-priority AE"
	case 4:
		return "Shutter speed priority AE"
	case 5:
		return "Creative (Slow speed)"
	case 6:
		return "Action (High speed)"
	case 7:
		return "Portrait"
	case 8:
		return "Landscape"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func meteringModeString(value uint16) string {
	switch value {
	case 0:
		return "Unknown"
	case 1:
		return "Average"
	case 2:
		return "Center-weighted average"
	case 3:
		return "Spot"
	case 4:
		return "Multi-spot"
	case 5:
		return "Multi-segment"
	case 6:
		return "Partial"
	case 255:
		return "Other"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func flashString(value uint16) string {
	if s, ok := flashDescriptions[value]; ok {
		return s
	}
	return fmt.Sprintf("Unknown (%d)", value)
}

func sensingMethodString(value uint16) string {
	switch value {
	case 1:
		return "Not defined"
	case 2:
		return "One-chip color area"
	case 3:
		return "Two-chip color area"
	case 4:
		return "Three-chip color area"
	case 5:
		return "Color sequential area"
	case 7:
		return "Trilinear"
	case 8:
		return "Color sequential linear"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func sceneTypeString(value byte) string {
	if value == 1 {
		return "Directly photographed"
	}
	return fmt.Sprintf("Unknown (%d)", value)
}

func customRenderedString(value uint16) string {
	switch value {
	case 0:
		return "Normal"
	case 1:
		return "Custom"
	case 2:
		return "HDR (no original saved)"
	case 3:
		return "HDR (original saved)"
	case 4:
		return "Original (for HDR)"
	case 6:
		return "Panorama"
	case 7:
		return "Portrait Processing"
	case 8:
		return "Portrait Processing (no original saved)"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func exposureModeString(value uint16) string {
	switch value {
	case 0:
		return "Auto"
	case 1:
		return "Manual"
	case 2:
		return "Auto bracket"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func whiteBalanceString(value uint16) string {
	switch value {
	case 0:
		return "Auto"
	case 1:
		return "Manual"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func sceneCaptureTypeString(value uint16) string {
	switch value {
	case 0:
		return "Standard"
	case 1:
		return "Landscape"
	case 2:
		return "Portrait"
	case 3:
		return "Night"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func lowNormalHighString(value uint16) string {
	switch value {
	case 0:
		return "Normal"
	case 1:
		return "Low"
	case 2:
		return "High"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func sharpnessString(value uint16) string {
	switch value {
	case 0:
		return "Normal"
	case 1:
		return "Soft"
	case 2:
		return "Hard"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func subjectDistanceRangeString(value uint16) string {
	switch value {
	case 0:
		return "Unknown"
	case 1:
		return "Macro"
	case 2:
		return "Close"
	case 3:
		return "Distant"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func compositeImageString(value uint16) string {
	switch value {
	case 0:
		return "Unknown"
	case 1:
		return "Not a Composite Image"
	case 2:
		return "General Composite Image"
	case 3:
		return "Composite Image Captured While Shooting"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func gpsRefName(ref string) string {
	switch ref {
	case "N":
		return "North"
	case "S":
		return "South"
	case "E":
		return "East"
	case "W":
		return "West"
	default:
		return ref
	}
}

func gpsAltitudeRefName(ref byte) string {
	if ref == 1 {
		return "Below Sea Level"
	}
	return "Above Sea Level"
}

func gpsDirectionRefName(ref string) string {
	switch ref {
	case "T":
		return "True North"
	case "M":
		return "Magnetic North"
	default:
		return ref
	}
}
