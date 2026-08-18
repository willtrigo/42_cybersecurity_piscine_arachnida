// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler_stats.go                                   :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/18 09:35:23 by dande-je          #+#    #+#             //
//   Updated: 2026/08/18 11:16:09 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import "sync/atomic"

type crawlStats struct {
	pagesVisited     atomic.Uint64
	imagesDownloaded atomic.Uint64
	errors           atomic.Uint64
}

func newCrawlStats() *crawlStats {
	return &crawlStats{}
}

func (s *crawlStats) incrementPagesVisited() {
	s.pagesVisited.Add(1)
}

func (s *crawlStats) incrementImagesDownloaded() {
	s.imagesDownloaded.Add(1)
}

func (s *crawlStats) incrementErrors() {
	s.errors.Add(1)
}

func (s *crawlStats) getPagesVisited() uint64 {
	return s.pagesVisited.Load()
}

func (s *crawlStats) getImagesDownloaded() uint64 {
	return s.imagesDownloaded.Load()
}

func (s *crawlStats) getErrors() uint64 {
	return s.errors.Load()
}
