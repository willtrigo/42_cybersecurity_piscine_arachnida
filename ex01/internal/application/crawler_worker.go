// ************************************************************************** //
//                                                                            //
//                                                        :::      ::::::::   //
//   crawler_worker.go                                  :+:      :+:    :+:   //
//                                                    +:+ +:+         +:+     //
//   By: dande-je <dande-je@student.42sp.org.br>    +#+  +:+       +#+        //
//                                                +#+#+#+#+#+   +#+           //
//   Created: 2026/08/19 16:31:24 by dande-je          #+#    #+#             //
//   Updated: 2026/08/19 23:58:03 by dande-je         ###   ########.fr       //
//                                                                            //
// ************************************************************************** //

package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/willtrigo/42_cybersecurity_piscine_arachnida/ex01/internal/domain"
)

func (c *Crawler) SetWorkerCount(count int) {
	if count < c.workerCount {
		c.workerCount = count
	}
}

func (c *Crawler) CrawlWithWorkers(ctx context.Context, start *domain.URL) error {
	if err := c.checkContext(ctx, "crawl canceled before starting"); err != nil {
		return err
	}

	c.logger.Printf("spider: starting crawl at %s (max depth: %d)", start.String(), c.maxDepth)

	state := newCrawlState(start)
	errChan := make(chan error, c.workerCount)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		state.cond.Broadcast()
	}()

	var wg sync.WaitGroup
	for i := range c.workerCount {
		wg.Add(1)
		go c.worker(ctx, i, &wg, errChan, state, start, cancel)
	}

	wg.Wait()
	close(errChan)

	var firstError error
	for err := range errChan {
		if err != nil && firstError == nil {
			firstError = err
		}
	}

	if firstError != nil {
		return firstError
	}

	c.logger.Printf("spider: crawl completed: visited %d pages, downloaded %d images, errors: %d",
		c.stats.getPagesVisited(), c.stats.getImagesDownloaded(), c.stats.getErrors())

	return nil
}

func (c *Crawler) worker(
	ctx context.Context,
	id int,
	wg *sync.WaitGroup,
	errChan chan<- error,
	state *crawlState,
	start *domain.URL,
	cancel context.CancelFunc,
) {
	defer wg.Done()

	for {
		task, ok := state.nextTask(ctx)
		if !ok {
			return
		}

		err := c.processPage(ctx, task, state, start)
		state.taskDone()

		if err != nil {
			if !c.isContextError(err) {
				c.logger.Printf("spider: skipping %s: %v", task.url.String(), err)
				c.stats.incrementErrors()
				continue
			}

			errChan <- fmt.Errorf("worker %d: %w", id, err)
			cancel()
			return
		}
	}
}

func (c *Crawler) downloadImagesConcurrently(ctx context.Context, pageURL *domain.URL, imageURLs []*domain.URL, depth int) {
	if len(imageURLs) == 0 {
		return
	}

	c.logger.Printf("spider: start download from %s (depth %d)", pageURL.String(), depth)

	var wg sync.WaitGroup
	var downloadedCount int

	imageWorkers := min(len(imageURLs), 3)

	workChan := make(chan *domain.URL, len(imageURLs))

	for range imageWorkers {
		wg.Go(func() {
			for imgURL := range workChan {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if loaded := c.imageCache.Has(imgURL.String()); loaded {
					continue
				}

				c.imageCache.Add(imgURL.String())

				if err := c.downloader.Download(ctx, imgURL); err != nil {
					c.logger.Printf("spider: failed to download \n %s\n %v", imgURL.String(), err)
					c.stats.incrementErrors()
				} else {
					c.stats.incrementImagesDownloaded()
					downloadedCount++
				}
			}
		})
	}

	for _, imgURL := range imageURLs {
		select {
		case workChan <- imgURL:
		case <-ctx.Done():
			close(workChan)
			wg.Wait()
			return
		}
	}
	close(workChan)
	wg.Wait()

	c.logger.Printf("spider: downloaded %d/%d images from %s (depth %d)", downloadedCount, len(imageURLs), pageURL.String(), depth)
}
