// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler.go                                         :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/17 17:54:45 by dande-je          #+#    #+#             //
//   Updated: 2026/08/18 11:16:46 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

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
	if err := checkContext(ctx); err != nil {
		return fmt.Errorf("crawl canceled before starting: %w", ctx.Err())
	}

	c.logger.Printf("spider: starting crawl at %s (max depth: %d)",
		start.String(), c.maxDepth)

	state := newCrawlState(start)
	var processed int

	for state.hasMoreTasks() {
		if err := checkContext(ctx); err != nil {
			c.logger.Printf("spider: crawl interrupted after processing %d pages: %v", processed, ctx.Err())
			return fmt.Errorf("crawl interrupted: %w", ctx.Err())
		}

		task := state.nextTask()
		processed++

		if task.depth > c.maxDepth {
			continue
		}

		if processed%10 == 0 {
			c.logger.Printf("spider: processed %d pages, queue size: %d, depth:%d", processed, state.queueSize(), state.currentDepth())
		}

		body, err := c.visitPage(ctx, task.url)
		if err != nil {
			if state.hasVisitedStart() {
				return fmt.Errorf("crawling starting page %s: %w", task.url.String(), err)
			}

			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return fmt.Errorf("crawl canceled: %w", err)
			}
			c.logger.Printf("spider: skipping %s: %v", task.url.String(), err)
			c.stats.incrementErrors()
			continue
		}
		state.markPageVisited()
		c.stats.incrementPagesVisited()

		if task.depth < c.maxDepth {
			if err := checkContext(ctx); err != nil {
				return fmt.Errorf("crawl interrupted: %w", ctx.Err())
			}

			links := c.discoverLinks(ctx, task.url, body)
			for _, link := range links {
				if err := checkContext(ctx); err != nil {
					return fmt.Errorf("crawl interrupted: %w", ctx.Err())
				}

				if state.isVisited(link) || link.Host() != start.Host() {
					continue
				}

				state.addTask(link, task.depth+1)
			}
		}
	}

	c.logger.Printf("spider: crawl completed: visited %d pages, downloaded %d images, errors: %d",
		c.stats.getPagesVisited(), c.stats.getImagesDownloaded(), c.stats.getErrors())

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
		c.stats.incrementImagesDownloaded()

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
			c.stats.incrementErrors()
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
