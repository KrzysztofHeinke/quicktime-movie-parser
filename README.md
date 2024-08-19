# QuickTime Movie Parser

## Overview

The **QuickTime Movie Parser** is a powerful command-line tool designed for parsing and analyzing MOV/MP4 files. The tool focuses on extracting metadata, such as track information, sample rates, video dimensions, and more, from the `moov` atom and its child atoms. This utility is particularly useful for developers, media professionals, and anyone involved in processing and analyzing media files.

## Features

- **Atom Parsing:** Parse and analyze the `moov` atom, as well as its child atoms (`trak`, `mdia`, `minf`, `stbl`, etc.), to extract detailed metadata.
- **Track Information Extraction:** Extract information about video and audio tracks, including width, height, and sample rates.
- **Customizable Search:** Search for specific atoms within a file and analyze their contents.

### Prerequisites

- [Go](https://golang.org/doc/install) 1.16 or later
- [Docker](https://docs.docker.com/engine/install/)

### Clone the Repository

```bash
git clone https://github.com/KrzysztofHeinke/quicktime-movie-parser.git
cd quicktime-movie-parser
```

### Build the Project

If you want build it on Windows machine better to build it locally but it needs to have go installed.
```bash
make build-local
```

That command will build it in docker but only for linux. If you have windows you can always run it in WSL.
```bash
make build
```

### Run the tests

```bash
make test
```

```bash
make test-docker
```

### Run the linter

```bash
make lint
```

```bash
make lint-docker
```

### Usage

```bash
./bin/quicktime-movie-parser parse example.mov
```

This command parses the moov atom in the provided example.mov file and prints out the extracted metadata, such as track information, sample rates, and video dimensions.

By default loggin is set to INFO. 
If you would like to take a look on DEBUG messages

```bash
./bin/quicktime-movie-parser --loglevel=debug parse example.mov
```

### License

This project is licensed under the MIT License. See the LICENSE file for details.
Contact

If you have any questions, issues, or suggestions, feel free to open an issue on GitHub or contact the project maintainer at krzysztof.heinke@gmail.com.