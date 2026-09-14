// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp_header.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/10 23:04:30 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 21:06:18 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"encoding/binary"
	"fmt"
)

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
