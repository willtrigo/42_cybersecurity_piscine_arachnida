// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp_header.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/10 23:04:30 by dande-je          #+#    #+#             //
//   Updated: 2026/09/24 16:52:30 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const bmpMinDIBHeaderSize = 40

const (
	bmpSignatureSize  = 2
	bmpSignature      = "BM"
	bmpFileHeaderSize = 14
)

const (
	dibRedMaskOffset         = 40
	dibAlphaMaskOffset       = 52
	dibColorSpaceOffset      = 56
	dibRenderingIntentOffset = 108
)

const bitmapV2HeaderSize = 52

const bitmapV3HeaderSize = 56

const bitmapV4HeaderSize = 108

const bitmapV5HeaderSize = 124

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

type bmpReader struct {
	*byteCursor
}

func decodeBMPHeader(data []byte) (header bmpHeader, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("decode BMP header: %w", asError(r))
		}
	}()

	reader := &bmpReader{newByteCursor(data)}
	header = parseBMPHeader(reader)
	return header, nil
}

func asError(v any) error {
	if err, ok := v.(error); ok {
		return err
	}
	return fmt.Errorf("%v", v)
}

func parseBMPHeader(reader *bmpReader) bmpHeader {
	reader.readFileHeader()

	dibHeaderSize := reader.mustUint32LE()
	if dibHeaderSize < bmpMinDIBHeaderSize {
		panic(fmt.Errorf("unsupported DIB header size: %d", dibHeaderSize))
	}
	if reader.remaining() < int(dibHeaderSize) {
		panic(domain.ErrTruncatedBMPHeader)
	}

	dibStart := reader.pos - uint32Size
	dib := &bmpReader{newByteCursor(data(reader, dibStart, int(dibHeaderSize)))}
	dib.mustSkip(uint32Size)

	header := bmpHeader{DIBHeaderSize: dibHeaderSize}
	header.readMandatoryFields(dib)
	header.readOptionalFields(dib)

	return header
}

func (r *bmpReader) readFileHeader() {
	signature, err := r.readBytes(bmpSignatureSize)
	if err != nil {
		panic(domain.ErrTruncatedBMPFileHeader)
	}
	if string(signature) != bmpSignature {
		panic(domain.ErrInvalidSignature)
	}
	r.mustSkip(bmpFileHeaderSize - bmpSignatureSize)
}

func (r *bmpReader) mustSkip(n int) {
	if err := r.skip(n); err != nil {
		panic(err)
	}
}

func (r *bmpReader) mustUint32LE() uint32 {
	v, err := r.readUint32LE()
	if err != nil {
		panic(err)
	}
	return v
}

func data(reader *bmpReader, start, length int) []byte {
	return reader.data[start : start+length]
}

func (h *bmpHeader) readMandatoryFields(r *bmpReader) {
	h.Width = r.mustUint32LE()
	h.Height = r.mustUint32LE()
	h.Planes = r.mustUint16LE()
	h.BitCount = r.mustUint16LE()
	h.Compression = r.mustUint32LE()
	h.SizeImage = r.mustUint32LE()
	h.XPelsPerMeter = r.mustUint32LE()
	h.YPelsPerMeter = r.mustUint32LE()
	h.ClrUsed = r.mustUint32LE()
	h.ClrImportant = r.mustUint32LE()
}

func (r *bmpReader) mustUint16LE() uint16 {
	v, err := r.readUint16LE()
	if err != nil {
		panic(err)
	}
	return v
}

func (h *bmpHeader) readOptionalFields(r *bmpReader) {
	if h.hasBitfieldMasks() {
		r.mustSkipTo(dibRedMaskOffset)
		h.RedMask = r.mustUint32LE()
		h.GreenMask = r.mustUint32LE()
		h.BlueMask = r.mustUint32LE()
	}

	if h.hasAlphaMask() {
		r.mustSkipTo(dibAlphaMaskOffset)
		h.AlphaMask = r.mustUint32LE()
	}

	if h.hasColorSpace() {
		r.mustSkipTo(dibColorSpaceOffset)
		h.ColorSpace = r.mustUint32LE()
	}

	if h.hasRenderingIntent() {
		r.mustSkipTo(dibRenderingIntentOffset)
		h.RenderingIntent = r.mustUint32LE()
	}
}

func (r *bmpReader) mustSkipTo(offset int) {
	if r.pos < offset {
		r.mustSkip(offset - r.pos)
	}
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
