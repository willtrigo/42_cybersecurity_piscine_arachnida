// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   stat.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 16:32:28 by dande-je          #+#    #+#             //
//   Updated: 2026/09/11 21:11:21 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package filesystem

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex02/internal/domain"
)

type OSStatReader struct{}

func NewOSStatReader() *OSStatReader {
	return &OSStatReader{}
}

func (OSStatReader) StatIsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("stat: %w", err)
	}
	if info.IsDir() {
		return true, nil
	}
	return false, nil
}

func (OSStatReader) FileStat(path string) (domain.FileStat, error) {
	info, err := os.Stat(path)
	if err != nil {
		return domain.FileStat{}, fmt.Errorf("stat: %w", err)
	}

	accessTime, changeTime := info.ModTime(), info.ModTime()
	if raw, ok := info.Sys().(*syscall.Stat_t); ok {
		accessTime = time.Unix(raw.Atim.Sec, raw.Atim.Nsec)
		changeTime = time.Unix(raw.Ctim.Sec, raw.Ctim.Nsec)
	}

	return domain.FileStat{
		Size:       info.Size(),
		ModTime:    info.ModTime(),
		AccessTime: accessTime,
		ChangeTime: changeTime,
		Mode:       info.Mode(),
	}, nil
}
