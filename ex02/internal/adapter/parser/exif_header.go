// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   exif_header.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/22 18:03:43 by dande-je          #+#    #+#             //
//   Updated: 2026/09/26 19:27:38 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	exifBlobPrefix = "Exif\x00\x00"

	tiffByteOrderLE = "II"
	tiffByteOrderBE = "MM"
	tiffMagicNumber = 42

	tiffHeaderSize    = 8
	ifdEntrySize      = 12
	ifdEntryCountSize = 2
	ifdValueFieldSize = 4

	exifIFDPointerTag = 0x8769
	gpsIFDPointerTag  = 0x8825
)

type tiffType uint16

const (
	tiffTypeByte      tiffType = 1
	tiffTypeASCII     tiffType = 2
	tiffTypeShort     tiffType = 3
	tiffTypeLong      tiffType = 4
	tiffTypeRational  tiffType = 5
	tiffTypeSByte     tiffType = 6
	tiffTypeUndefined tiffType = 7
	tiffTypeSShort    tiffType = 8
	tiffTypeSLong     tiffType = 9
	tiffTypeSRational tiffType = 10
	tiffTypeFloat     tiffType = 11
	tiffTypeDouble    tiffType = 12
)

var tiffTypeSize = map[tiffType]int{
	tiffTypeByte:      1,
	tiffTypeASCII:     1,
	tiffTypeShort:     2,
	tiffTypeLong:      4,
	tiffTypeRational:  8,
	tiffTypeSByte:     1,
	tiffTypeUndefined: 1,
	tiffTypeSShort:    2,
	tiffTypeSLong:     4,
	tiffTypeSRational: 8,
	tiffTypeFloat:     4,
	tiffTypeDouble:    8,
}

type rational struct {
	Num uint32
	Den uint32
}

type sRational struct {
	Num int32
	Den int32
}

type tiffEntry struct {
	order binary.ByteOrder
	raw   []byte
	Type  tiffType
	Count uint32
}

type exifHeader struct {
	IFD0      map[uint16]tiffEntry
	ExifIFD   map[uint16]tiffEntry
	GPSIFD    map[uint16]tiffEntry
	IFD1      map[uint16]tiffEntry
	Thumbnail []byte
	order     binary.ByteOrder
	data      []byte
}

func decodeEXIFData(data []byte) (exifHeader, error) {
	if strings.HasPrefix(string(data), exifBlobPrefix) {
		data = data[len(exifBlobPrefix):]
	}
	if len(data) < tiffHeaderSize {
		return exifHeader{}, fmt.Errorf("exif: truncated TIFF header")
	}

	var order binary.ByteOrder
	switch string(data[0:2]) {
	case tiffByteOrderLE:
		order = binary.LittleEndian
	case tiffByteOrderBE:
		order = binary.BigEndian
	default:
		return exifHeader{}, fmt.Errorf("exif: invalid byte order marker")
	}

	view := newTIFFView(data, order)
	if magic, _ := view.u16(2); magic != tiffMagicNumber {
		return exifHeader{}, fmt.Errorf("exif: invalid TIFF magic number")
	}

	ifd0Offset, ok := view.u32(4)
	if !ok {
		return exifHeader{}, fmt.Errorf("exif: truncated TIFF header")
	}

	ifd0, nextIFD, err := readIFD(view, ifd0Offset)
	if err != nil {
		return exifHeader{}, fmt.Errorf("exif: %w", err)
	}

	header := exifHeader{IFD0: ifd0, order: order, data: data}

	if off, ok := pointerOffset(ifd0, exifIFDPointerTag); ok {
		if sub, _, err := readIFD(view, off); err == nil {
			header.ExifIFD = sub
		}
	}
	if off, ok := pointerOffset(ifd0, gpsIFDPointerTag); ok {
		if gps, _, err := readIFD(view, off); err == nil {
			header.GPSIFD = gps
		}
	}

	if nextIFD != 0 {
		if ifd1, _, err := readIFD(view, nextIFD); err == nil {
			header.IFD1 = ifd1
			header.Thumbnail = extractThumbnail(view, ifd1)
		}
	}

	return header, nil
}

func pointerOffset(entries map[uint16]tiffEntry, tag uint16) (uint32, bool) {
	e, ok := entries[tag]
	if !ok {
		return 0, false
	}
	return e.longValue()
}

