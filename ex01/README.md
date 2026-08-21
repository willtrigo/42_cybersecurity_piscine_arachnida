# 42 School — Cybersecurity Piscine — Arachnida — Exercise 01.

## Description
 
**Spider** is the first exercise of the Cybersecurity Piscine's Arachnida project. It's a command-line web scraper that downloads every image found on a given page and, optionally, recurses through the links it discovers to pull images from the rest of the site as well.
 
The program must:
- Accept a URL and download every image referenced on that page
- Recursively follow links to other pages on the same host when `-r` is set
- Cap how deep that recursion goes with `-l`
- Save downloaded images to a configurable output directory with `-p`
- Only download the extensions `.jpg`/`.jpeg`, `.png`, `.gif` and `.bmp`
No scraping libraries or `wget`-style shortcuts are allowed — HTTP requests and file handling can rely on the standard library, but the crawling, parsing and downloading logic all had to be written from scratch.
 
## Algorithm Explanation
 
The crawler is a **worker pool feeding on a shared BFS frontier**, structured as a small hexagonal/clean-architecture layout (`domain` → `application` → `adapter`):
 
- **Shared frontier (`crawlState`)**: a queue of `(url, depth)` tasks guarded by a `sync.Cond`. Workers call `nextTask` to block until either new work is enqueued or the crawl is provably finished — the queue is empty *and* no other worker still has a task in flight. This avoids a classic race where a naive producer loop would decide "no more work" the instant the queue looks empty, even though a worker is still mid-page and about to discover child links.
- **Deduplication**: adding a task and checking whether a URL was already visited happen as a single atomic step under one lock, so two workers discovering the same link from different pages at the same moment can't both win and queue a duplicate crawl of the same page.
- **URL normalization**: every URL (including the one passed on the command line) is normalized to a canonical form — fragment stripped, trailing slash collapsed — before it's used as a dedup key, so `https://example.com` and `https://example.com/` are recognized as the same page.
- **Host-scoped recursion**: only links whose host matches the starting URL's host are followed, keeping the crawl on the target site instead of wandering off to every subdomain or external site it happens to link to.
- **Per-page concurrent downloads**: images found on a page are downloaded through a small bounded worker pool of their own, deduplicated via an in-memory cache so the same image URL is never fetched twice.
- **Graceful shutdown**: the whole crawl runs under a `context.Context` cancelled on `SIGINT`/`SIGTERM`/`SIGQUIT`, so an interrupt stops in-flight work cleanly instead of leaving goroutines dangling.

## Instructions
 
The project is built and run entirely through Docker, orchestrated via `Makefile` (or the equivalent `Taskfile.yml`, if you prefer [Task](https://taskfile.dev) over `make`) — no local Go toolchain is required.
 
### Compilation
 
```sh
# Development: build (if needed) and start a hot-reloading container.
# Source is bind-mounted, so edits on the host rebuild via `air` automatically.
make up
```
```sh
# Production: build the minimal runtime image and drop into a shell inside it
make prod
```
```sh
# Rebuild all images from scratch (parallel, no cache, latest base images)
make build
```
```sh
# Static analysis (golangci-lint), run in an isolated container
make lint
```
```sh
# Test suite with race detection and coverage
make test
```
 
Every target above has a matching `task` equivalent (`task up`, `task prod`, `task build`, `task lint`, `task test`) — run `make help` or `task --list` to see the full set, including `shell`, `logs`, `down`, `clean`/`fclean` and `re`.
 
## Usage Example
 
```sh
./spider [-r] [-l N] [-p PATH] URL
```
 
| Flag | Description | Default |
|------|-------------|---------|
| `-r` | Recursively follow links found on the page | off |
| `-l N` | Maximum recursion depth (requires `-r`) | `5` |
| `-p PATH` | Directory downloaded images are saved to | `./data/` |
 
```sh
# Download every image on a single page
./spider https://example.com
```
```sh
# Recursively crawl up to depth 3, saving into ./images/
./spider -r -l 3 -p ./images/ https://example.com
```
 
Once you're inside the container (`make shell` for the dev image, or `make prod` for the production one), the binary is available at `/app/spider` (`./spider` from that working directory).
 
## Resources
 
- [Go `flag` package](https://pkg.go.dev/flag)
- [Go `sync.Cond`](https://pkg.go.dev/sync#Cond)
- [Go `net/http` client](https://pkg.go.dev/net/http)
- [Docker multi-stage builds](https://docs.docker.com/build/building/multi-stage/)
- [golangci-lint](https://golangci-lint.run/)
- [air - Live reload for Go apps](https://github.com/air-verse/air)

## AI Usage
 
AI (Claude) was used in the following parts:
 
- Initial project structure setup and architecture design
- Refactoring the name of some variables for better structure and readability
- Improving and updating this `README.md`
