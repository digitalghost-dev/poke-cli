<div align="center">
    <img height="250" width="350" src="pokemon.svg" alt="pokemon-logo"/>
    <h1>Pokémon CLI</h1>
    <img src="https://img.shields.io/github/v/release/digitalghost-dev/poke-cli?style=flat-square&logo=git&logoColor=FFCC00&label=Release%20Version&labelColor=EEE&color=FFCC00" alt="version-label">
    <img src="https://img.shields.io/docker/image-size/digitalghostdev/poke-cli/v1.4.0?arch=arm64&style=flat-square&logo=docker&logoColor=FFCC00&labelColor=EEE&color=FFCC00" alt="docker-image-size">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/ci.yml?branch=main&style=flat-square&logo=github&logoColor=FFCC00&label=CI&labelColor=EEE&color=FFCC00" alt="ci-status-badge">
</div>
<div align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/go_test.yml?style=flat-square&logo=go&logoColor=00ADD8&label=Tests&labelColor=EEE&color=00ADD8" alt="tests-label">
    <img src="https://img.shields.io/github/go-mod/go-version/digitalghost-dev/poke-cli?style=flat-square&logo=Go&labelColor=EEE&color=00ADD8" alt="go-version"/>
    <img src="https://img.shields.io/codecov/c/github/digitalghost-dev/poke-cli?token=05GBSAOQIT&style=flat-square&logo=codecov&logoColor=00ADD8&labelColor=EEE&color=00ADD8" alt="codecov"/>
</div>

## Overview
`poke-cli` is a hybrid of a classic CLI and a modern TUI tool for viewing data about Pokémon! This is my first Go project.
View the [documentation](https://docs.poke-cli.com)!

The architecture behind how the tool works is straight forward:
1. Each command indicates which [API](https://pokeapi.co/) endpoint to use.
2. Flags provide more information and can be stacked together or used individually.
3. Each command has a `-h | --help` flag that is built-in with Golang's `flag` package.

View future plans in the [Roadmap](#roadmap) section.

---
## Demo
![demo](https://poke-cli-s3-bucket.s3.us-west-2.amazonaws.com/demo-v1.3.3.gif)

---
## Installation

* [Binary](#binary)
* [Docker Image](#docker-image)
* [Homebrew](#homebrew)
* [Source](#source)

### Binary

1. Head to the [releases](https://github.com/digitalghost-dev/poke-cli/releases) page of the project.
2. Choose a version to download. The latest is best.
3. Choose an operating system and click on the matching zipped folder to start the download.
4. Extract the folder. The tool is ready to use.
5. Either change directories into the extracted folder or move the binary to a chosen directory.
6. Run the tool!

> [!IMPORTANT]
> For macOS, you may have to allow the executable to run as it is not signed. Head to System Settings > Privacy & Security > scroll down and allow executable to run.

<details>

<summary>View Image of Settings</summary>

![settings](https://poke-cli-s3-bucket.s3.us-west-2.amazonaws.com/macos_privacy_settings.png)

</details>


 #### Example usage
  ```bash
  # Windows
  .\poke-cli.exe pokemon charizard --types --abilities
   
  # Unix
  .\poke-cli ability airlock --pokemon
  ```

### Docker Image

1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/).
2. Once installed, use the command below to pull the image and run the container!
   * `--rm`: Automatically remove the container when it exits. 
     * Optional.
   * `-i`: Interactive mode, keeps STDIN open for input.
     * Necessary.
   * `-t`: Allocates a terminal (TTY) for a terminal-like session.
     * Necessary.
3. Choose how to interact with the container:
   * Run a single command and exit:
    ```bash
    docker run --rm -it digitalghostdev/poke-cli:v1.4.0 <command> [subcommand] flag]
    ```
   * Enter the container and use its shell:
    ```bash
    docker run --rm -it --name poke-cli --entrypoint /bin/sh digitalghostdev/poke-cli:v1.4.0 -c "cd /app && exec sh"
   # placed into the /app directory, run the program with './poke-cli'
   # example: ./poke-cli ability swift-swim
    ```
   
### Homebrew
1. Install the Cask:
    ```bash
    brew install --cask digitalghost-dev/tap/poke-cli
    ````
2. Verify install!
    ```bash
    poke-cli -v
    ```

### Source

1. Run the following command:
   ```bash
   go install github.com/digitalghost-dev/poke-cli@latest
   ```
2. The tool is ready to use!

---
## Usage
By running `poke-cli [-h | --help]`, it'll display information on how to use the tool. 
```
╭──────────────────────────────────────────────────────────╮
│Welcome! This tool displays data related to Pokémon!      │
│                                                          │
│ USAGE:                                                   │
│    poke-cli [flag]                                       │
│    poke-cli <command> [flag]                             │
│    poke-cli <command> <subcommand> [flag]                │
│                                                          │
│ FLAGS:                                                   │
│    -h, --help      Shows the help menu                   │
│    -l, --latest    Prints the latest version available   │
│    -v, --version   Prints the current version            │
│                                                          │
│ COMMANDS:                                                │
│    ability         Get details about an ability          │
│    item            Get details about an item             │
│    move            Get details about a move              │
│    natures         Get details about all natures         │
│    pokemon         Get details about a Pokémon           │
│    search          Search for a resource                 │
│    types           Get details about a typing            │
│                                                          │
│ hint: when calling a resource with a space, use a hyphen │
│ example: poke-cli ability strong-jaw                     │
│ example: poke-cli pokemon flutter-mane                   │
│                                                          │
│ ↓ ctrl/cmd + click for docs/guides                       │
│ docs.poke-cli.com                                        │
╰──────────────────────────────────────────────────────────╯
```

---

## Roadmap
Below is a list of the planned/completed commands and flags:

- [x] `ability`: get data about an ability.
    - [x] `-p | --pokemon`: display Pokémon that learn this ability.
- [ ] `berry`: get data about a berry.
- [x] `item`: get data about an item.
- [x] `move`: get data about a move.
    - [ ] `-p | --pokemon`: display Pokémon that learn this move.
- [x] `natures`: get data about natures.
- [x] `pokemon`: get data about a Pokémon.
    - [x] `-a | --abilities`: display the Pokémon's abilities.
    - [x] `-i | --image`: display a pixel image of the Pokémon.
    - [x] `-s | --stats`: display the Pokémon's base stats.
    - [x] `-t | --types`: display the Pokémon's typing.
    - [x] `-m | --moves`: display learnable moves.
- [ ] `search`: search for a resource 
    - [x] `ability`
    - [ ] `berry`
    - [ ] `item`
    - [x] `move`
    - [x] `pokemon`
- [ ] `speed`: compare speed stats between two Pokémon.
- [x] `types`: get data about a specific typing.

---
## Tested Terminals
| OS      | Terminal          | Status | Issues                        |
|---------|-------------------|:------:|-------------------------------|
| macOS   | Ghostty           |   ✅    | None                          | 
| macOS   | Alacritty         |   ✅    | None                          |
| macOS   | HyperJS           |   ✅    | None                          |
| macOS   | iTerm2            |   ✅    | None                          |
| macOS   | macOS Terminal    |   ⚠️   | Images do not render properly |
| Windows | Windows Terminal  |   ✅    | None                          |
| Ubuntu  | Standard Terminal |   ✅    | None                          |
| Ubuntu  | Tabby             |   ✅    | None                          |