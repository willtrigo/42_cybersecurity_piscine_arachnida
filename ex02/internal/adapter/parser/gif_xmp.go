// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif_xmp.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/12 22:09:49 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 22:09:03 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import "bytes"

const (
	xmpPacketEndMarker = "<?xpacket end="
	xmpMetaEndTag      = "</x:xmpmeta>"
)

func readXMPPacket(c *byteCursor) ([]byte, error) {
	packetEnd := xmpPacketEnd(c.peekRemaining())
	if packetEnd == -1 {
		return nil, c.skipSubBlocks()
	}

	packet := append([]byte(nil), c.peekRemaining()[:packetEnd]...)
	if err := c.skip(packetEnd); err != nil {
		return nil, err
	}

	_ = c.skipSubBlocks()

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
