// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   jpeg_header.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/22 21:08:10 by dande-je          #+#    #+#             //
//   Updated: 2026/09/22 22:40:37 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	jpegMarkerSize     = 2
	jpegSegmentLenSize = 2

	jpegSOIMarker = 0xFFD8
	jpegEOIMarker = 0xFFD9
	jpegSOSMarker = 0xFFDA

	jpegAPP0Marker = 0xFFE0
	jpegAPP1Marker = 0xFFE1
	jpegAPP2Marker = 0xFFE2
	jpegCOMMarker  = 0xFFFE

	jfifSignature = "JFIF\x00"
	xmpSignature  = "http://ns.adobe.com/xap/1.0/\x00"
	iccSignature  = "ICC_PROFILE\x00"

	sofMarkerBase = 0xC0
	sofMarkerEnd  = 0xCF
	sofMarkerDHT  = 0xC4
	sofMarkerJPG  = 0xC8
	sofMarkerDAC  = 0xCC
)

type jpegComponent struct {
	ID         byte
	HSampling  byte
	VSampling  byte
	QuantTable byte
}

type jpegHeader struct {
	Components      []jpegComponent
	EncodingProcess string
	ICCProfileName  string
	Comment         string
	XMPPacket       []byte
	ICCProfile      []byte
	EXIFData        []byte
	Width           uint16
	Height          uint16
	NumComponents   uint8
	BitsPerSample   uint8
	IsProgressive   bool
	HasJFIF         bool
}

func decodeJPEGHeader(data []byte) (jpegHeader, error) {
	if len(data) < jpegMarkerSize {
		return jpegHeader{}, fmt.Errorf("truncated JPEG")
	}
	if binary.BigEndian.Uint16(data[0:2]) != jpegSOIMarker {
		return jpegHeader{}, fmt.Errorf("invalid JPEG SOI marker")
	}

	var header jpegHeader
	pos := jpegMarkerSize
	sawSOF := false

	for pos+jpegMarkerSize <= len(data) {
		if data[pos] != 0xFF {
			return jpegHeader{}, fmt.Errorf("invalid JPEG marker at offset %d", pos)
		}
		marker := binary.BigEndian.Uint16(data[pos : pos+2])
		pos += jpegMarkerSize

		if marker == jpegEOIMarker || marker == jpegSOSMarker {
			break
		}
		if marker >= 0xFFD0 && marker <= 0xFFD7 {
			continue
		}
		if marker == 0xFF01 {
			continue
		}

		if pos+jpegSegmentLenSize > len(data) {
			return jpegHeader{}, fmt.Errorf("truncated JPEG segment length")
		}
		segLen := int(binary.BigEndian.Uint16(data[pos : pos+2]))
		if segLen < jpegSegmentLenSize || pos+segLen > len(data) {
			return jpegHeader{}, fmt.Errorf("invalid JPEG segment length")
		}
		payload := data[pos+jpegSegmentLenSize : pos+segLen]
		pos += segLen

		switch marker {
		case jpegAPP0Marker:
			if bytes.HasPrefix(payload, []byte(jfifSignature)) {
				header.HasJFIF = true
			}
		case jpegAPP1Marker:
			applyAPP1(payload, &header)
		case jpegAPP2Marker:
			applyAPP2(payload, &header)
		case jpegCOMMarker:
			header.Comment = string(payload)
		default:
			if isSOFMarker(marker) {
				if err := decodeSOFSegment(payload, marker, &header); err != nil {
					return jpegHeader{}, err
				}
				sawSOF = true
			}
		}
	}

	if !sawSOF {
		return jpegHeader{}, fmt.Errorf("missing JPEG SOF marker")
	}

	return header, nil
}

func isSOFMarker(marker uint16) bool {
	b := byte(marker & 0xFF)
	if b < sofMarkerBase || b > sofMarkerEnd {
		return false
	}
	if b == sofMarkerDHT || b == sofMarkerJPG || b == sofMarkerDAC {
		return false
	}
	return true
}

func applyAPP1(payload []byte, header *jpegHeader) {
	switch {
	case bytes.HasPrefix(payload, []byte(exifBlobPrefix)):
		header.EXIFData = append([]byte(nil), payload...)
	case bytes.HasPrefix(payload, []byte(xmpSignature)):
		header.XMPPacket = append([]byte(nil), payload[len(xmpSignature):]...)
	}
}

func applyAPP2(payload []byte, header *jpegHeader) {
	if !bytes.HasPrefix(payload, []byte(iccSignature)) {
		return
	}
	rest := payload[len(iccSignature):]
	if len(rest) < 2 {
		return
	}
	header.ICCProfile = append(header.ICCProfile, rest[2:]...)
}

func decodeSOFSegment(payload []byte, marker uint16, header *jpegHeader) error {
	if len(payload) < 6 {
		return fmt.Errorf("truncated SOF segment")
	}
	header.BitsPerSample = payload[0]
	header.Height = binary.BigEndian.Uint16(payload[1:3])
	header.Width = binary.BigEndian.Uint16(payload[3:5])
	header.NumComponents = payload[5]
	header.EncodingProcess = sofEncodingProcess(byte(marker & 0xFF))
	header.IsProgressive = byte(marker&0xFF) == 0xC2

	if len(payload) < 6+int(header.NumComponents)*3 {
		return fmt.Errorf("truncated SOF component list")
	}
	header.Components = make([]jpegComponent, 0, header.NumComponents)
	for i := 0; i < int(header.NumComponents); i++ {
		base := 6 + i*3
		header.Components = append(header.Components, jpegComponent{
			ID:         payload[base],
			HSampling:  payload[base+1] >> 4,
			VSampling:  payload[base+1] & 0x0F,
			QuantTable: payload[base+2],
		})
	}
	return nil
}

func sofEncodingProcess(b byte) string {
	switch b {
	case 0xC0:
		return "Baseline DCT, Huffman coding"
	case 0xC1:
		return "Extended sequential DCT, Huffman coding"
	case 0xC2:
		return "Progressive DCT, Huffman coding"
	case 0xC3:
		return "Lossless, Huffman coding"
	case 0xC5:
		return "Sequential DCT, differential Huffman coding"
	case 0xC6:
		return "Progressive DCT, differential Huffman coding"
	case 0xC7:
		return "Lossless, differential Huffman coding"
	case 0xC9:
		return "Extended sequential DCT, arithmetic coding"
	case 0xCA:
		return "Progressive DCT, arithmetic coding"
	case 0xCB:
		return "Lossless, arithmetic coding"
	case 0xCD:
		return "Sequential DCT, differential arithmetic coding"
	case 0xCE:
		return "Progressive DCT, differential arithmetic coding"
	case 0xCF:
		return "Lossless, differential arithmetic coding"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", b)
	}
}
