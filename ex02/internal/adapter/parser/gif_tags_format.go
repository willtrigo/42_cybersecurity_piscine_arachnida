// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   gif_tags_format.go                                 :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/09/24 17:29:44 by dande-je          #+#    #+#             //
//   Updated: 2026/09/24 17:30:22 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package parser

import "fmt"

const centisecondsPerSecond = 100

func formatYesNo(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func formatAnimationIterations(loopCount uint16) string {
	if loopCount == 0 {
		return "Infinite"
	}
	return fmt.Sprintf("%d", loopCount)
}

func formatDuration(centiseconds int) string {
	return fmt.Sprintf("%.2f s", float64(centiseconds)/centisecondsPerSecond)
}
