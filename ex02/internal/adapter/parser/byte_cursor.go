// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   byte_cursor.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/12 19:53:17 by dande-je          #+#    #+#             //
//   Updated: 2026/09/23 09:43:10 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bytes"
	"encoding/binary"
	"io"
)

type byteCursor struct {
	data []byte
	pos  int
}

func newByteCursor(data []byte) *byteCursor {
	return &byteCursor{data: data}
}

func (c *byteCursor) readBytes(n int) ([]byte, error) {
	if c.remaining() < n {
		return nil, io.ErrUnexpectedEOF
	}
	b := c.data[c.pos : c.pos+n]
	c.pos += n
	return b, nil
}

func (c *byteCursor) remaining() int {
	return len(c.data) - c.pos
}

func (c *byteCursor) readUint16LE() (uint16, error) {
	b, err := c.readBytes(uint16Size)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func (c *byteCursor) readUint16BE() (uint16, error) {
	b, err := c.readBytes(uint16Size)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}

func (c *byteCursor) readUint32LE() (uint32, error) {
	b, err := c.readBytes(uint32Size)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}

func (c *byteCursor) readUint32BE() (uint32, error) {
	b, err := c.readBytes(uint32Size)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

func (c *byteCursor) readByte() (byte, error) {
	if c.remaining() < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	b := c.data[c.pos]
	c.pos++
	return b, nil
}

func (c *byteCursor) readUntil(delim byte) ([]byte, error) {
	idx := bytes.IndexByte(c.data[c.pos:], delim)
	if idx == -1 {
		return nil, io.ErrUnexpectedEOF
	}
	b := c.data[c.pos : c.pos+idx]
	c.pos += idx + 1
	return b, nil
}

func (c *byteCursor) readNullTerminatedString() (string, error) {
	b, err := c.readUntil(0)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *byteCursor) skip(n int) error {
	if c.remaining() < n {
		return io.ErrUnexpectedEOF
	}
	c.pos += n
	return nil
}

func (c *byteCursor) readSubBlocks() ([]byte, error) {
	var buf bytes.Buffer
	for {
		size, err := c.readByte()
		if err != nil {
			return nil, err
		}
		if size == 0 {
			return buf.Bytes(), nil
		}
		chunk, err := c.readBytes(int(size))
		if err != nil {
			return nil, err
		}
		buf.Write(chunk)
	}
}

func (c *byteCursor) peekRemaining() []byte {
	return c.data[c.pos:]
}

func (c *byteCursor) skipSubBlocks() error {
	_, err := c.readSubBlocks()
	return err
}
