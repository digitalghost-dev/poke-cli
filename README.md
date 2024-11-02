<div align="center">
    <img height="250" width="350" src="https://cdn.simpleicons.org/pokemon/FFCC00" alt="pokemon-logo"/>
    <h1>Pokémon CLI</h1>
    <img src="https://img.shields.io/github/v/release/digitalghost-dev/poke-cli?style=flat-square&logo=git&logoColor=FFCC00&label=Release%20Version&labelColor=EEE&color=FFCC00" alt="version-label">
    <img src="https://img.shields.io/docker/image-size/digitalghostdev/poke-cli/v0.7.0?arch=arm64&style=flat-square&logo=docker&logoColor=FFCC00&labelColor=EEE&color=FFCC00" alt="docker-image-size">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/ci.yml?branch=main&style=flat-square&logo=github&logoColor=FFCC00&label=CI&labelColor=EEE&color=FFCC00">
</div>

<div align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/go_test.yml?style=flat-square&logo=go&logoColor=00ADD8&label=Tests&labelColor=EEE&color=00ADD8" alt="tests-label">
    <img src="https://img.shields.io/github/go-mod/go-version/digitalghost-dev/poke-cli?style=flat-square&logo=Go&labelColor=EEE&color=00ADD8" alt="go-version">
   
</div>

## Overview
A CLI tool for viewing data about Pokémon from your terminal!

## Demo
![demo](https://pokemon-objects.nyc3.digitaloceanspaces.com/demo.gif)

## Install

### Taskfile
_Taskfile can build the executable for you_

1. Install [Taskfile](https://taskfile.dev/installation/).
2. Once installed, clone the repository and `cd` into it.
3. Then, simply run `task build` and an executable for your machine type will be created. 
    * Example usage:
   ```bash
   # Windows
   .\poke-cli.exe pokemon charizard --types --abilities
   
   # Unix
   .\poke-cli pokemon vespiquen -t -a
   ```


### Docker
_Use a Docker Image_

```bash
docker run --rm -it digitalghostdev/poke-cli:v0.7.0 [command] [subcommand] [flag]
```

### Go Build
_Build the executable yourself_

1. Install [Golang](https://go.dev/dl/).
2. Once installed, clone the repository and `cd` into it.
3. Run `go build .` or `go install` to build the executable.
    * Keep in mind that `go install` will place the executable in your `$GOPATH/bin` directory. [Read More](https://www.golang.company/blog/what-is-the-difference-between-go-run-go-build-and-go-install-commands).
4. An executable will be created then the tool can be used! It can also be added to your path to run the binary from anywhere.
   * Example usage:
   ```bash
   # Windows
   .\poke-cli.exe pokemon charizard --types --abilities
   
   # Unix
   .\poke-cli pokemon vespiquen -t -a
   
   # If built with go install
   poke-cli pokemon slugma -t
   ```

## Usage
By running `poke-cli [-h | --help]`, it'll display information on how to use the tool. 
```
╭──────────────────────────────────────────────────────╮
│Welcome! This tool displays data related to Pokémon!  │
│                                                      │
│ USAGE:                                               │
│    poke-cli [flag]                                   │
│    poke-cli [command] [flag]                         │
│    poke-cli [command] [subcommand] [flag]            │
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
