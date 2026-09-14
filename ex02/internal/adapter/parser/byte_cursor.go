// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   byte_cursor.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/12 19:53:17 by dande-je          #+#    #+#             //
//   Updated: 2026/09/12 22:55:24 by dande-je         ###   ########.fr       //
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

func (c *byteCursor) readByte() (byte, error) {
	if c.remaining() < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	b := c.data[c.pos]
	c.pos++
	return b, nil
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
