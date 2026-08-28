// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif.go                                             :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/27 11:09:02 by dande-je          #+#    #+#             //
//   Updated: 2026/08/28 15:35:15 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const (
	gifHeaderSize                  = 6
	gifLogicalScreenDescriptorSize = 7
	gifImageDescriptorSize         = 9
)

const (
	logicalScreenWidthOffsetBegin  = 0
	logicalScreenWidthOffsetEnd    = logicalScreenWidthOffsetBegin + 2
	logicalScreenHeightOffsetBegin = 2
	logicalScreenHeightOffsetEnd   = logicalScreenHeightOffsetBegin + 2
	logicalScreenPackedOffset      = 4
)

const (
	colorTableMask      = 0x80
	colorTableFlagShift = 1
	colorTableSizeMask  = 0x07
	colorTableShift     = 1
	colorTableEntrySize = 3
)

const (
	gifTrailer             = 0x3B
	gifExtensionIntroducer = 0x21
	gifCommentLabel        = 0xFE
	gifApplicationLabel    = 0xFF
	gifImageSeparator      = 0x2C
)

const (
	gifXMPIdentifier = "XMP DataXMP"
	gifXMPEndMarker  = "</x:xmpmeta>"
)

const regexFullMatchCount = 2

const imageDescriptorPackedOffset = 8

var xmpAttributeNames = map[string]map[string]string{
	"xmpmta": {
		"xmptk": "XMPToolkit",
	},
	"xmpmeta": {
		"xmptk": "XMPToolkit",
	},
	"Description": {
		"CreatorTool": "CreatorTool",
		"InstanceID":  "InstanceID",
		"DocumentID":  "DocumentID",
	},
	"DerivedFrom": {
		"instanceID": "DerivedFromInstanceID",
		"documentID": "DerivedFromDocumentID",
	},
}

type GIFParser struct{}

func NewGIFParser() *GIFParser {
	return &GIFParser{}
}

func (GIFParser) Read(path string) (metadata *domain.Metadata, err error) {
	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("gif: invalid file path")
	}

	data, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}
	defer func() {
		err = errors.Join(err, data.Close())
	}()

	r := bufio.NewReader(data)

	header := make([]byte, gifHeaderSize)
	if _, err = io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("gif: reading signature: %w", err)
	}

	lsd := make([]byte, gifLogicalScreenDescriptorSize)
	if _, err = io.ReadFull(r, lsd); err != nil {
		return nil, fmt.Errorf("gif: reading logical screen descriptor: %w", err)
	}
	width := int(binary.LittleEndian.Uint16(lsd[logicalScreenWidthOffsetBegin:logicalScreenWidthOffsetEnd]))
	height := int(binary.LittleEndian.Uint16(lsd[logicalScreenHeightOffsetBegin:logicalScreenHeightOffsetEnd]))
	packed := lsd[logicalScreenPackedOffset]

	if packed&colorTableMask != 0 {
		size := globalColorTableSize(packed)
		if _, err = io.CopyN(io.Discard, r, int64(size)); err != nil {
			return nil, fmt.Errorf("gif: skipping global color table: %w", err)
		}
	}

	tags, err := walkBlocks(r)
	if err != nil {
		return nil, fmt.Errorf("gif: %w", err)
	}

	return &domain.Metadata{
		Format:     domain.FormatGIF,
		Dimensions: domain.Dimensions{Width: width, Height: height},
		Tags:       tags,
	}, nil
}

func globalColorTableSize(packed byte) int {
	entries := colorTableFlagShift << ((packed & colorTableSizeMask) + colorTableShift)
	return colorTableEntrySize * entries
}

func walkBlocks(r *bufio.Reader) ([]domain.Tag, error) {
	var tags []domain.Tag

	for {
		introducer, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return tags, nil
			}
			return nil, fmt.Errorf("reading block introducer: %w", err)
		}

		switch introducer {
		case gifTrailer:
			return tags, nil
		case gifExtensionIntroducer:
			extTags, err := readExtension(r)
			if err != nil {
				return nil, err
			}
			tags = append(tags, extTags...)
		case gifImageSeparator:
			if err := skipImageBlock(r); err != nil {
				return nil, fmt.Errorf("skipping image block: %w", err)
			}
		default:
			return nil, fmt.Errorf("unexpected block introducer 0x%02X", introducer)
		}
	}
}

