// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   tiff.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/26 19:20:09 by dande-je          #+#    #+#             //
//   Updated: 2026/09/26 19:22:01 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import "encoding/binary"

type tiffView struct {
	data  []byte
	order binary.ByteOrder
}

func newTIFFView(data []byte, order binary.ByteOrder) tiffView {
	return tiffView{data: data, order: order}
}

func (v tiffView) slice(off, n int) ([]byte, bool) {
	if off < 0 || n < 0 || off > len(v.data) || off+n > len(v.data) {
		return nil, false
	}
	return v.data[off : off+n], true
}

func (v tiffView) u16(off int) (uint16, bool) {
	b, ok := v.slice(off, uint16Size)
	if !ok {
		return 0, false
	}
	return v.order.Uint16(b), true
}

func (v tiffView) u32(off int) (uint32, bool) {
	b, ok := v.slice(off, uint32Size)
	if !ok {
		return 0, false
	}
	return v.order.Uint32(b), true
}

func (v tiffView) copySlice(off, n int) ([]byte, bool) {
	b, ok := v.slice(off, n)
	if !ok {
		return nil, false
	}
	return append([]byte(nil), b...), true
}
