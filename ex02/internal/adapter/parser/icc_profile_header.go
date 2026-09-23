// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   icc_profile_header.go                              :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/22 21:13:26 by dande-je          #+#    #+#             //
//   Updated: 2026/09/22 22:49:24 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const iccHeaderSize = 128

type iccTag struct {
	Signature [4]byte
	Offset    uint32
	Size      uint32
}

type iccProfile struct {
	Raw        []byte
	Version    string
	CMMType    string
	DeviceCls  string
	ColorSpace string
	ConnSpace  string
	DateTime   string
	Platform   string
	CMMFlags   string
	Creator    string
	ProfileID  string
	Tags       []iccTag
}

func decodeICCProfile(data []byte) (iccProfile, error) {
	if len(data) < iccHeaderSize {
		return iccProfile{}, fmt.Errorf("icc: header too short")
	}

	p := iccProfile{Raw: data}

	p.CMMType = iccCMMType(data[4:8])

	major := data[8]
	minor := data[9] >> 4
	patch := data[9] & 0x0F
	p.Version = fmt.Sprintf("%d.%d.%d", major, minor, patch)

	p.DeviceCls = iccDeviceClass(data[12:16])

	p.ColorSpace = iccColorSpace(data[16:20])

	p.ConnSpace = iccColorSpace(data[20:24])

	y := binary.BigEndian.Uint16(data[24:26])
	mo := binary.BigEndian.Uint16(data[26:28])
	d := binary.BigEndian.Uint16(data[28:30])
	h := binary.BigEndian.Uint16(data[30:32])
	mi := binary.BigEndian.Uint16(data[32:34])
	s := binary.BigEndian.Uint16(data[34:36])
	p.DateTime = fmt.Sprintf("%04d:%02d:%02d %02d:%02d:%02d", y, mo, d, h, mi, s)

	p.Platform = iccPlatform(data[40:44])

	p.CMMFlags = iccCMMFlags(binary.BigEndian.Uint32(data[44:48]))

	p.Creator = iccCMMType(data[80:84])

	p.ProfileID = strings.ToLower(fmt.Sprintf("%x", data[84:100]))

	if len(data) < iccHeaderSize+4 {
		return p, nil
	}
	count := binary.BigEndian.Uint32(data[128:132])
	pos := 132
	for i := uint32(0); i < count; i++ {
		if pos+12 > len(data) {
			break
		}
		var sig [4]byte
		copy(sig[:], data[pos:pos+4])
		p.Tags = append(p.Tags, iccTag{
			Signature: sig,
			Offset:    binary.BigEndian.Uint32(data[pos+4 : pos+8]),
			Size:      binary.BigEndian.Uint32(data[pos+8 : pos+12]),
		})
		pos += 12
	}

	return p, nil
}

func (p iccProfile) tags() []domain.Tag {
	tags := []domain.Tag{
		newICCTag("Profile CMM Type", p.CMMType),
		newICCTag("Profile Version", p.Version),
		newICCTag("Profile Class", p.DeviceCls),
		newICCTag("Color Space Data", p.ColorSpace),
		newICCTag("Profile Connection Space", p.ConnSpace),
		newICCTag("Profile Date Time", p.DateTime),
		newICCTag("Profile File Signature", "acsp"),
		newICCTag("Primary Platform", p.Platform),
		newICCTag("CMM Flags", p.CMMFlags),
		newICCTag("Profile Creator", p.Creator),
		newICCTag("Profile ID", p.ProfileID),
	}

	for _, t := range p.Tags {
		tags = append(tags, p.decodeTag(t)...)
	}
	return tags
}

func (p iccProfile) decodeTag(t iccTag) []domain.Tag {
	if int(t.Offset+8) > len(p.Raw) || t.Size < 8 {
		return nil
	}
	typ := string(p.Raw[t.Offset : t.Offset+4])
	body := p.Raw[t.Offset+8 : t.Offset+t.Size]

	switch string(t.Signature[:]) {
	case "desc":
		if v := iccTextDescription(body); v != "" {
			return []domain.Tag{newICCTag("Profile Description", v)}
		}
	case "cprt":
		if v := iccTextDescription(body); v != "" {
			return []domain.Tag{newICCTag("Profile Copyright", v)}
		}
	case "dmnd", "dmdd":
		if v := iccTextDescription(body); v != "" {
			return []domain.Tag{newICCTag("Device Manufacturer" /* or Model */, v)}
		}
	case "wtpt":
		if v := iccXYZ(body, typ); v != "" {
			return []domain.Tag{newICCTag("Media White Point", v)}
		}
	case "rXYZ":
		if v := iccXYZ(body, typ); v != "" {
			return []domain.Tag{newICCTag("Red Matrix Column", v)}
		}
	case "gXYZ":
		if v := iccXYZ(body, typ); v != "" {
			return []domain.Tag{newICCTag("Green Matrix Column", v)}
		}
	case "bXYZ":
		if v := iccXYZ(body, typ); v != "" {
			return []domain.Tag{newICCTag("Blue Matrix Column", v)}
		}
	case "chad":
		if v := iccChromaticAdaptation(body); v != "" {
			return []domain.Tag{newICCTag("Chromatic Adaptation", v)}
		}
	case "rTRC", "gTRC", "bTRC":
		label := map[string]string{
			"rTRC": "Red Tone Reproduction Curve",
			"gTRC": "Green Tone Reproduction Curve",
			"bTRC": "Blue Tone Reproduction Curve",
		}[string(t.Signature[:])]
		return []domain.Tag{newICCTag(label, fmt.Sprintf("(Binary data %d bytes)", len(body)))}
	}
	return nil
}

