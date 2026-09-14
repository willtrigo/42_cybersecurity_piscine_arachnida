// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   system_tags.go                                     :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/03 17:33:49 by dande-je          #+#    #+#             //
//   Updated: 2026/09/13 23:31:00 by dande-je         ###   ########.fr       //
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

func buildSystemTags(path string, stat domain.FileStat) []domain.Tag {
	return []domain.Tag{
		newSystemTags("File Name", filepath.Base(path)),
		newSystemTags("Directory", filepath.Dir(path)),
		newSystemTags("File Size", formatFileSize(stat.Size)),
		newSystemTags("File Modification Date/Time", stat.ModTime.Format(dateTimeLayout)),
		newSystemTags("File Access Date/Time", stat.AccessTime.Format(dateTimeLayout)),
		newSystemTags("File Inode Change Date/Time", stat.ChangeTime.Format(dateTimeLayout)),
		newSystemTags("File Permissions", stat.Mode.String()),
	}
}

func newSystemTags(name, value string) domain.Tag {
	return domain.NewTag(fileInfoIFDpath, name, value)
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
