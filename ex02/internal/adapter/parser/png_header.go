// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   png_header.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/14 11:40:33 by dande-je          #+#    #+#             //
//   Updated: 2026/09/14 13:10:26 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	pngSignatureSize = 8

	pngChunkLengthSize = 4
	pngChunkTypeSize   = 4
	pngChunkCRCSize    = 4
	pngChunkHeaderSize = pngChunkLengthSize + pngChunkTypeSize

	chunkTypeIHDR  = "IHDR"
	chunkTypeICCP  = "iCCP"
	chunkTypeText  = "tEXt"
	chunkTypeZTXt  = "zTXt"
	chunkTypeITest = "iTXt"
	chunkTypeEXIf  = "eXIf"
	chunkTypeIEND  = "IEND"
)

const (
	pngIHDRSize                 = 13
	ihdrWidthOffset             = 0
	ihdrHeightOffset            = 4
	ihdrBitDepthOffset          = 8
	ihdrColorTypeOffset         = 9
	ihdrCompressionMethodOffset = 10
	ihdrFilterMethodOffset      = 11
	ihdrInterlaceMethodOffset   = 12
)

const (
	pngITXtFlagsSize = 2
	xmpKeyword       = "XML:com.adobe.xmp"
)

var pngSignature = [pngSignatureSize]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}

type pngChunk struct {
	Type string
	Data []byte
}

type pngTextEntry struct {
	Keyword string
	Value   string
}

type pngHeader struct {
	TextEntries       []pngTextEntry
	ICCProfileName    string
	XMPPacket         []byte
	EXIFData          []byte
	Width             uint32
	Height            uint32
	BitDepth          uint8
	ColorType         uint8
	CompressionMethod uint8
	FilterMethod      uint8
	InterlaceMethod   uint8
}

func decodePNGHeader(data []byte) (pngHeader, error) {
	if len(data) < pngSignatureSize || !bytes.Equal(data[:pngSignatureSize], pngSignature[:]) {
		return pngHeader{}, fmt.Errorf("invalid signature")
	}

	var header pngHeader
	sawIHDR := false
	pos := pngSignatureSize

	for pos+pngChunkHeaderSize+pngChunkCRCSize <= len(data) {
		chunk, nextPos, err := readPNGChunk(data, pos)
		if err != nil {
			return pngHeader{}, err
		}
		pos = nextPos

		done, err := applyPNGChunk(chunk, &header, &sawIHDR)
		if err != nil {
			return pngHeader{}, err
		}
		if done {
			break
		}
	}

	if !sawIHDR {
		return pngHeader{}, fmt.Errorf("missing IHDR chunk")
	}

	return header, nil
}

func readPNGChunk(data []byte, pos int) (pngChunk, int, error) {
	length := binary.BigEndian.Uint32(data[pos : pos+pngChunkLengthSize])
	chunkType := string(data[pos+pngChunkLengthSize : pos+pngChunkHeaderSize])

	dataStart := pos + pngChunkHeaderSize
	dataEnd := dataStart + int(length)
	if dataEnd+pngChunkCRCSize > len(data) {
		return pngChunk{}, 0, fmt.Errorf("truncated chunk %q", chunkType)
	}

	chunk := pngChunk{
		Type: chunkType,
		Data: data[dataStart:dataEnd],
	}

	return chunk, dataEnd + pngChunkCRCSize, nil
}

func applyPNGChunk(chunk pngChunk, header *pngHeader, sawIHDR *bool) (done bool, err error) {
	switch chunk.Type {
	case chunkTypeIHDR:
		if err := decodeIHDRChunk(chunk.Data, header); err != nil {
			return false, err
		}
		*sawIHDR = true

	case chunkTypeICCP:
		header.ICCProfileName = decodeICCPChunk(chunk.Data)

	case chunkTypeText:
		if entry, ok := decodeTEXtChunk(chunk.Data); ok {
			header.TextEntries = append(header.TextEntries, entry)
		}

	case chunkTypeZTXt:
		if entry, ok := decodeZTXtChunk(chunk.Data); ok {
			header.TextEntries = append(header.TextEntries, entry)
		}

	case chunkTypeITest:
		if entry, packet, isXMP := decodeITXtChunk(chunk.Data); entry.Keyword != "" || isXMP {
			if isXMP {
				header.XMPPacket = packet
			} else {
				header.TextEntries = append(header.TextEntries, entry)
			}
		}

	case chunkTypeEXIf:
		header.EXIFData = append([]byte(nil), chunk.Data...)

	case chunkTypeIEND:
		return true, nil
	}

	return false, nil
}

func decodeIHDRChunk(data []byte, header *pngHeader) error {
	if len(data) < pngIHDRSize {
		return fmt.Errorf("truncated IHDR chunk")
	}

	header.Width = binary.BigEndian.Uint32(data[ihdrWidthOffset : ihdrWidthOffset+uint32Size])
	header.Height = binary.BigEndian.Uint32(data[ihdrHeightOffset : ihdrHeightOffset+uint32Size])
	header.BitDepth = data[ihdrBitDepthOffset]
	header.ColorType = data[ihdrColorTypeOffset]
	header.CompressionMethod = data[ihdrCompressionMethodOffset]
	header.FilterMethod = data[ihdrFilterMethodOffset]
	header.InterlaceMethod = data[ihdrInterlaceMethodOffset]

	return nil
}

func decodeICCPChunk(data []byte) string {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return ""
	}
	return string(data[:nullIndex])
}

func decodeTEXtChunk(data []byte) (pngTextEntry, bool) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return pngTextEntry{}, false
	}

	return pngTextEntry{
		Keyword: string(data[:nullIndex]),
		Value:   string(data[nullIndex+1:]),
	}, true
}

func decodeZTXtChunk(data []byte) (pngTextEntry, bool) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 || nullIndex+2 > len(data) {
		return pngTextEntry{}, false
	}

	keyword := string(data[:nullIndex])
	compressed := data[nullIndex+2:]

	text, err := inflateZlib(compressed)
	if err != nil {
		return pngTextEntry{}, false
	}

	return pngTextEntry{Keyword: keyword, Value: string(text)}, true
}

func inflateZlib(compressed []byte) (text []byte, err error) {
	reader, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, reader.Close())
	}()

	return io.ReadAll(reader)
}

func decodeITXtChunk(data []byte) (entry pngTextEntry, packet []byte, isXMP bool) {
	keywordEnd := bytes.IndexByte(data, 0)
	if keywordEnd == -1 {
		return pngTextEntry{}, nil, false
	}
	keyword := string(data[:keywordEnd])

	rest := data[keywordEnd+1:]
	if len(rest) < pngITXtFlagsSize {
		return pngTextEntry{}, nil, false
	}
	compressionFlag := rest[0]
	rest = rest[pngITXtFlagsSize:]

	languageTagEnd := bytes.IndexByte(rest, 0)
	if languageTagEnd == -1 {
		return pngTextEntry{}, nil, false
	}
	rest = rest[languageTagEnd+1:]

	translatedKeywordEnd := bytes.IndexByte(rest, 0)
	if translatedKeywordEnd == -1 {
		return pngTextEntry{}, nil, false
	}
	text := rest[translatedKeywordEnd+1:]

	if compressionFlag != 0 {
		inflated, err := inflateZlib(text)
		if err != nil {
			return pngTextEntry{}, nil, false
		}
		text = inflated
	}

	if keyword == xmpKeyword {
		return pngTextEntry{}, text, true
	}

	return pngTextEntry{Keyword: keyword, Value: string(text)}, nil, false
}
