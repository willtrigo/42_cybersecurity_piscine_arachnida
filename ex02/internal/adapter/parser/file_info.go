// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   file_info.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/03 17:33:49 by dande-je          #+#    #+#             //
//   Updated: 2026/09/11 19:13:12 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import (
	"fmt"
	"path/filepath"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

const fileInfoIFDpath = "File"

const dateTimeLayout = "2006:01:02 15:04:05-07:00"

const (
	bytesPerKiB               = 1024
	decimalPrecisionThreshold = 10
)

var sizeSuffixes = []string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

func buildFileInfoTags(path string, stat domain.FileStat) []domain.Tag {
	return []domain.Tag{
		newFileInfoTag("File Name", filepath.Base(path)),
		newFileInfoTag("Directory", filepath.Dir(path)),
		newFileInfoTag("File Size", formatFileSize(stat.Size)),
		newFileInfoTag("File Modification Date/Time", stat.ModTime.Format(dateTimeLayout)),
		newFileInfoTag("File Access Date/Time", stat.AccessTime.Format(dateTimeLayout)),
		newFileInfoTag("File Inode Change Date/Time", stat.ChangeTime.Format(dateTimeLayout)),
		newFileInfoTag("File Permissions", stat.Mode.String()),
	}
}

func newFileInfoTag(name, value string) domain.Tag {
	return domain.Tag{IDFPath: fileInfoIFDpath, Name: name, Value: value}
}

func formatFileSize(size int64) string {
	if size < bytesPerKiB {
		return fmt.Sprintf("%d bytes", size)
	}

	div, exp := int64(bytesPerKiB), 0
	for n := size / bytesPerKiB; n >= bytesPerKiB && exp < len(sizeSuffixes)-1; n /= bytesPerKiB {
		div *= bytesPerKiB
		exp++
	}

	value := float64(size) / float64(div)
	suffix := sizeSuffixes[exp]

	if value < decimalPrecisionThreshold {
		return fmt.Sprintf("%.1f %s", value, suffix)
	}
	return fmt.Sprintf("%.0f %s", value, suffix)
}