func readExtension(r *bufio.Reader) ([]domain.Tag, error) {
	label, err := r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("reading extension label: %w", err)
	}

	body, err := collectSubBlocks(r)
	if err != nil {
		return nil, fmt.Errorf("reading extension body: %w", err)
	}

	switch label {
	case gifCommentLabel:
		if len(body) == 0 {
			return nil, nil
		}
		return []domain.Tag{{Name: "Comment", Value: string(body)}}, nil
	case gifApplicationLabel:
		return extractXMPTags(body), nil
	default:
		return nil, nil
	}
}

func collectSubBlocks(r *bufio.Reader) ([]byte, error) {
	var out []byte
	for {
		size, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if size == 0 {
			return out, nil
		}
		chunk := make([]byte, size)
		if _, err = io.ReadFull(r, chunk); err != nil {
			return nil, err
		}
		out = append(out, chunk...)
	}
}

func extractXMPTags(body []byte) []domain.Tag {
	if !strings.HasPrefix(string(body), gifXMPIdentifier) {
		return nil
	}

	packet := body[len(gifXMPIdentifier):]
	if end := strings.LastIndex(string(packet), gifXMPEndMarker); end != -1 {
		packet = packet[:end+len(gifXMPEndMarker)]
	}
	if len(packet) == 0 {
		return nil
	}

	tags := parseXMPPackeet(packet)
	if len(tags) == 0 {
		return []domain.Tag{{Name: "XMP", Value: string(packet)}}
	}
	return tags
}

func parseXMPPackeet(packet []byte) []domain.Tag {
	var tags []domain.Tag

	dec := xml.NewDecoder(bytes.NewReader(packet))
	dec.Strict = false

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}

		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}

		wanted, ok := xmpAttributeNames[start.Name.Local]
		if !ok {
			continue
		}
		for _, attr := range start.Attr {
			if tagName, found := wanted[attr.Name.Local]; found {
				tags = append(tags, domain.Tag{Name: tagName, Value: strings.TrimSpace(attr.Value)})
			}
		}
	}

	if len(tags) <= 1 {
		tags = append(tags, extractBrokenXMPTags(string(packet))...)
	}

	return tags
}

func extractBrokenXMPTags(s string) []domain.Tag {
	patterns := []struct {
		re   *regexp.Regexp
		name string
	}{
		{regexp.MustCompile(`xmp:CreatorTool="([^"]+?)(?:\s+xmpMM:|")`), "CreatorTool"},
		{regexp.MustCompile(`xmpMM:InstanceID="([^"]+)"`), "InstanceID"},
		{regexp.MustCompile(`xmpMM:DocumentID="([^"]+)"`), "DocumentID"},
		{regexp.MustCompile(`stRef:instane?ID="([^"]+)"`), "DerivedFromInstanceID"},
		{regexp.MustCompile(`stRef:documentID="([^"]+)"`), "DerivedFromDocumentID"},
	}

	var tags []domain.Tag
	for _, p := range patterns {
		if m := p.re.FindStringSubmatch(s); len(m) == regexFullMatchCount {
			tags = append(tags, domain.Tag{Name: p.name, Value: strings.TrimSpace(m[1])})
		}
	}
	return tags
}

func skipImageBlock(r *bufio.Reader) error {
	descriptor := make([]byte, gifImageDescriptorSize)
	if _, err := io.ReadFull(r, descriptor); err != nil {
		return fmt.Errorf("reading image descriptor: %w", err)
	}
	packed := descriptor[imageDescriptorPackedOffset]

	if packed&colorTableMask != 0 {
		size := globalColorTableSize(packed)
		if _, err := io.CopyN(io.Discard, r, int64(size)); err != nil {
			return fmt.Errorf("skipping local color table; %w", err)
		}
	}

	if _, err := r.ReadByte(); err != nil {
		return fmt.Errorf("reading LZW code size: %w", err)
	}

	_, err := collectSubBlocks(r)
	return err
}
