// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   jpeg_header.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/22 21:08:10 by dande-je          #+#    #+#             //
//   Updated: 2026/09/26 18:23:09 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bytes"
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

type jpegReader struct {
	*byteCursor
}

func decodeJPEGHeader(data []byte) (jpegHeader, error) {
	return decodeWithRecover("JPEG", func() jpegHeader {
		return parseJpegHeader(&jpegReader{newByteCursor(data)})
	})
}

func parseJpegHeader(reader *jpegReader) jpegHeader {
	signature := reader.mustUint16BE()
	if signature != jpegSOIMarker {
		panic(fmt.Errorf("invalid JPEG SOI marker"))
	}

	var header jpegHeader
	sawSOF := false

	for reader.remaining() >= jpegMarkerSize {
		if marker, ok := reader.readMarker(); ok {
			if marker == jpegEOIMarker || marker == jpegSOSMarker {
				break
			}
			if isStandaloneMarker(marker) {
				continue
			}

			payload := reader.mustReadSegment()
			applySegment(marker, payload, &header, &sawSOF)
		}
	}

	if !sawSOF {
		panic(fmt.Errorf("missing JPEG SOF marker"))
	}

	return header
}

func (r *jpegReader) mustUint16BE() uint16 {
	v, err := r.readUint16BE()
	if err != nil {
		panic(fmt.Errorf("truncated JPEG"))
	}
	return v
}

func (r *jpegReader) readMarker() (uint16, bool) {
	marker, err := r.readUint16BE()
	if err != nil {
		return 0, false
	}
	return marker, true
}

func isStandaloneMarker(marker uint16) bool {
	return (marker >= 0xFFD0 && marker <= 0xFFD7) || marker == 0xFF01
}

func (r *jpegReader) mustReadSegment() []byte {
	segLen, err := r.readUint16BE()
	if err != nil || int(segLen) < jpegSegmentLenSize {
		panic(fmt.Errorf("invalid JPEG segment length"))
	}
	payload, err := r.readBytes(int(segLen) - jpegSegmentLenSize)
	if err != nil {
		panic(fmt.Errorf("truncated JPEG segment"))
	}
	return payload
}

func applySegment(marker uint16, payload []byte, header *jpegHeader, sawSOF *bool) {
	switch marker {
	case jpegAPP0Marker:
		if bytes.HasPrefix(payload, []byte(jfifSignature)) {
			header.HasJFIF = true
		}
	case jpegAPP1Marker:
		applyAPP1(payload, header)
	case jpegAPP2Marker:
		applyAPP2(payload, header)
	case jpegCOMMarker:
		header.Comment = string(payload)
	default:
		if isSOFMarker(marker) {
			if err := decodeSOFSegment(payload, marker, header); err != nil {
				panic(err)
			}
			*sawSOF = true
		}
	}
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
	reader := newByteCursor(payload)

	bitsPerSample, err := reader.readByte()
	if err != nil {
		return fmt.Errorf("truncated SOF segment")
	}
	height, err := reader.readUint16BE()
	if err != nil {
		return fmt.Errorf("truncated SOF segment")
	}
	width, err := reader.readUint16BE()
	if err != nil {
		return fmt.Errorf("truncated SOF segment")
	}
	numComponents, err := reader.readByte()
	if err != nil {
		return fmt.Errorf("truncated SOF segment")
	}

	header.BitsPerSample = bitsPerSample
	header.Height = height
	header.Width = width
	header.NumComponents = numComponents
	header.EncodingProcess = sofEncodingProcess(byte(marker & 0xFF))
	header.IsProgressive = byte(marker&0xFF) == 0xC2

	header.Components = make([]jpegComponent, 0, numComponents)
	for i := 0; i < int(numComponents); i++ {
		id, err := reader.readByte()
		if err != nil {
			return fmt.Errorf("truncated SOF component list")
		}
		sampling, err := reader.readByte()
		if err != nil {
			return fmt.Errorf("truncated SOF component list")
		}
		quantTable, err := reader.readByte()
		if err != nil {
			return fmt.Errorf("truncated SOF component list")
		}

		header.Components = append(header.Components, jpegComponent{
			ID:         id,
			HSampling:  sampling >> 4,
			VSampling:  sampling & 0x0F,
			QuantTable: quantTable,
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
