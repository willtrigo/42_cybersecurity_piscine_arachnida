// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   bmp.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 18:23:57 by dande-je          #+#    #+#             //
//   Updated: 2026/08/28 10:23:26 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	bmpFileHeaderSize = 14

	bmpDIBSizeFieldSize = 4

	bmpCoreHeaderSize = 12

	bmpCoreWidthOffsetBegin    = 0
	bmpCoreWidthOffsetEnd      = bmpCoreWidthOffsetBegin + 2
	bmpCoreHeightOffsetBegin   = 2
	bmpCoreHeightOffsetEnd     = bmpCoreHeightOffsetBegin + 2
	bmpCorePlanesOffsetBegin   = 4
	bmpCorePlanesOffsetEnd     = bmpCorePlanesOffsetBegin + 2
	bmpCoreBitCountOffsetBegin = 6
	bmpCoreBitCountOffsetEnd   = bmpCoreBitCountOffsetBegin + 2

	bmpInfoWidthOffsetBegin    = 0
	bmpInfoWidthOffsetEnd      = bmpInfoWidthOffsetBegin + 4
	bmpInfoHeightOffsetBegin   = 4
	bmpInfoHeightOffsetEnd     = bmpInfoHeightOffsetBegin + 4
	bmpInfoPlanesOffsetBegin   = 8
	bmpInfoPlanesOffsetEnd     = bmpInfoPlanesOffsetBegin + 2
	bmpInfoBitCountOffsetBegin = 10
	bmpInfoBitCountOffsetEnd   = bmpInfoBitCountOffsetBegin + 2
)

type BMPParser struct{}

func NewBMPParser() *BMPParser {
	return &BMPParser{}
}

func (BMPParser) Read(path string) (metadata *domain.Metadata, err error) {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("bmp: invalid file path")
	}

	data, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("bmp: %w", err)
	}
	defer func() {
		err = errors.Join(err, data.Close())
	}()

	buf := make([]byte, bmpFileHeaderSize+bmpDIBSizeFieldSize)
	if _, err := io.ReadFull(data, buf); err != nil {
		return nil, fmt.Errorf("bmp: reading header: %w", err)
	}

	dibSize := binary.LittleEndian.Uint32(buf[bmpFileHeaderSize:])

	dib := make([]byte, dibSize-bmpDIBSizeFieldSize)
	if _, err := io.ReadFull(data, dib); err != nil {
		return nil, fmt.Errorf("bmp: reading DIB header: %w", err)
	}

	var width, height int
	var planes, bitCount uint16

	if dibSize == bmpCoreHeaderSize {
		width = int(binary.LittleEndian.Uint16(dib[bmpCoreWidthOffsetBegin:bmpCoreWidthOffsetEnd]))
		height = int(binary.LittleEndian.Uint16(dib[bmpCoreHeightOffsetBegin:bmpCoreHeightOffsetEnd]))
		planes = binary.LittleEndian.Uint16(dib[bmpCorePlanesOffsetBegin:bmpCorePlanesOffsetEnd])
		bitCount = binary.LittleEndian.Uint16(dib[bmpCoreBitCountOffsetBegin:bmpCoreBitCountOffsetEnd])
	} else {
		// Width/height are stored as signed LONG per the BITMAPINFOHEADER spec;
		// a negative height indicates a top-down bitmap. Reinterpreting the raw
		// bits as int32 is intentional, not an overflow.
		width = int(int32(binary.LittleEndian.Uint32(dib[bmpInfoWidthOffsetBegin:bmpInfoWidthOffsetEnd])))    // #nosec G115
		height = int(int32(binary.LittleEndian.Uint32(dib[bmpInfoHeightOffsetBegin:bmpInfoHeightOffsetEnd]))) // #nosec G115
		planes = binary.LittleEndian.Uint16(dib[bmpInfoPlanesOffsetBegin:bmpInfoPlanesOffsetEnd])
		bitCount = binary.LittleEndian.Uint16(dib[bmpInfoBitCountOffsetBegin:bmpInfoBitCountOffsetEnd])
	}

	return &domain.Metadata{
		Format:     domain.FormatBMP,
		Dimensions: domain.Dimensions{Width: width, Height: height},
		Tags: []domain.Tag{
			{Name: "BitsPerPixel", Value: strconv.Itoa(int(bitCount))},
			{Name: "ColorPlanes", Value: strconv.Itoa(int(planes))},
			{Name: "DIBHeaderSize", Value: strconv.Itoa(int(dibSize))},
		},
	}, nil
}
