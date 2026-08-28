// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   image.go                                           :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 05:58:25 by dande-je          #+#    #+#             //
//   Updated: 2026/08/27 18:08:50 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import (
	"bytes"
	"fmt"
)

type Format uint8

const (
	FormatUnknown Format = iota
	FormatJPEG
	FormatPNG
	FormatGIF
	FormatBMP
)

const maxSignatureLen = 9

var signatures = []signature{
	{[]byte{0xFF, 0xD8, 0xFF}, FormatJPEG},
	{[]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}, FormatPNG},
	{[]byte("GIF8"), FormatGIF},
	{[]byte("BM"), FormatBMP},
}

type signature struct {
	magic  []byte
	format Format
}

func (f Format) String() string {
	switch f {
	case FormatJPEG:
		return "JPEG"
	case FormatPNG:
		return "PNG"
	case FormatGIF:
		return "GIF"
	case FormatBMP:
		return "BMP"
	default:
		return "unknown"
	}
}

func DectectFormat(header []byte) (Format, error) {
	for _, sig := range signatures {
		if len(header) >= len(sig.magic) && bytes.Equal(header[:len(sig.magic)], sig.magic) {
			return sig.format, nil
		}
	}
	return FormatUnknown, fmt.Errorf("%w", ErrUnknownFormat)
}

func MaxSignatureLen() int {
	return maxSignatureLen
}
