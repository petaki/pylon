# Pylon

A small CLI based bookmark manager.

## Badges

[![Build Status](https://github.com/petaki/pylon/workflows/tests/badge.svg)](https://github.com/petaki/pylon/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## Getting Started

Follow the steps below to install and configure Pylon.

### Prerequisites

- MeiliSearch: `Version >= 1.24` for data storage
- chromedp/headless-shell: `Version >= 141.0` for data crawling

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
