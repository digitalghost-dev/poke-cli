<div align="center">
    <img width="425" src="poke-cli.png" alt="pokemon-logo"/>
    <h4></h4>
    <img src="https://img.shields.io/github/v/release/digitalghost-dev/poke-cli?style=flat-square&logo=git&logoColor=FFCC00&label=Release%20Version&labelColor=EEE&color=FFCC00" alt="version-label">
    <img src="https://img.shields.io/docker/image-size/digitalghostdev/poke-cli/v2.0.0?arch=arm64&style=flat-square&logo=docker&logoColor=FFCC00&labelColor=EEE&color=FFCC00" alt="docker-image-size">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/ci.yml?branch=main&style=flat-square&logo=github&logoColor=FFCC00&label=CI&labelColor=EEE&color=FFCC00" alt="ci-status-badge">
</div>
<div align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/digitalghost-dev/poke-cli/go_test.yml?style=flat-square&logo=go&logoColor=00ADD8&label=Tests&labelColor=EEE&color=00ADD8" alt="tests-label">
    <img src="https://img.shields.io/github/go-mod/go-version/digitalghost-dev/poke-cli?style=flat-square&logo=Go&labelColor=EEE&color=00ADD8" alt="go-version"/>
    <img src="https://img.shields.io/codecov/c/github/digitalghost-dev/poke-cli?token=05GBSAOQIT&style=flat-square&logo=codecov&logoColor=00ADD8&labelColor=EEE&color=00ADD8" alt="codecov"/>
</div>

## Pokemon CLI
`poke-cli` is a hybrid of a classic CLI and a modern TUI tool for viewing VG and TCG data about Pokémon!

