# Pylon

[![Build Status](https://github.com/petaki/pylon/workflows/tests/badge.svg)](https://github.com/petaki/pylon/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

A small CLI based bookmark manager.

## Getting Started

Before you start, you need to install the prerequisites.

### Prerequisites

- MeiliSearch: `Version >= 0.28` for data storage
- chromedp/headless-shell: `Version >= 90.0` for data crawling

### Install from binary

Downloads can be found at releases page on [GitHub](https://github.com/petaki/pylon/releases).

---

### Install from source

#### Prerequisites for building

- GO: `Version >= 1.19`

#### 1. Clone the repository:

```
git clone git@github.com:petaki/pylon.git
```

#### 2. Open the folder:

```
cd pylon
```

#### 3. Build the Pylon:

```
go build
```

## Configuration

Initialize the `.pylonfile` file in your home directory or use environment variables:

```
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

```
pylon https://github.com/petaki
```

Or add with tags:

```
pylon --tags="go,code" https://golang.org
```

### Search links

```
pylon link search <query>
```

### Update the link

```
pylon link update <id>
```

### Delete the link

```
pylon link delete <id>
```

### Delete all links

```
pylon link delete-all
```

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/pylon/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
