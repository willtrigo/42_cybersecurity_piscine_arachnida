// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif_header.go                                      :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/12 14:45:46 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 23:49:43 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	gifSignatureSize = 6
	gifSignature87a  = "GIF87a"
	gifSignature89a  = "GIF89a"

	globalColorTableFlagMask = 0x80
	colorTableSizeMask       = 0x07
	colorResolutionMask      = 0x70
	colorResolutionShift     = 4
	bytesPerEntry            = 3

	trailer             = 0x3b
	extensionIntroducer = 0x21

	graphicControlLabel     = 0xf9
	applicationLabel        = 0xff
	plainTextLabel          = 0x01
	plainTextFixedBlockSize = 12
	imageSeparator          = 0x2c

	graphicControlPayloadSize = 4
	packedFieldOffset         = 0
	delayFieldOffset          = 1
	delayFieldSize            = 2
	transparentIndexOffset    = 3
	transparentColorFlagMask  = 0x01

	xmpApplicationID = "XMP DataXMP"

	netscapeApplicationID  = "NETSCAPE2.0"
	loopPayloadSize        = 3
	loopFieldIndexOffset   = 0
	netscapeLoopSubBlockID = 0x01
	loopFieldOffset        = 1
	loopFieldSize          = 2

	imageDescriptorFixedSize = 9
)

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

func decodeGIFHeader(data []byte) (gifHeader, error) {
	cursor := newByteCursor(data)

	signature, err := cursor.readBytes(gifSignatureSize)
	if err != nil {
		return gifHeader{}, fmt.Errorf("truncated signature")
	}
	version := string(signature)
	if version != gifSignature87a && version != gifSignature89a {
		return gifHeader{}, fmt.Errorf("invalid signature")
	}

	header := gifHeader{Version: strings.TrimPrefix(version, "GIF")}

	if err := decodeLogicalScreenDescriptor(cursor, &header); err != nil {
		return gifHeader{}, fmt.Errorf("truncated logical screen descriptor")
	}

	decodeGIFBlocks(cursor, &header)

	return header, nil
}

func decodeLogicalScreenDescriptor(c *byteCursor, header *gifHeader) error {
	width, err := c.readUint16LE()
	if err != nil {
		return err
	}
	height, err := c.readUint16LE()
	if err != nil {
		return err
	}
	packed, err := c.readByte()
	if err != nil {
		return err
	}
	backgroundColorIndex, err := c.readByte()
	if err != nil {
		return err
	}
	if err := c.skip(1); err != nil {
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
		if err := c.skip(colorTableByteSize(tableSize)); err != nil {
			return err
		}
	}

	return nil
}

func colorTableByteSize(sizeField int) int {
	return bytesPerEntry * (1 << (sizeField + 1))
}

func decodeGIFBlocks(c *byteCursor, header *gifHeader) {
	var totalDelay int

	for {
		introducer, err := c.readByte()
		if err != nil || introducer == trailer {
			break
		}

		switch introducer {
		case extensionIntroducer:
			if !decodeExtensionBlock(c, header, &totalDelay) {
				header.DurationCentiseconds = totalDelay
				return
			}
		case imageSeparator:
			if err := skipImageDescriptor(c); err != nil {
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

func decodeExtensionBlock(c *byteCursor, header *gifHeader, totalDelay *int) bool {
	label, err := c.readByte()
	if err != nil {
		return false
	}

	switch label {
	case graphicControlLabel:
		delay, transparentIndex, hasTransparent, err := readGraphicControlExtension(c)
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
		appID, loopCount, hasLoop, xmpPacket, err := readApplicationExtension(c)
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
		if err := c.skip(plainTextFixedBlockSize); err != nil {
			return false
		}
		return c.skipSubBlocks() == nil

	default:
		return c.skipSubBlocks() == nil
	}
}

func readGraphicControlExtension(c *byteCursor) (delayCentiseconds int, transparentIndex uint8, hasTransparent bool, err error) {
	payload, err := c.readSubBlocks()
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

func readApplicationExtension(c *byteCursor) (appID string, loopCount uint16, hasLoop bool, xmpPacket []byte, err error) {
	idSize, err := c.readByte()
	if err != nil {
		return "", 0, false, nil, err
	}
	idBytes, err := c.readBytes(int(idSize))
	if err != nil {
		return "", 0, false, nil, err
	}
	appID = string(idBytes)

	if appID == xmpApplicationID {
		xmpPacket, err = readXMPPacket(c)
		return appID, 0, false, xmpPacket, err
	}

	data, err := c.readSubBlocks()
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

func skipImageDescriptor(c *byteCursor) error {
	descriptor, err := c.readBytes(imageDescriptorFixedSize)
	if err != nil {
		return err
	}

	packed := descriptor[imageDescriptorFixedSize-1]
	if packed&globalColorTableFlagMask != 0 {
		tableSize := int(packed & colorTableSizeMask)
		if err := c.skip(colorTableByteSize(tableSize)); err != nil {
			return err
		}
	}

	if _, err := c.readByte(); err != nil {
		return err
	}

	return c.skipSubBlocks()
}
