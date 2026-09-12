// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp_header.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/10 23:04:30 by dande-je          #+#    #+#             //
//   Updated: 2026/09/12 22:55:44 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"encoding/binary"
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const bmpIFDpath = "BMP"

const (
	bmpFileHeaderSize = 14

	bmpSignatureOffset = 0
	bmpSignatureSize   = 2
	bmpSignature       = "BM"

	bmpDIBHeaderSizeOffset = 14
	bmpMinDIBHeaderSize    = 40

	dibWidthOffset         = 4
	dibHeightOffset        = 8
	dibPlanesOffset        = 12
	dibBitCountOffset      = 14
	dibCompressionOffset   = 16
	dibSizeImageOffset     = 20
	dibXPelsPerMeterOffset = 24
	dibYPelsPerMeterOffset = 28
	dibClrUsedOffset       = 32
	dibClrImportantOffset  = 36

	bitmapCoreHeaderSize = 12
	bitmapInfoHeaderSize = 40

	bitmapV2HeaderSize = 52
	dibRedMaskOffset   = 40
	dibGreenMaskOffset = 44
	dibBlueMaskOffset  = 48

	bitmapV3HeaderSize = 56
	dibAlphaMaskOffset = 52

	bitmapV4HeaderSize  = 108
	dibColorSpaceOffset = 56

	bitmapV5HeaderSize       = 124
	dibRenderingIntentOffset = 108
)

const (
	biRGB            = 0
	biRLE8           = 1
	biRLE4           = 2
	biBitfields      = 3
	biJPEG           = 4
	biPNG            = 5
	biAlphaBitfields = 6
	biCMYK           = 11
	biCMYKRLE8       = 12
	biCMYKRLE4       = 13
)

const (
	lcsCalibratedRGB = 0x00000000
	lcsSRGB          = 0x73524742
	lcsWindows       = 0x57696e20
	lcsProfileLinked = 0x4c494e4b
	lcsProfileEmbed  = 0x4d424544
)

const (
	lcsGMBusiness        = 1
	lcsGMGraphics        = 2
	lcsGMImages          = 4
	lcsGMAbsColorimetric = 8
)

const pixelsPerMegapixel = 1_000_000

type bmpHeader struct {
	Width           uint32
	Height          uint32
	XPelsPerMeter   uint32
	YPelsPerMeter   uint32
	Planes          uint16
	BitCount        uint16
	Compression     uint32
	SizeImage       uint32
	ClrUsed         uint32
	ClrImportant    uint32
	RedMask         uint32
	GreenMask       uint32
	BlueMask        uint32
	AlphaMask       uint32
	ColorSpace      uint32
	RenderingIntent uint32
	DIBHeaderSize   uint32
}

