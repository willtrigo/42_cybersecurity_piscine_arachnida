// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 17:54:45 by dande-je          #+#    #+#             //
//   Updated: 2026/08/18 00:11:13 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync/atomic"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

type Crawler struct {
	httpClient domain.HTTPClient
	parser     domain.HTMLParser
	downloader *Downloader
	logger     *log.Logger
	maxDepth   int
	stats      atomicCrawlStats
}

type atomicCrawlStats struct {
	pagesVisited     atomic.Uint64
	imagesDownloaded atomic.Uint64
	errors           atomic.Uint64
}

type crawlTask struct {
	url   *domain.URL
	depth int
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
	}
	return c
}

func (c *Crawler) Crawler(ctx context.Context, start *domain.URL) error {
	if err := checkContext(ctx); err != nil {
		return fmt.Errorf("crawl canceled before starting: %w", ctx.Err())
	}

	c.logger.Printf("spider: starting crawl at %s (max depth: %d)",
		start.String(), c.maxDepth)

	visited := map[string]bool{start.String(): true}
	queue := []crawlTask{{url: start, depth: 0}}
	visitedStart := false

	var processed int

	for len(queue) > 0 {
		if err := checkContext(ctx); err != nil {
			c.logger.Printf("spider: crawl interrupted after processing %d pages: %v", processed, ctx.Err())
			return fmt.Errorf("crawl interrupted: %w", ctx.Err())
		}

		task := queue[0]
		queue = queue[1:]
		processed++

		if task.depth > c.maxDepth {
			continue
		}

		if processed%10 == 0 {
			c.logger.Printf("spider: processed %d pages, queue size: %d, depth:%d", processed, len(queue), task.depth)
		}

		body, err := c.visitPage(ctx, task.url)
		if err != nil {
			if !visitedStart {
				return fmt.Errorf("crawling starting page %s: %w", task.url.String(), err)
			}

			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("crawl canceled: %w", err)
			}
			c.logger.Printf("spider: skipping %s: %v", task.url.String(), err)
			c.stats.errors.Add(1)
			continue
		}
		visitedStart = true
		c.stats.pagesVisited.Add(1)

		if task.depth < c.maxDepth {
			if err := checkContext(ctx); err != nil {
				return fmt.Errorf("crawl interrupted: %w", ctx.Err())
			}

			links := c.discoverLinks(ctx, task.url, body)
			for _, link := range links {
				if err := checkContext(ctx); err != nil {
					return fmt.Errorf("crawl interrupted: %w", ctx.Err())
				}

				if visited[link.String()] {
					continue
				}

				if link.Host() != start.Host() {
					continue
				}

				visited[link.String()] = true
				queue = append(queue, crawlTask{url: link, depth: task.depth + 1})
			}
		}
	}

	c.logger.Printf("spider: crawl completed: visited %d pages, downloaded %d images, errors: %d",
		c.stats.pagesVisited.Load(), c.stats.imagesDownloaded.Load(), c.stats.errors.Load())

	return nil
}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (c *Crawler) visitPage(ctx context.Context, pageURL *domain.URL) ([]byte, error) {
	if err := checkContext(ctx); err != nil {
		return nil, fmt.Errorf("visit canceled: %w", ctx.Err())
	}

	body, contentType, err := c.httpClient.Get(ctx, pageURL)
	if err != nil {
		return nil, fmt.Errorf("fetching page: %w", err)
	}

	if !isHTMLContent(contentType) {
		c.logger.Printf("spider: skipping non-HTML content from %s(content-type: %s", pageURL.String(), contentType)
		return body, nil
	}

	imagesURLs, err := c.parser.ExtractImagesURLs(pageURL, body)
	if err != nil {
		return nil, fmt.Errorf("extraction images: %w", err)
	}

	for i := range imagesURLs {
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("download interrupted: %w", ctx.Err())
		}

		// TODO: logic to download image

		if (i+1)%10 == 0 {
			c.logger.Printf("spider: downloaded %d/%d images from %s", i+1, len(imagesURLs), pageURL.String())
		}
	}

	return body, nil
}

func (c *Crawler) discoverLinks(ctx context.Context, pageURL *domain.URL, body []byte) []*domain.URL {
	if err := checkContext(ctx); err != nil {
		c.logger.Printf("spider: link discovery canceled for %s", pageURL.String())
	}

	links, err := c.parser.ExtractLinks(pageURL, body)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			c.logger.Printf("spider: failed to extract links from %s: %v", pageURL.String(), err)
			c.stats.errors.Add(1)
		}
		return nil
	}

	return links
}

func isHTMLContent(contentType string) bool {
	return contentType == "" ||
		contentType == "text/html" ||
		contentType == "application/xhtml+xml" ||
		strings.HasPrefix(contentType, "text/html;") ||
		strings.HasPrefix(contentType, "application/xhtml+xml;")
}
