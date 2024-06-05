<div align="center">
    <img height="250" width="350" src="https://cdn.simpleicons.org/pokemon/FFCC00" alt="pokemon-logo"/>
    <h1>Pokémon CLI</h1>
    <img src="https://img.shields.io/github/v/release/digitalghost-dev/poke-cli?style=flat-square&logo=git&logoColor=FFCC00&label=Release%20Version&labelColor=EEE&color=FFCC00" alt="version-label">
    <img src="https://img.shields.io/docker/image-size/digitalghostdev/poke-cli/v0.1.1?arch=arm64&style=flat-square&logo=docker&logoColor=FFCC00&labelColor=EEE&color=FFCC00" alt="docker-image-size">
</div>

<div align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/go_test.yml?style=flat-square&logo=go&logoColor=00ADD8&label=Tests&labelColor=EEE&color=00ADD8" alt="tests-label">
    <img src="https://img.shields.io/github/go-mod/go-version/digitalghost-dev/poke-cli?style=flat-square&logo=Go&labelColor=EEE&color=00ADD8" alt="go-version">
</div>

## Overview
A CLI tool for viewing data about Pokémon from your terminal!

## Install

### Go Build
1. Make sure [Go is installed](https://go.dev/dl/) on your machine. This project uses `v1.21`.
2. Clone the repository in a root directory: `git clone https://github.com/digitalghost-dev/poke-cli.git`
3. Change directories into the `poke-cli` directory.
4. Run `go build -o poke-cli`
5. A binary will be created then the tool can be used! It can also be added to your path to run the binary from anywhere.

### Docker
Use a Docker Image instead:
```bash
docker run --rm -it digitalghostdev/poke-cli:0.1.1 [command] [flag]
```

> [!NOTE]
> Currently working on more ways to distribute the binary.

## Usage
By running `poke-cli --help`, it'll display information on how to use the tool. 
```
Welcome! This tool displays data about a selected Pokémon in the terminal!
      
USAGE:
         poke-cli [flag]
         poke-cli [pokemon name] [flag]
         ----------
         Example: poke-cli bulbasaur or poke-cli flutter-mane --types
             
GLOBAL FLAGS:
         -h, --help      Shows the help menu

POKEMON NAME FLAGS:
         Add a flag after declaring a Pokémon's name for more details!
        --types
```
