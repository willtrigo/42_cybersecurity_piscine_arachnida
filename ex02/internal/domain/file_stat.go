// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   file_stat.go                                       :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/10 21:59:19 by dande-je          #+#    #+#             //
//   Updated: 2026/09/10 21:59:43 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package domain

import (
	"os"
	"time"
)

type FileStat struct {
	ModTime    time.Time
	AccessTime time.Time
	ChangeTime time.Time
	Mode       os.FileMode
	Size       int64
}
