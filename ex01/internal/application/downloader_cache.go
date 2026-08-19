// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   downloader_cache.go                                :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/18 22:43:16 by dande-je          #+#    #+#             //
//   Updated: 2026/08/19 00:34:24 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

type ImageCache struct {
	cache map[string]bool
}

type LinkCache struct {
	cache map[string]bool
}

func newImageCache() *ImageCache {
	return &ImageCache{
		cache: make(map[string]bool),
	}
}

func newLinkCache() *LinkCache {
	return &LinkCache{
		cache: make(map[string]bool),
	}
}

func (ic *ImageCache) Has(url string) bool {
	_, ok := ic.cache[url]
	return ok
}

func (ic *ImageCache) Add(url string) {
	ic.cache[url] = true
}

func (lc *LinkCache) Has(url string) bool {
	_, ok := lc.cache[url]
	return ok
}

func (lc *LinkCache) Add(url string) {
	lc.cache[url] = true
}
