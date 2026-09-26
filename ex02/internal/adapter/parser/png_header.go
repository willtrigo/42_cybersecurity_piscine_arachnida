// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   png_header.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/14 11:40:33 by dande-je          #+#    #+#             //
//   Updated: 2026/09/26 20:41:15 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const pngSignatureSize = 8

const (
	pngChunkLengthSize = 4
	pngChunkTypeSize   = 4
)

const (
	pngChunkHeaderSize = pngChunkLengthSize + pngChunkTypeSize
	pngChunkCRCSize    = 4
)

const (
	chunkTypeIHDR = "IHDR"
	chunkTypeICCP = "iCCP"
	chunkTypeText = "tEXt"
	chunkTypeZTXt = "zTXt"
	chunkTypeITXt = "iTXt"
	chunkTypeEXIf = "eXIf"
	chunkTypeIEND = "IEND"
)

const (
	ihdrByteFieldCount   = 5
	bitDepthIdx          = 0
	colorTypeIdx         = 1
	compressionMethodIdx = 2
	filterMethodIdx      = 3
	interlaceMethodIdx   = 4
)

const (
	rawProfileTypeEXIf = "exif"
	rawProfileTypeXMP  = "xmp"
)

const (
	rawProfileKeywordPrefix = "Raw profile type "
	rawProfileHeaderLines   = 3
	profileTypeLengthIndex  = 1
	profileTypeDecodedIndex = 2
	profileTypeIndex        = 0
)

const (
	pngITXtFlagsSize = 2
	xmpKeyword       = "XML:com.adobe.xmp"
)

var pngSignature = [pngSignatureSize]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}

type pngTextEntry struct {
	Keyword string
	Value   string
}

type pngHeader struct {
	TextEntries       []pngTextEntry
	ICCProfileName    string
	XMPPacket         []byte
	EXIFData          exifHeader
	Width             uint32
	Height            uint32
	BitDepth          uint8
	ColorType         uint8
	CompressionMethod uint8
	FilterMethod      uint8
	InterlaceMethod   uint8
}

type pngReader struct {
	*byteCursor
}

type pngChunk struct {
	Type string
	Data []byte
}

func decodePNGHeader(data []byte) (pngHeader, error) {
	return decodeWithRecover("PNG", func() pngHeader {
		return parsePngHeader(&pngReader{newByteCursor(data)})
	})
}

func parsePngHeader(reader *pngReader) pngHeader {
	sig, err := reader.slice(0, pngSignatureSize)
	if err != nil || !bytes.Equal(sig, pngSignature[:]) {
		panic(domain.ErrInvalidSignature)
	}

	var header pngHeader
	sawIHDR := false
	pos := pngSignatureSize
	raw := reader.bytes()

	for pos+pngChunkHeaderSize+pngChunkCRCSize <= len(raw) {
		chunk, nextPos, err := readPNGChunk(raw, pos)
		if err != nil {
			panic(err)
		}
		pos = nextPos

		done, err := applyPNGChunk(chunk, &header, &sawIHDR)
		if err != nil {
			panic(err)
		}
		if done {
			break
		}
	}

	if !sawIHDR {
		panic(domain.ErrMissingIHDRChunk)
	}

	return header
}

func readPNGChunk(data []byte, pos int) (pngChunk, int, error) {
	reader := newByteCursor(data[pos:])

	length, err := reader.readUint32BE()
	if err != nil {
		return pngChunk{}, 0, domain.ErrTruncatedChunkHeader
	}

	chunkType, err := reader.readBytes(pngChunkTypeSize)
	if err != nil {
		return pngChunk{}, 0, domain.ErrTruncatedChunkType
	}

	payload, err := reader.readBytes(int(length))
	if err != nil {
		return pngChunk{}, 0, fmt.Errorf("truncated chunk %q", string(chunkType))
	}

	chunk := pngChunk{
		Type: string(chunkType),
		Data: payload,
	}

	return chunk, pos + pngChunkHeaderSize + int(length) + pngChunkCRCSize, nil
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
			recordPNGText(header, entry)
		}

	case chunkTypeZTXt:
		if entry, ok := decodeZTXtChunk(chunk.Data); ok {
			recordPNGText(header, entry)
		}

	case chunkTypeITXt:
		if entry, packet, isXMP := decodeITXtChunk(chunk.Data); entry.Keyword != "" || isXMP {
			if isXMP {
				header.XMPPacket = packet
			} else {
				recordPNGText(header, entry)
			}
		}

	case chunkTypeEXIf:
		exif, err := decodeEXIFData(append([]byte(nil), chunk.Data...))
		if err != nil {
			return false, err
		}
		header.EXIFData = exif

	case chunkTypeIEND:
		return true, nil
	}

	return false, nil
}

