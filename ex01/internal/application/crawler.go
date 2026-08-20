// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 17:54:45 by dande-je          #+#    #+#             //
//   Updated: 2026/08/21 11:07:16 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"context"
	"fmt"
	"log"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

type Crawler struct {
	httpClient  domain.HTTPClient
	parser      domain.HTMLParser
	downloader  *Downloader
	logger      *log.Logger
	stats       *CrawlStats
	imageCache  *ImageCache
	linkCache   *LinkCache
	maxDepth    int
	workerCount int
}

func NewCrawler(
	httpClient domain.HTTPClient,
	parser domain.HTMLParser,
	downloader *Downloader,
	maxDepth int,
) *Crawler {
	c := &Crawler{
		httpClient:  httpClient,
		parser:      parser,
		downloader:  downloader,
		maxDepth:    maxDepth,
		logger:      log.Default(),
		stats:       newCrawlStats(),
		imageCache:  newImageCache(),
		linkCache:   newLinkCache(),
		workerCount: 10,
	}
	return c
}

func (c *Crawler) Crawl(ctx context.Context, start *domain.URL) error {
	fmt.Printf("spider: using %d concurrent workers\n", c.workerCount)
	return c.CrawlWithWorkers(ctx, start)
}

func (c *Crawler) processPage(ctx context.Context, task crawlTask, state *crawlState, start *domain.URL) error {
	body, err := c.fetchPage(ctx, task.url, task)
	if err != nil {
		return err
	}

	c.stats.incrementPagesVisited()

	if task.depth < c.maxDepth {
		return c.enqueueLinks(ctx, task, body, state, start)
	}

	return nil
}

func (c *Crawler) fetchPage(ctx context.Context, url *domain.URL, task crawlTask) ([]byte, error) {
	body, contentType, err := c.httpClient.Get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetching page: %w", err)
	}

	if !isHTMLContent(contentType) {
		c.logger.Printf("spider: skipping non-HTML content from %s(content-type: %s", url.String(), contentType)
		return body, nil
	}

	imageURLs, err := c.parser.ExtractImagesURLs(url, body)
	if err != nil {
		return nil, fmt.Errorf("extraction images: %w", err)
	}

	c.downloadImagesConcurrently(ctx, url, imageURLs, task.depth)

	return body, nil
}

func (c *Crawler) enqueueLinks(ctx context.Context, task crawlTask, body []byte, state *crawlState, start *domain.URL) error {
	links := c.discoverLinks(ctx, task.url, body)

	nextDepth := task.depth + 1
	for _, link := range links {
		if err := c.checkContext(ctx, "crawl interrupted"); err != nil {
			return err
		}

		if link.Host() != start.Host() {
			continue
		}

		state.addTask(link, nextDepth)
	}

	return nil
}

func (c *Crawler) discoverLinks(ctx context.Context, pageURL *domain.URL, body []byte) []*domain.URL {
	if err := c.checkContext(ctx, "link discover canceled"); err != nil {
		c.logger.Printf("spider: link discovery canceled for %s", pageURL.String())
		return nil
	}

	if visited := c.linkCache.Has(pageURL.String()); visited {
		return nil
	}

	c.linkCache.Add(pageURL.String())

	links, err := c.parser.ExtractLinks(pageURL, body)
	if err != nil {
		if !c.isContextError(err) {
			c.logger.Printf("spider: failed to extract links from %s: %v", pageURL.String(), err)
			c.stats.incrementErrors()
		}
		return nil
	}

	return links
}
