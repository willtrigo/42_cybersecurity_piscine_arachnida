// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   stat.go                                            :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/26 16:32:28 by dande-je          #+#    #+#             //
//   Updated: 2026/08/26 16:41:37 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package filesystem

import (
	"fmt"
	"os"
	"time"
)

type OSStatReader struct{}

func NewOSSStatReader() *OSStatReader {
	return &OSStatReader{}
}

func (OSStatReader) Stat(path string) (int64, time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("stat: %w", err)
	}
	return info.Size(), info.ModTime(), nil
}
