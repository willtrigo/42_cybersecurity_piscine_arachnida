// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler_stats.go                                   :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/18 09:35:23 by dande-je          #+#    #+#             //
//   Updated: 2026/08/19 22:35:26 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import "sync/atomic"

type CrawlStats struct {
	pagesVisited     atomic.Uint64
	imagesDownloaded atomic.Uint64
	errors           atomic.Uint64
}

func newCrawlStats() *CrawlStats {
	return &CrawlStats{}
}

func (s *CrawlStats) incrementPagesVisited() {
	s.pagesVisited.Add(1)
}

func (s *CrawlStats) incrementImagesDownloaded() {
	s.imagesDownloaded.Add(1)
}

func (s *CrawlStats) incrementErrors() {
	s.errors.Add(1)
}

func (s *CrawlStats) getPagesVisited() uint64 {
	return s.pagesVisited.Load()
}

func (s *CrawlStats) getImagesDownloaded() uint64 {
	return s.imagesDownloaded.Load()
}

func (s *CrawlStats) getErrors() uint64 {
	return s.errors.Load()
}
