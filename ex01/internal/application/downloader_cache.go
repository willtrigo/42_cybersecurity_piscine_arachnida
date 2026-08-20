// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   downloader_cache.go                                :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/18 22:43:16 by dande-je          #+#    #+#             //
//   Updated: 2026/08/19 19:06:13 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import "sync"

type ImageCache struct {
	cache map[string]bool
	mu    sync.RWMutex
}

type LinkCache struct {
	cache map[string]bool
	mu    sync.RWMutex
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
	ic.mu.RLock()
	defer ic.mu.RUnlock()
	_, ok := ic.cache[url]
	return ok
}

func (ic *ImageCache) Add(url string) {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	ic.cache[url] = true
}

func (lc *LinkCache) Has(url string) bool {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	_, ok := lc.cache[url]
	return ok
}

func (lc *LinkCache) Add(url string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.cache[url] = true
}
