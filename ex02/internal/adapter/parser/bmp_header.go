// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp_header.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/10 23:04:30 by dande-je          #+#    #+#             //
//   Updated: 2026/09/23 10:32:12 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import "fmt"

const (
	bmpFileHeaderSize = 14

	bmpSignatureOffset = 0
	bmpSignatureSize   = 2
	bmpSignature       = "BM"

	bmpDIBHeaderSizeOffset = 14
	bmpMinDIBHeaderSize    = 40

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
	var err error

	cursor := newByteCursor(data)

	signature, err := cursor.readBytes(bmpSignatureSize)
	if err != nil {
		return bmpHeader{}, fmt.Errorf("truncated file header")
	}
	if string(signature) != bmpSignature {
		return bmpHeader{}, fmt.Errorf("invalid signature")
	}
	if err = cursor.skip(bmpFileHeaderSize - bmpSignatureSize); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated file header")
	}

	dibHeaderSize, err := cursor.readUint32LE()
	if err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header size")
	}
	if dibHeaderSize < bmpMinDIBHeaderSize {
		return bmpHeader{}, fmt.Errorf("unsupported DIB header size: %d", dibHeaderSize)
	}
	if cursor.remaining() < int(dibHeaderSize) {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}

	dibStart := cursor.pos - uint32Size
	dibCursor := newByteCursor(data[dibStart : dibStart+int(dibHeaderSize)])
	if _, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}

	header := bmpHeader{
		DIBHeaderSize: dibHeaderSize,
	}
	if header.Width, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.Height, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.Planes, err = dibCursor.readUint16LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.BitCount, err = dibCursor.readUint16LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.Compression, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.SizeImage, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.XPelsPerMeter, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.YPelsPerMeter, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.ClrUsed, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}
	if header.ClrImportant, err = dibCursor.readUint32LE(); err != nil {
		return bmpHeader{}, fmt.Errorf("truncated DIB header")
	}

	if header.hasBitfieldMasks() {
		if err = skipTo(dibCursor, dibRedMaskOffset); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
		if header.RedMask, err = dibCursor.readUint32LE(); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
		if header.GreenMask, err = dibCursor.readUint32LE(); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
		if header.BlueMask, err = dibCursor.readUint32LE(); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
	}

	if header.hasAlphaMask() {
		if err = skipTo(dibCursor, dibAlphaMaskOffset); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
		if header.AlphaMask, err = dibCursor.readUint32LE(); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
	}

	if header.hasColorSpace() {
		if err = skipTo(dibCursor, dibColorSpaceOffset); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
		if header.ColorSpace, err = dibCursor.readUint32LE(); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
	}

	if header.hasRenderingIntent() {
		if err = skipTo(dibCursor, dibRenderingIntentOffset); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
		if header.RenderingIntent, err = dibCursor.readUint32LE(); err != nil {
			return bmpHeader{}, fmt.Errorf("truncated DIB header")
		}
	}

	return header, nil
}

func skipTo(cursor *byteCursor, offset int) error {
	if cursor.pos >= offset {
		return nil
	}
	return cursor.skip(offset - cursor.pos)
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
