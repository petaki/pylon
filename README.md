# Pylon

A small CLI based bookmark manager.

```
pylon https://github.com

=> Get User-Agent and WebSocket Debugger URL
=> Get URL content
=> Parse HTML source
=> Send data to MeiliSearch
+-------------------------------------------------------------------------------------------------------------------+
| ID                                   | URL                | Title                                                 |
+-------------------------------------------------------------------------------------------------------------------+
| 7e6a8e11-09da-4e59-a199-8589aefdbfc0 | https://github.com | GitHub · Change is constant. GitHub keeps you ahead.  |
+-------------------------------------------------------------------------------------------------------------------+
```

## Badges

[![Build Status](https://github.com/petaki/pylon/workflows/tests/badge.svg)](https://github.com/petaki/pylon/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## Features

- **Add** - save bookmarks with automatic title and metadata extraction
- **Tags** - organize bookmarks with custom tags
- **Search** - full-text search powered by MeiliSearch
- **Update** - re-crawl and refresh bookmark metadata
- **Delete** - remove individual or all bookmarks
- **Headless Crawling** - extract page content using chromedp/headless-shell

## Getting Started

Follow the steps below to install and configure Pylon.

### Prerequisites

- MeiliSearch: `Version >= 1.38` for data storage
- chromedp/headless-shell: `Version >= 145.0` for data crawling

### Install from Binary

Download the latest release for your platform from the [GitHub Releases](https://github.com/petaki/pylon/releases) page.

---

### Install from Source

#### Prerequisites

- Go: `Version >= 1.26`

#### Steps

1. Clone the repository:

```bash
git clone git@github.com:petaki/pylon.git
```

2. Build the binary:

```bash
cd pylon
go build
```

## Configuration

Initialize the `.pylonfile` file in your home directory or use environment variables:

```bash
pylon config init
```

### Headless Shell Host

```
HEADLESS_SHELL_HOST=http://127.0.0.1:9222
```

### MeiliSearch Host

```
MEILISEARCH_HOST=http://127.0.0.1:7700
```

### MeiliSearch API Key

```
MEILISEARCH_API_KEY=
```

### MeiliSearch Index

```
MEILISEARCH_INDEX=pylon
```

## Usage

The following commands show how to use the package.

### Add a link

```bash
pylon https://github.com/petaki
```

Or add with tags:

```bash
pylon --tags="go,code" https://golang.org
```

### Search links

```bash
pylon link search <query>
```

### Update the link

```bash
pylon link update <id>
```

### Delete the link

```bash
pylon link delete <id>
```

### Delete all links

```bash
pylon link delete-all
```

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/pylon/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