func decodeIHDRChunk(data []byte, header *pngHeader) error {
	reader := newByteCursor(data)

	width, err := reader.readUint32BE()
	if err != nil {
		return domain.ErrTruncatedIHDRChunk
	}
	height, err := reader.readUint32BE()
	if err != nil {
		return domain.ErrTruncatedIHDRChunk
	}
	fields, err := reader.readBytes(ihdrByteFieldCount)
	if err != nil {
		return domain.ErrTruncatedIHDRChunk
	}

	header.Width = width
	header.Height = height
	header.BitDepth = fields[bitDepthIdx]
	header.ColorType = fields[colorTypeIdx]
	header.CompressionMethod = fields[compressionMethodIdx]
	header.FilterMethod = fields[filterMethodIdx]
	header.InterlaceMethod = fields[interlaceMethodIdx]

	return nil
}

func recordPNGText(header *pngHeader, entry pngTextEntry) {
	if profileType, payload, ok := decodeImageMagickRawProfile(entry.Keyword, entry.Value); ok {
		switch profileType {
		case rawProfileTypeEXIf:
			if exif, err := decodeEXIFData(payload); err == nil {
				header.EXIFData = exif
			}
		case rawProfileTypeXMP:
			header.XMPPacket = payload
		default:
			header.TextEntries = append(header.TextEntries, entry)
		}
		return
	}

	header.TextEntries = append(header.TextEntries, entry)
}

func decodeImageMagickRawProfile(keyword, value string) (string, []byte, bool) {
	if !strings.HasPrefix(keyword, rawProfileKeywordPrefix) {
		return "", nil, false
	}

	lines := strings.SplitN(strings.TrimLeft(value, "\n"), "\n", rawProfileHeaderLines)
	if len(lines) < rawProfileHeaderLines {
		return "", nil, false
	}

	length, err := strconv.Atoi(strings.TrimSpace(lines[profileTypeLengthIndex]))
	if err != nil || length < 0 {
		return "", nil, false
	}

	decoded, err := hex.DecodeString(stripNonHexDigitis(lines[profileTypeDecodedIndex]))
	if err != nil || len(decoded) < length {
		return "", nil, false
	}

	return strings.ToLower(strings.TrimSpace(lines[profileTypeIndex])), decoded[:length], true
}

func stripNonHexDigitis(s string) string {
	digits := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if c := s[i]; isHexDigit(c) {
			digits = append(digits, c)
		}
	}
	return string(digits)
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func decodeICCPChunk(data []byte) string {
	name, err := newByteCursor(data).readNullTerminatedString()
	if err != nil {
		return ""
	}
	return name
}

func decodeTEXtChunk(data []byte) (pngTextEntry, bool) {
	reader := newByteCursor(data)

	keyword, err := reader.readNullTerminatedString()
	if err != nil {
		return pngTextEntry{}, false
	}

	return pngTextEntry{
		Keyword: keyword,
		Value:   string(reader.peekRemaining()),
	}, true
}

func decodeZTXtChunk(data []byte) (pngTextEntry, bool) {
	reader := newByteCursor(data)

	keyword, err := reader.readNullTerminatedString()
	if err != nil {
		return pngTextEntry{}, false
	}

	compressionMethod, err := reader.readByte()
	if err != nil || compressionMethod != 0 {
		return pngTextEntry{}, false
	}

	text, err := inflateZlib(reader.peekRemaining())
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
	reader := newByteCursor(data)

	keyword, err := reader.readNullTerminatedString()
	if err != nil {
		return pngTextEntry{}, nil, false
	}

	flags, err := reader.readBytes(pngITXtFlagsSize)
	if err != nil {
		return pngTextEntry{}, nil, false
	}
	compressionFlag := flags[0]

	if _, err := reader.readNullTerminatedString(); err != nil {
		return pngTextEntry{}, nil, false
	}
	if _, err := reader.readNullTerminatedString(); err != nil {
		return pngTextEntry{}, nil, false
	}

	text := reader.peekRemaining()

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
