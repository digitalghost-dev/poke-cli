<div align="center">
    <img height="250" width="350" src="https://cdn.simpleicons.org/pokemon/FFCC00" alt="pokemon-logo"/>
    <h1>Pokémon CLI</h1>
    <img src="https://img.shields.io/github/v/release/digitalghost-dev/poke-cli?style=flat-square&logo=git&logoColor=FFCC00&label=Release%20Version&labelColor=EEE&color=FFCC00" alt="version-label">
    <img src="https://img.shields.io/docker/image-size/digitalghostdev/poke-cli/v0.7.2?arch=arm64&style=flat-square&logo=docker&logoColor=FFCC00&labelColor=EEE&color=FFCC00" alt="docker-image-size">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/ci.yml?branch=main&style=flat-square&logo=github&logoColor=FFCC00&label=CI&labelColor=EEE&color=FFCC00">
</div>

<div align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/go_test.yml?style=flat-square&logo=go&logoColor=00ADD8&label=Tests&labelColor=EEE&color=00ADD8" alt="tests-label">
    <img src="https://img.shields.io/github/go-mod/go-version/digitalghost-dev/poke-cli?style=flat-square&logo=Go&labelColor=EEE&color=00ADD8" alt="go-version">
   
</div>

## Overview
A CLI tool for viewing data about Pokémon from your terminal! I am new to writing Go and taking my time in building this 
project. 

My aim is to have four commands finished for `v1.0.0`. Read more in the [Roadmap](#roadmap) section.

---
## Demo
![demo](https://pokemon-objects.nyc3.digitaloceanspaces.com/demo_0.7.1.gif)

---
## Install

### Docker Image
_Use a Docker Image_

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/).
2. Once installed, use the command below to pull the image and run the container!
   * `--rm`: Automatically remove the container when it exits. 
     * Optional.
   * `-i`: Interactive mode, keeps STDIN open for input.
     * Necessary.
   * `-t`: Allocates a terminal (TTY) for a terminal-like session.
     * Necessary.

```bash
docker run --rm -i -t digitalghostdev/poke-cli:v0.7.2 <command> [subcommand] flag]
```

### Go Install
_Install the executable yourself_

1. Install [Golang](https://go.dev/dl/).
2. Once installed, run the following command:
   ```bash
   go install github.com/digitalghost-dev/poke-cli@v0
   ```
3. The tool is ready to use!
---
## Usage
By running `poke-cli [-h | --help]`, it'll display information on how to use the tool. 
```
╭──────────────────────────────────────────────────────╮
│Welcome! This tool displays data related to Pokémon!  │
│                                                      │
│ USAGE:                                               │
│    poke-cli [flag]                                   │
│    poke-cli <command> [flag]                         │
│    poke-cli <command> <subcommand> [flag]            │
│                                                      │
│ FLAGS:                                               │
│    -h, --help      Shows the help menu               │
│    -l, --latest    Prints the latest available       │
│                    version of the program            │
│                                                      │
│ AVAILABLE COMMANDS:                                  │
│    pokemon         Get details of a specific Pokémon │
│    types           Get details of a specific typing  │
╰──────────────────────────────────────────────────────╯
```

---
## Roadmap
The architecture behind how the tool works is straight forward.
1. Commands indicate which data endpoint to focus on.
2. Flags provide more information and can be all stacked together or chosen.

### Planned for Version 1.0.0
- [ ] `pokemon`: get data about a specific Pokémon.
   - [x] `--abilities | -a`: display the Pokémon's abilities.
   - [x] `--types | -t`: display the Pokémon's typing.
   - [ ] `--stats | -s`: 
- [x] `types`: get data about a specific typing.
- [ ] `ability`: get data about a specific ability.
- [ ] `move`: get data about a specific move.