func iccTextDescription(body []byte) string {
	if len(body) < 4 {
		return ""
	}
	count := binary.BigEndian.Uint32(body[:4])
	if int(count) > len(body)-4 {
		return ""
	}
	return strings.TrimRight(string(body[4:4+count]), "\x00")
}

func iccXYZ(body []byte, typ string) string {
	if typ != "XYZ " || len(body) < 12 {
		return ""
	}
	x := float64(int32(binary.BigEndian.Uint32(body[0:4]))) / 65536.0  // #nosec G115
	y := float64(int32(binary.BigEndian.Uint32(body[4:8]))) / 65536.0  // #nosec G115
	z := float64(int32(binary.BigEndian.Uint32(body[8:12]))) / 65536.0 // #nosec G115
	return fmt.Sprintf("%.5g %.5g %.5g", x, y, z)
}

func iccChromaticAdaptation(body []byte) string {
	if len(body) < 36 {
		return ""
	}
	parts := make([]string, 9)
	for i := 0; i < 9; i++ {
		v := float64(int32(binary.BigEndian.Uint32(body[i*4:i*4+4]))) / 65536.0 // #nosec G115
		parts[i] = strconv.FormatFloat(v, 'g', -1, 64)
	}
	return strings.Join(parts, " ")
}

func iccCMMType(b []byte) string {
	if len(b) < 4 {
		return ""
	}
	switch string(b[:4]) {
	case "appl":
		return "Apple Computer Inc."
	case "ADBE":
		return "Adobe Systems Inc."
	case "MSFT":
		return "Microsoft Corp."
	case "\x00\x00\x00\x00":
		return ""
	default:
		return strings.TrimRight(string(b[:4]), "\x00 ")
	}
}

func iccDeviceClass(b []byte) string {
	switch string(b) {
	case "scnr":
		return "Input Device Profile"
	case "mntr":
		return "Display Device Profile"
	case "prtr":
		return "Output Device Profile"
	case "link":
		return "DeviceLink Profile"
	case "spac":
		return "ColorSpace Conversion Profile"
	case "abst":
		return "Abstract Profile"
	case "nmcl":
		return "Named Color Profile"
	default:
		return string(b)
	}
}

func iccColorSpace(b []byte) string {
	switch string(b) {
	case "XYZ ":
		return "XYZ"
	case "Lab ":
		return "Lab"
	case "Luv ":
		return "Luv"
	case "YCbr":
		return "YCbCr"
	case "Yxy ":
		return "Yxy"
	case "RGB ":
		return "RGB"
	case "GRAY":
		return "Gray"
	case "HSV ":
		return "HSV"
	case "HLS ":
		return "HLS"
	case "CMYK":
		return "CMYK"
	case "CMY ":
		return "CMY"
	case "2CLR":
		return "2 color"
	case "3CLR":
		return "3 color"
	case "4CLR":
		return "4 color"
	case "5CLR":
		return "5 color"
	case "6CLR":
		return "6 color"
	case "7CLR":
		return "7 color"
	case "8CLR":
		return "8 color"
	case "9CLR":
		return "9 color"
	case "ACLR":
		return "10 color"
	case "BCLR":
		return "11 color"
	case "CCLR":
		return "12 color"
	case "DCLR":
		return "13 color"
	case "ECLR":
		return "14 color"
	case "FCLR":
		return "15 color"
	default:
		return string(b)
	}
}

func iccPlatform(b []byte) string {
	switch string(b) {
	case "APPL":
		return "Apple Computer Inc."
	case "MSFT":
		return "Microsoft Corp."
	case "SUNW":
		return "Sun Microsystems Inc."
	case "SGI ":
		return "Silicon Graphics Inc."
	case "TGNT":
		return "Taligent Inc."
	case "\x00\x00\x00\x00":
		return ""
	default:
		return strings.TrimRight(string(b), "\x00 ")
	}
}

func iccCMMFlags(v uint32) string {
	embedded := v&1 != 0
	independent := v&2 != 0
	parts := []string{}
	if embedded {
		parts = append(parts, "Embedded")
	} else {
		parts = append(parts, "Not Embedded")
	}
	if independent {
		parts = append(parts, "Independent")
	} else {
		parts = append(parts, "Not Independent")
	}
	return strings.Join(parts, ", ")
}
