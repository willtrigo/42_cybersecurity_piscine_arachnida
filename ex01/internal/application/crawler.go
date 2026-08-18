// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 17:54:45 by dande-je          #+#    #+#             //
//   Updated: 2026/08/18 14:32:56 by dande-je         ###   ########.fr       //
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
	httpClient domain.HTTPClient
	parser     domain.HTMLParser
	downloader *Downloader
	logger     *log.Logger
	stats      *crawlStats
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
		logger:     log.Default(),
		stats:      newCrawlStats(),
	}
	return c
}

func (c *Crawler) Crawler(ctx context.Context, start *domain.URL) error {
	if err := c.checkContext(ctx, "crawl canceled before starting"); err != nil {
		return err
	}

	c.logger.Printf("spider: starting crawl at %s (max depth: %d)", start.String(), c.maxDepth)

	state := newCrawlState(start)
	var processed int

	for state.hasMoreTasks() {
		if err := c.checkContext(ctx, "crawl interrupted"); err != nil {
			c.logger.Printf("spider: crawl interrupted after processing %d pages: %v", processed, ctx.Err())
			return err
		}

		task := state.nextTask()
		processed++

		if task.depth > c.maxDepth {
			continue
		}

		c.logProgess(processed, state, task.depth)

		if err := c.processPage(ctx, task, state, start); err != nil {
			if state.hasVisitedStart() && !c.isContextError(err) {
				c.logger.Printf("spider: skipping %s: %v", task.url.String(), err)
				c.stats.incrementErrors()
				continue
			}
			return err
		}

		state.markPageVisited()
		c.stats.incrementPagesVisited()
	}

	c.logger.Printf("spider: crawl completed: visited %d pages, downloaded %d images, errors: %d",
		c.stats.getPagesVisited(), c.stats.getImagesDownloaded(), c.stats.getErrors())

	return nil
}

func (c *Crawler) processPage(ctx context.Context, task crawlTask, state *crawlState, start *domain.URL) error {
	body, err := c.fetchPage(ctx, task.url, task)
	if err != nil {
		return err
	}

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

	c.downloadImages(ctx, url, imageURLs, task.depth)

	return body, nil
}

func (c *Crawler) downloadImages(ctx context.Context, pageURL *domain.URL, imageURLs []*domain.URL, depth int) {
	var i int

	for i = range imageURLs {
		if err := c.checkContext(ctx, "download interrupted"); err != nil {
			return
		}

		// TODO: logic to download image
		c.stats.incrementImagesDownloaded()

	}

	downloads := 0
	if i != 0 {
		downloads = i + 1
	}
	c.logger.Printf("spider: downloaded %d/%d images from %s (depth %d)", downloads, len(imageURLs), pageURL.String(), depth)
}

func (c *Crawler) enqueueLinks(ctx context.Context, task crawlTask, body []byte, state *crawlState, start *domain.URL) error {
	links := c.discoverLinks(ctx, task.url, body)

	nextDepth := task.depth + 1
	for _, link := range links {
		if err := c.checkContext(ctx, "crawl interrupted"); err != nil {
			return err
		}

		if state.isVisited(link) || link.Host() != start.Host() {
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
