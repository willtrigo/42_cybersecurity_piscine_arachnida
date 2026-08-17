// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 17:54:45 by dande-je          #+#    #+#             //
//   Updated: 2026/08/17 19:24:23 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

type Crawler struct {
	httpClient domain.HTTPClient
	parser     domain.HTMLParser
	downloader *Downloader
	maxDepth   int
}

func NewCrawler(
	httpClient domain.HTTPClient,
	parser domain.HTMLParser,
	downloader *Downloader,
	maxDepth int,
) *Crawler {
	c := &Crawler{
		httpClient: httpClient,
		parser:     parser,
		downloader: downloader,
		maxDepth:   maxDepth,
	}
	return c
}