func extractThumbnail(view tiffView, ifd1 map[uint16]tiffEntry) []byte {
	offEntry, ok := ifd1[tagThumbnailOffset]
	if !ok {
		return nil
	}
	lenEntry, ok := ifd1[tagThumbnailLength]
	if !ok {
		return nil
	}

	off, ok := offEntry.longValue()
	if !ok {
		return nil
	}
	length, ok := lenEntry.longValue()
	if !ok {
		return nil
	}

	thumb, ok := view.copySlice(int(off), int(length))
	if !ok {
		return nil
	}
	return thumb
}

func readIFD(view tiffView, offset uint32) (map[uint16]tiffEntry, uint32, error) {
	start := int(offset)
	count, ok := view.u16(start)
	if !ok {
		return nil, 0, fmt.Errorf("truncated IFD")
	}

	entries := make(map[uint16]tiffEntry, count)
	pos := start + ifdEntryCountSize

	for i := 0; i < int(count); i++ {
		raw, ok := view.slice(pos, ifdEntrySize)
		if !ok {
			return nil, 0, fmt.Errorf("truncated IFD entry")
		}

		tag := view.order.Uint16(raw[0:2])
		typ := tiffType(view.order.Uint16(raw[2:4]))
		valueCount := view.order.Uint32(raw[4:8])
		valueField := raw[8 : 8+ifdValueFieldSize]

		if size, known := tiffTypeSize[typ]; known {
			if value, ok := resolveEntryValue(view, valueField, size, int(valueCount)); ok {
				entries[tag] = tiffEntry{Type: typ, Count: valueCount, raw: value, order: view.order}
			}
		}
		pos += ifdEntrySize
	}

	next, _ := view.u32(pos)
	return entries, next, nil
}

func resolveEntryValue(view tiffView, valueField []byte, elemSize, count int) ([]byte, bool) {
	total := elemSize * count
	if total <= ifdValueFieldSize {
		return append([]byte(nil), valueField[:total]...), true
	}

	offset := int(view.order.Uint32(valueField))
	return view.copySlice(offset, total)
}

func (e tiffEntry) ascii() string {
	return strings.TrimRight(string(e.raw), "\x00")
}

func (e tiffEntry) shorts() []uint16 {
	out := make([]uint16, len(e.raw)/2)
	for i := range out {
		out[i] = e.order.Uint16(e.raw[i*2 : i*2+2])
	}
	return out
}

func (e tiffEntry) longs() []uint32 {
	out := make([]uint32, len(e.raw)/4)
	for i := range out {
		out[i] = e.order.Uint32(e.raw[i*4 : i*4+4])
	}
	return out
}

func (e tiffEntry) rationals() []rational {
	out := make([]rational, len(e.raw)/8)
	for i := range out {
		out[i] = rational{
			Num: e.order.Uint32(e.raw[i*8 : i*8+4]),
			Den: e.order.Uint32(e.raw[i*8+4 : i*8+8]),
		}
	}
	return out
}

func (e tiffEntry) sRationals() []sRational {
	out := make([]sRational, len(e.raw)/8)
	for i := range out {
		out[i] = sRational{
			Num: int32(e.order.Uint32(e.raw[i*8 : i*8+4])),   // #nosec G115
			Den: int32(e.order.Uint32(e.raw[i*8+4 : i*8+8])), // #nosec G115
		}
	}
	return out
}

func (e tiffEntry) shortValue() (uint16, bool) {
	s := e.shorts()
	if len(s) == 0 {
		return 0, false
	}
	return s[0], true
}

func (e tiffEntry) longValue() (uint32, bool) {
	l := e.longs()
	if len(l) == 0 {
		return 0, false
	}
	return l[0], true
}

func (e tiffEntry) numericValue() (uint32, bool) {
	switch e.Type {
	case tiffTypeShort:
		v, ok := e.shortValue()
		return uint32(v), ok
	case tiffTypeLong:
		return e.longValue()
	default:
		return 0, false
	}
}

func (e tiffEntry) rationalValue() (rational, bool) {
	r := e.rationals()
	if len(r) == 0 {
		return rational{}, false
	}
	return r[0], true
}

func (e tiffEntry) sRationalValue() (sRational, bool) {
	r := e.sRationals()
	if len(r) == 0 {
		return sRational{}, false
	}
	return r[0], true
}

func (h exifHeader) empty() bool { return h.IFD0 == nil }