func decodeBMPHeader(data []byte) (bmpHeader, error) {
	if len(data) < bmpFileHeaderSize+uint32Size {
		return bmpHeader{}, fmt.Errorf("truncated file header")
	}
	if string(data[bmpSignatureOffset:bmpSignatureOffset+bmpSignatureSize]) != bmpSignature {
		return bmpHeader{}, fmt.Errorf("invalid signature")
	}

	dibHeaderSize := readUint32At(data, bmpDIBHeaderSizeOffset)
	if dibHeaderSize < bmpMinDIBHeaderSize {
		return bmpHeader{}, fmt.Errorf("unsupported DIB header size: %d", dibHeaderSize)
	}

	end := bmpFileHeaderSize + int(dibHeaderSize)
	if len(data) < end {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	dib := data[bmpFileHeaderSize:end]

	header := bmpHeader{
		DIBHeaderSize: dibHeaderSize,
		Width:         readUint32At(dib, dibWidthOffset),
		Height:        readUint32At(dib, dibHeightOffset),
		Planes:        readUint16At(dib, dibPlanesOffset),
		BitCount:      readUint16At(dib, dibBitCountOffset),
		Compression:   readUint32At(dib, dibCompressionOffset),
		SizeImage:     readUint32At(dib, dibSizeImageOffset),
		XPelsPerMeter: readUint32At(dib, dibXPelsPerMeterOffset),
		YPelsPerMeter: readUint32At(dib, dibYPelsPerMeterOffset),
		ClrUsed:       readUint32At(dib, dibClrUsedOffset),
		ClrImportant:  readUint32At(dib, dibClrImportantOffset),
	}

	if header.hasBitfieldMasks() {
		header.RedMask = readUint32At(dib, dibRedMaskOffset)
		header.GreenMask = readUint32At(dib, dibGreenMaskOffset)
		header.BlueMask = readUint32At(dib, dibBlueMaskOffset)
	}
	if header.hasAlphaMask() {
		header.AlphaMask = readUint32At(dib, dibAlphaMaskOffset)
	}
	if header.hasColorSpace() {
		header.ColorSpace = readUint32At(dib, dibColorSpaceOffset)
	}
	if header.hasRenderingIntent() {
		header.RenderingIntent = readUint32At(dib, dibRenderingIntentOffset)
	}

	return header, nil
}

func readUint32At(b []byte, offset int) uint32 {
	return binary.LittleEndian.Uint32(b[offset : offset+uint32Size])
}

func readUint16At(b []byte, offset int) uint16 {
	return binary.LittleEndian.Uint16(b[offset : offset+uint16Size])
}

func (h bmpHeader) hasBitfieldMasks() bool {
	return h.dibAtLeast(bitmapV2HeaderSize)
}

func (h bmpHeader) hasAlphaMask() bool {
	return h.dibAtLeast(bitmapV3HeaderSize)
}

func (h bmpHeader) hasColorSpace() bool {
	return h.dibAtLeast(bitmapV4HeaderSize)
}

func (h bmpHeader) hasRenderingIntent() bool {
	return h.dibAtLeast(bitmapV5HeaderSize)
}

func (h bmpHeader) dibAtLeast(size uint32) bool {
	return h.DIBHeaderSize >= size
}

func (h bmpHeader) tags() []domain.Tag {
	tags := []domain.Tag{
		newBMPTag("File Type", "BMP"),
		newBMPTag("File Type Extension", "bmp"),
		newBMPTag("MIME Type", "image/bmp"),
		newBMPTag("BMP Version", bmpVersion(h.DIBHeaderSize)),
		newBMPTag("Image Width", fmt.Sprintf("%d", h.Width)),
		newBMPTag("Image Height", fmt.Sprintf("%d", h.Height)),
		newBMPTag("Planes", fmt.Sprintf("%d", h.Planes)),
		newBMPTag("Bit Depth", fmt.Sprintf("%d", h.BitCount)),
		newBMPTag("Compression", bmpCompression(h.Compression)),
		newBMPTag("Image Length", fmt.Sprintf("%d", h.SizeImage)),
		newBMPTag("Pixels Per Meter X", fmt.Sprintf("%d", h.XPelsPerMeter)),
		newBMPTag("Pixels Per Meter Y", fmt.Sprintf("%d", h.YPelsPerMeter)),
		newBMPTag("Num Colors", bmpColorCount(h.ClrUsed, "Use Bitdepth")),
		newBMPTag("Num Important Colors", bmpColorCount(h.ClrImportant, "All")),
	}

	if h.hasBitfieldMasks() {
		tags = append(tags,
			newBMPTag("Red Mask", bmpMask(h.RedMask)),
			newBMPTag("Green Mask", bmpMask(h.GreenMask)),
			newBMPTag("Blue Mask", bmpMask(h.BlueMask)),
		)
	}
	if h.hasAlphaMask() {
		tags = append(tags, newBMPTag("Alpha Mask", bmpMask(h.AlphaMask)))
	}
	if h.hasColorSpace() {
		tags = append(tags, newBMPTag("Color Space", bmpColorSpace(h.ColorSpace)))
	}
	if h.hasRenderingIntent() {
		tags = append(tags, newBMPTag("Rendering Intent", bmlRenderingIntent(h.RenderingIntent)))
	}

	tags = append(tags,
		newBMPTag("Image Size", fmt.Sprintf("%d x %d", h.Width, h.Height)),
		newBMPTag("Megapixels", fmt.Sprintf("%.1f", megapixels(h.Width, h.Height))),
	)

	return tags
}

func newBMPTag(name, value string) domain.Tag {
	return domain.Tag{IDFPath: bmpIFDpath, Name: name, Value: value}
}

func bmpVersion(dibHeaderSize uint32) string {
	switch dibHeaderSize {
	case bitmapCoreHeaderSize:
		return "Windows V2"
	case bitmapInfoHeaderSize:
		return "Windows V3"
	case bitmapV2HeaderSize, bitmapV3HeaderSize:
		return "Windows V3 (Adobe)"
	case bitmapV4HeaderSize:
		return "Windows V4"
	case bitmapV5HeaderSize:
		return "Windows V5"
	default:
		return fmt.Sprintf("Unknown (%d)", dibHeaderSize)
	}
}

func bmpCompression(value uint32) string {
	switch value {
	case biRGB:
		return "None"
	case biRLE8:
		return "8-bit RLE"
	case biRLE4:
		return "4-bit RLE"
	case biBitfields:
		return "Bitfields"
	case biJPEG:
		return "JPEG"
	case biPNG:
		return "PNG"
	case biAlphaBitfields:
		return "Alpha Bitfields"
	case biCMYK:
		return "CMYK"
	case biCMYKRLE8:
		return "CMYK 8-bit RLE"
	case biCMYKRLE4:
		return "CMYK 4-bit RLE"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func bmpColorCount(value uint32, zeroLabel string) string {
	if value == 0 {
		return zeroLabel
	}
	return fmt.Sprintf("%d", value)
}

func bmpMask(value uint32) string {
	return fmt.Sprintf("0x%08x", value)
}

func bmpColorSpace(value uint32) string {
	switch value {
	case lcsCalibratedRGB:
		return "Calibrated RGB"
	case lcsSRGB:
		return "sRGB"
	case lcsWindows:
		return "Windows Color Space"
	case lcsProfileLinked:
		return "Linked"
	case lcsProfileEmbed:
		return "Embedded"
	default:
		return fmt.Sprintf("Unknown (0x%08x)", value)
	}
}

func bmlRenderingIntent(value uint32) string {
	switch value {
	case lcsGMBusiness:
		return "Saturation (LGS_GM_BUSINESS)"
	case lcsGMGraphics:
		return "Graphic (LCS_GM_GRAPHICS)"
	case lcsGMImages:
		return "Picture (LCS_GM_IMAGES)"
	case lcsGMAbsColorimetric:
		return "Absolute (LCS_GM_ABS_COLORIMETRIC)"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func megapixels(width, height uint32) float64 {
	return float64(int64(width)*int64(height)) / pixelsPerMegapixel
}
