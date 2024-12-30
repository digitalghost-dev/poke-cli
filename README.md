<div align="center">
    <img height="250" width="350" src="https://cdn.simpleicons.org/pokemon/FFCC00" alt="pokemon-logo"/>
    <h1>Pokémon CLI</h1>
    <img src="https://img.shields.io/github/v/release/digitalghost-dev/poke-cli?style=flat-square&logo=git&logoColor=FFCC00&label=Release%20Version&labelColor=EEE&color=FFCC00" alt="version-label">
    <img src="https://img.shields.io/docker/image-size/digitalghostdev/poke-cli/v0.11.0?arch=arm64&style=flat-square&logo=docker&logoColor=FFCC00&labelColor=EEE&color=FFCC00" alt="docker-image-size">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/ci.yml?branch=main&style=flat-square&logo=github&logoColor=FFCC00&label=CI&labelColor=EEE&color=FFCC00" alt="ci-status-badge">
</div>
<div align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/go_test.yml?style=flat-square&logo=go&logoColor=00ADD8&label=Tests&labelColor=EEE&color=00ADD8" alt="tests-label">
    <img src="https://img.shields.io/github/go-mod/go-version/digitalghost-dev/poke-cli?style=flat-square&logo=Go&labelColor=EEE&color=00ADD8" alt="go-version"/>
    <img src="https://img.shields.io/codecov/c/github/digitalghost-dev/poke-cli?token=05GBSAOQIT&style=flat-square&logo=codecov&logoColor=00ADD8&labelColor=EEE&color=00ADD8" alt="codecov"/>
</div>

## Overview
A CLI tool for viewing data about Pokémon from your terminal! I am new to writing Go and taking my time in building this 
project. 

My aim is to have five commands finished for `v1.0.0`. Read more in the [Roadmap](#roadmap) section.

---
## Demo
![demo](https://poke-cli-s3-bucket.s3.us-west-2.amazonaws.com/demo-v0.10.0.gif)

---
## Install

### Binary
_Download a pre-built binary_

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

![settings](https://pokemon-objects.nyc3.digitaloceanspaces.com/macos_settings.png)

</details>


 #### Example usage
  ```bash
  # Windows
  .\poke-cli.exe pokemon charizard --types --abilities
   
  # Unix
  .\poke-cli pokemon vespiquen -t -a
  ```

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
docker run --rm -i -t digitalghostdev/poke-cli:v0.11.0 <command> [subcommand] flag]
```

### Go Install
_If you have Go already, install the executable yourself_

1. Run the following command:
   ```bash
   go install github.com/digitalghost-dev/poke-cli@latest
   ```
2. The tool is ready to use!
---
## Usage
By running `poke-cli [-h | --help]`, it'll display information on how to use the tool. 
```
╭─────────────────────────────────────────────────────────╮
│Welcome! This tool displays data related to Pokémon!     │
│                                                         │
│ USAGE:                                                  │
│    poke-cli [flag]                                      │
│    poke-cli <command> [flag]                            │
│    poke-cli <command> <subcommand> [flag]               │
│                                                         │
│ FLAGS:                                                  │
│    -h, --help      Shows the help menu                  │
│    -l, --latest    Prints the latest version available  │
│    -v, --version   Prints the current version           │
│                                                         │
│ COMMANDS:                                               │
│    natures         Get details about Pokémon natures    │
│    pokemon         Get details about a specific Pokémon │
│    types           Get details about a specific typing  │
╰─────────────────────────────────────────────────────────╯

```

---
## Roadmap
The architecture behind how the tool works is straight forward.
1. Commands indicate which data endpoint to focus on.
2. Flags provide more information and can be all stacked together or chosen.

### Planned for Version 1.0.0

_Not 100% up-to-date, may add or remove some of these choices_

- [ ] `ability`: get data about a specific ability.
    - [ ] `-p | --pokemon`: display Pokémon that learn this ability.
- [ ] `move`: get data about a specific move.
    - [ ] `-p | --pokemon`: display Pokémon that learn this move.
- [x] `natures`: get data about natures.
- [ ] `pokemon`: get data about a specific Pokémon.
   - [x] `-a | --abilities`: display the Pokémon's abilities.
   - [ ] `-i | --image`: display a pixel image of the Pokémon.
   - [x] `-s | --stats`: display the Pokémon's base stats.
   - [x] `-t | --types`: display the Pokémon's typing.
   - [ ] `-m | --moves`: display learnable moves.
- [x] `types`: get data about a specific typing.