View the [documentation](https://docs.poke-cli.com) on the data infrastructure in [data_platform/](https://github.com/digitalghost-dev/poke-cli/tree/main/data_platform) if you're interested.

* [Demo](#demo)
* [Installation](#installation)
* [Usage](#usage)
* [Roadmap](#roadmap)
* [Tested Terminals](#tested-terminals)

> [!NOTE]
> A version 2 is being planned and built that will remove/update some commands and flags. Refer to the changes under the [Roadmap](#version-2-changes) for more information.

---

## Demo
### Video Game Data

![demo-vg](https://dc8hq8aq7pr04.cloudfront.net/demo-v1.6.0.gif)

### Trading Card Game Data

![demo-tcg](https://dc8hq8aq7pr04.cloudfront.net/poke-cli-card-v1.8.8.gif)

---

## Installation

* [Homebrew](#homebrew)
* [Scoop](#scoop)
* [Linux Packages](#linux-packages)
* [Docker Image](#docker-image)
* [Binary](#binary)
* [Source](#source)


### Homebrew
1. Install the Cask:
    ```bash
    brew install --cask digitalghost-dev/tap/poke-cli
    ````
2. Verify installation:
    ```bash
    poke-cli -v
    ```

### Scoop
1. Add the bucket:
    ```bash
    scoop bucket add digitalghost https://github.com/digitalghost-dev/scoop-bucket.git
    ```

2. Install poke-cli:
    ```bash
    scoop install poke-cli
    ```
   
3. Verify installation:
    ```bash
    poke-cli -v
    ```

### Linux Packages
[![Hosted By: Cloudsmith](https://img.shields.io/badge/OSS%20hosting%20by-cloudsmith-blue?logo=cloudsmith&style=flat-square)](https://cloudsmith.com)

This package repository is generously hosted by Cloudsmith.
Cloudsmith is a fully cloud-based service that lets you easily create, store, and share packages in any format, anywhere.

1. Run the **Repository Setup** script first for the correct Linux distribution.
2. Run the corresponding **Installation Command** afterward.

| Package Type | Distributions                     | Repository Setup                                                                                                                        | Installation Command                   |
|:------------:|-----------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------|
|    `apk`     | Alpine                            | `sudo apk add --no-cache bash && curl -1sLf 'https://dl.cloudsmith.io/basic/digitalghost-dev/poke-cli/setup.alpine.sh' \| sudo -E bash` | `sudo apk add poke-cli --update-cache` |
|    `deb`     | Ubuntu, Debian                    | `curl -1sLf 'https://dl.cloudsmith.io/public/digitalghost-dev/poke-cli/setup.deb.sh' \| sudo -E bash`                                   | `sudo apt-get install poke-cli`        |
|    `rpm`     | Fedora, CentOS, Red Hat, openSUSE | `curl -1sLf 'https://dl.cloudsmith.io/public/digitalghost-dev/poke-cli/setup.rpm.sh' \| sudo -E bash`                                   | `sudo yum install poke-cli`            |

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
    docker run --rm -it digitalghostdev/poke-cli:v2.0.0 <command> [subcommand] [flag]
    ```
   * Enter the container and use its shell:
    ```bash
    docker run --rm -it --name poke-cli --entrypoint /bin/sh digitalghostdev/poke-cli:v2.0.0 -c "cd /app && exec sh"
   # placed into the /app directory, run the program with './poke-cli'
   # example: ./poke-cli ability swift-swim
    ```

> [!NOTE]
> The `card` command renders TCG card images using your terminal's graphics protocol. When running inside Docker, pass your terminal's environment variables so image rendering works correctly:
> ```bash
> # Kitty
> docker run --rm -it -e TERM -e KITTY_WINDOW_ID digitalghostdev/poke-cli:v2.0.0 card
>
> # WezTerm, iTerm2, Ghostty, Konsole, Rio, Tabby
> docker run --rm -it -e TERM -e TERM_PROGRAM digitalghostdev/poke-cli:v2.0.0 card
>
> # Windows Terminal (Sixel)
> docker run --rm -it -e WT_SESSION digitalghostdev/poke-cli:v2.0.0 card
> ```
> If your terminal is not listed above, image rendering is not supported inside Docker.

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

### Source

1. Run the following command:
   ```bash
   go install github.com/digitalghost-dev/poke-cli@latest
   ```
2. The tool should be ready to use if `$PATH` is set up.

> [!TIP]
> `go install` builds only the `poke-cli` binary, **not** the `poke-cache` caching helper (a separate binary that every packaged install bundles). `poke-cli` works the same without it; it just calls PokéAPI directly instead of caching responses on disk. To enable caching, download the `poke-cache` archive for your platform from the [releases](https://github.com/digitalghost-dev/poke-cli/releases/latest) page, extract it, and move the `poke-cache` binary onto your `$PATH`.


---
## Usage
By running `poke-cli [-h | --help]`, it'll display information on how to use the tool or check out the [docs](https://docs.poke-cli.com/)!
```
╭───────────────────────────────────────────────────────────────╮
│Welcome! This tool displays data related to Pokémon!           │
│                                                               │
│ USAGE:                                                        │
│    poke-cli [flag]                                            │
│    poke-cli <command> [flag]                                  │
│    poke-cli <command> <subcommand> [flag]                     │
│                                                               │
│ FLAGS:                                                        │
│    -h, --help      Shows the help menu                        │
│    -l, --latest    Prints the latest version available        │
│    -v, --version   Prints the current version                 │
│                                                               │
│ COMMANDS:                                                     │
│    ability         Get details about an ability               │
│    berry           Get details about a berry                  │
│    card            Get details about a TCG card               │
│    item            Get details about an item                  │
│    mechanics       Get details about video game mechanics     │
│    move            Get details about a move                   │
│    pokemon         Get details about a Pokémon                │
│    search          Search for a resource                      │
│    speed           Calculate the speed of a Pokémon in battle │
│    tcg             Get details about TCG tournaments          │
│    types           Get details about a typing                 │
│                                                               │
│ hint: when calling a resource with a space, use a hyphen      │
│ example: poke-cli ability strong-jaw                          │
│ example: poke-cli pokemon flutter-mane                        │
│                                                               │
│ ↓ ctrl/cmd + click for docs/guides                            │
│ docs.poke-cli.com                                             │
╰───────────────────────────────────────────────────────────────╯
```

---

## Roadmap
Below is a list of the planned/completed commands and flags:

- [x] `ability`: get data about an ability.
    - [x] `-p | --pokemon`: display Pokémon that learn this ability.
- [x] `berry`: get data about a berry.
- [ ] `card`: get data about a TCG card.
    - [x] add mega evolution data
    - [x] add scarlet & violet data
    - [x] add sword & shield data
    - [x] add sun & moon data
    - [ ] add x & y data
- [x] `item`: get data about an item.
- [x] `mechanics`: get data about game mechanics.
    - [x] `-n | --natures`: display a table of all natures.
- [ ] `move`: get data about a move.
    - [ ] `-p | --pokemon`: display Pokémon that learn this move.
- [ ] `pokemon`: get data about a Pokémon.
    - [x] `-a | --abilities`: display the Pokémon's abilities.
    - [ ] `-c | --cry`: play the Pokémon's cry.
    - [x] `-d | --defenses`: display the Pokémon's type defences.
    - [x] `-i | --image`: display a pixel image of the Pokémon.
    - [x] `-m | --moves`: display learnable moves.
    - [x] `-s | --stats`: display the Pokémon's base stats.
- [ ] `search`: search for a resource 
    - [x] `ability`
    - [ ] `berry`
    - [ ] `item`
    - [x] `move`
    - [x] `pokemon`
- [x] `speed`: compare speed stats between two Pokémon.
- [x] `tcg`: get data about TCG tournaments.
- [x] `types`: get data about a specific typing.

### Version 2 Changes
The following planned changes in `v2`:

- `pokemon <name> -t | --types` — removed; typing is included by default.
- `pokemon <name> --defense` - being renamed to `--defenses` to keep consistency with other flags in the `pokemon` command.
- `natures` — moves to a flag under a new `mechanics` command.
- `tcg` — moves to a new `comp` command (covers competitive TCG *and* VGC data).
- Adding `pflag` library to enforce POSIX style flags.


---
## Tested Terminals
| Terminal           | OS                            | Status | Issues                                                                            |
|--------------------|-------------------------------|:------:|-----------------------------------------------------------------------------------|
| Alacritty          | macOS, Ubuntu, Windows        |   🟡   | No support for TCG images                                                         |
| Foot               | Ubuntu, Fedora                |   🟢   | None                                                                              |
| Ghostty            | macOS                         |   🟢   | None                                                                              |
| iTerm2             | macOS                         |   🟢   | None                                                                              |
| Kitty              | macOS, Ubuntu, Debian, Fedora |   🟢   | None                                                                              |
| Rio                | macOS                         |   🟢   | None                                                                              |
| Tabby              | Ubuntu                        |   🟢   | None                                                                              |
| Terminal (Alpine)  | Alpine                        |   🟡   | Some colors aren't supported<br>`pokemon <name> --image=xx` flag has pixel issues |
| Terminal (Linux)   | Ubuntu, Debian, Fedora        |   🟡   | No support for TCG images                                                         |
| Terminal (macOS)   | macOS                         |   🟠   | No support for TCG images<br>`pokemon <name> --image=xx` flag has pixel issues    |
| Terminal (Windows) | Windows                       |   🟢   | None                                                                              |
| WezTerm            | macOS, Windows                |   🟡   | Windows version has issues with displaying TCG images                             |