// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif_header.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/12 14:45:46 by dande-je          #+#    #+#             //
//   Updated: 2026/09/24 19:40:22 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	gifSignatureSize = 6
	gifSignature87a  = "GIF87a"
	gifSignature89a  = "GIF89a"
)

const (
	globalColorTableFlagMask = 0x80
	colorTableSizeMask       = 0x07
	colorResolutionMask      = 0x70
	colorResolutionShift     = 4
)

const bytesPerEntry = 3

const (
	trailer             = 0x3b
	extensionIntroducer = 0x21
	imageSeparator      = 0x2c
)

const (
	graphicControlLabel     = 0xf9
	applicationLabel        = 0xff
	netscapeApplicationID   = "NETSCAPE2.0"
	xmpApplicationID        = "XMP DataXMP"
	plainTextLabel          = 0x01
	plainTextFixedBlockSize = 12
)

const (
	graphicControlPayloadSize = 4
	packedFieldOffset         = 0
	delayFieldOffset          = 1
	delayFieldSize            = 2
	transparentIndexOffset    = 3
	transparentColorFlagMask  = 0x01
)

const (
	loopPayloadSize        = 3
	loopFieldIndexOffset   = 0
	netscapeLoopSubBlockID = 0x01
	loopFieldOffset        = 1
	loopFieldSize          = 2
)

const (
	xmpPacketEndMarker = "<?xpacket end="
	xmpMetaEndTag      = "</x:xmpmeta>"
)

const imageDescriptorFixedSize = 9

type gifHeader struct {
	Version               string
	XMPPacket             []byte
	FrameCount            int
	ColorResolutionDepth  int
	DurationCentiseconds  int
	BitsPerPixel          int
	Width                 uint32
	Height                uint32
	LoopCount             uint16
	BackgroundColorIndex  uint8
	TransparentColorIndex uint8
	HasGlobalColorTable   bool
	HasLoopCount          bool
	HasTransparantColor   bool
}

type gifReader struct {
	*byteCursor
}

func decodeGIFHeader(data []byte) (header gifHeader, err error) {
	return decodeWithRecover("GIF", func() gifHeader {
		return parseGifHeader(&gifReader{newByteCursor(data)})
	})
}

func parseGifHeader(reader *gifReader) gifHeader {
	signature, err := reader.readBytes(gifSignatureSize)
	if err != nil {
		panic(domain.ErrTruncatedSignature)
	}
	version := string(signature)
	if version != gifSignature87a && version != gifSignature89a {
		panic(domain.ErrInvalidSignature)
	}

	header := gifHeader{Version: strings.TrimPrefix(version, "GIF")}

	if err := decodeLogicalScreenDescriptor(reader, &header); err != nil {
		panic(fmt.Errorf("truncated logical screen descriptor: %w", err))
	}

	decodeGIFBlocks(reader, &header)

	return header
}

func decodeLogicalScreenDescriptor(r *gifReader, header *gifHeader) error {
	width, err := r.readUint16LE()
	if err != nil {
		return err
	}
	height, err := r.readUint16LE()
	if err != nil {
		return err
	}
	packed, err := r.readByte()
	if err != nil {
		return err
	}
	backgroundColorIndex, err := r.readByte()
	if err != nil {
		return err
	}
	if err := r.skip(1); err != nil {
		return err
	}

	hasGlobalColorTable := packed&globalColorTableFlagMask != 0
	tableSize := int(packed & colorTableSizeMask)

	header.Width = uint32(width)
	header.Height = uint32(height)
	header.HasGlobalColorTable = hasGlobalColorTable
	header.ColorResolutionDepth = int((packed&colorResolutionMask)>>colorResolutionShift) + 1
	header.BitsPerPixel = tableSize + 1
	header.BackgroundColorIndex = backgroundColorIndex

	if hasGlobalColorTable {
		if err := r.skip(colorTableByteSize(tableSize)); err != nil {
			return err
		}
	}

	return nil
}

func colorTableByteSize(sizeField int) int {
	return bytesPerEntry * (1 << (sizeField + 1))
}

func decodeGIFBlocks(r *gifReader, header *gifHeader) {
	var totalDelay int

	for {
		introducer, err := r.readByte()
		if err != nil || introducer == trailer {
			break
		}

		switch introducer {
		case extensionIntroducer:
			if !decodeExtensionBlock(r, header, &totalDelay) {
				header.DurationCentiseconds = totalDelay
				return
			}
		case imageSeparator:
			if err := skipImageDescriptor(r); err != nil {
				header.DurationCentiseconds = totalDelay
				return
			}
			header.FrameCount++
		default:
			header.DurationCentiseconds = totalDelay
			return
		}
	}

	header.DurationCentiseconds = totalDelay
}

func decodeExtensionBlock(r *gifReader, header *gifHeader, totalDelay *int) bool {
	label, err := r.readByte()
	if err != nil {
		return false
	}

	switch label {
	case graphicControlLabel:
		delay, transparentIndex, hasTransparent, err := readGraphicControlExtension(r)
		if err != nil {
			return false
		}
		*totalDelay += delay
		if hasTransparent && !header.HasTransparantColor {
			header.HasTransparantColor = true
			header.TransparentColorIndex = transparentIndex
		}
		return true

	case applicationLabel:
		appID, loopCount, hasLoop, xmpPacket, err := readApplicationExtension(r)
		if err != nil {
			return false
		}
		switch appID {
		case netscapeApplicationID:
			header.HasLoopCount = hasLoop
			header.LoopCount = loopCount
		case xmpApplicationID:
			header.XMPPacket = xmpPacket
		}
		return true

	case plainTextLabel:
		if err := r.skip(plainTextFixedBlockSize); err != nil {
			return false
		}
		return r.skipSubBlocks() == nil

	default:
		return r.skipSubBlocks() == nil
	}
}

func readGraphicControlExtension(r *gifReader) (delayCentiseconds int, transparentIndex uint8, hasTransparent bool, err error) {
	payload, err := r.readSubBlocks()
	if err != nil {
		return 0, 0, false, err
	}
	if len(payload) < graphicControlPayloadSize {
		return 0, 0, false, fmt.Errorf("truncated graphic control extension")
	}

	packed := payload[packedFieldOffset]
	delayCentiseconds = int(binary.LittleEndian.Uint16(payload[delayFieldOffset : delayFieldOffset+delayFieldSize]))
	transparentIndex = payload[transparentIndexOffset]
	hasTransparent = packed&transparentColorFlagMask != 0

	return delayCentiseconds, transparentIndex, hasTransparent, nil
}

func readApplicationExtension(r *gifReader) (appID string, loopCount uint16, hasLoop bool, xmpPacket []byte, err error) {
	idSize, err := r.readByte()
	if err != nil {
		return "", 0, false, nil, err
	}
	idBytes, err := r.readBytes(int(idSize))
	if err != nil {
		return "", 0, false, nil, err
	}
	appID = string(idBytes)

	if appID == xmpApplicationID {
		xmpPacket, err = readXMPPacket(r)
		return appID, 0, false, xmpPacket, err
	}

	data, err := r.readSubBlocks()
	if err != nil {
		return appID, 0, false, nil, err
	}

	if appID == netscapeApplicationID {
		if len(data) >= loopPayloadSize && data[loopFieldIndexOffset] == netscapeLoopSubBlockID {
			loopCount = binary.BigEndian.Uint16(data[loopFieldOffset : loopFieldOffset+loopFieldSize])
			hasLoop = true
		}
	}

	return appID, loopCount, hasLoop, nil, nil
}

func readXMPPacket(r *gifReader) ([]byte, error) {
	packetEnd := xmpPacketEnd(r.peekRemaining())
	if packetEnd == -1 {
		return nil, r.skipSubBlocks()
	}

	packet := append([]byte(nil), r.peekRemaining()[:packetEnd]...)
	if err := r.skip(packetEnd); err != nil {
		return nil, err
	}

	_ = r.skipSubBlocks()

	return packet, nil
}

func xmpPacketEnd(data []byte) int {
	if piStart := bytes.Index(data, []byte(xmpPacketEndMarker)); piStart != -1 {
		if closeOffset := bytes.Index(data[piStart:], []byte("?>")); closeOffset != -1 {
			return piStart + closeOffset + len("?>")
		}
	}
	if metaEnd := bytes.Index(data, []byte(xmpMetaEndTag)); metaEnd != -1 {
		return metaEnd + len(xmpMetaEndTag)
	}
	return -1
}

func skipImageDescriptor(r *gifReader) error {
	descriptor, err := r.readBytes(imageDescriptorFixedSize)
	if err != nil {
		return err
	}

	packed := descriptor[imageDescriptorFixedSize-1]
	if packed&globalColorTableFlagMask != 0 {
		tableSize := int(packed & colorTableSizeMask)
		if err := r.skip(colorTableByteSize(tableSize)); err != nil {
			return err
		}
	}

	if _, err := r.readByte(); err != nil {
		return err
	}

	return r.skipSubBlocks()
}
