# Installation Methods
There are a variety of ways to install the tool and support for different operating systems.

### Binary

1. Head to the [releases](https://github.com/digitalghost-dev/poke-cli/releases) page of the project.
2. Choose a version to download. The latest is best.
3. Choose an operating system and click on the matching zipped folder to start the download.
4. Extract the folder. The tool is ready to use.
5. Either change directories into the extracted folder or move the binary to a chosen directory.
6. Run the tool!

??? info "View Image of Settings Screen"
    For macOS, you may have to allow the executable to run as it is not signed. 
    Head to System Settings > Privacy & Security > scroll down and allow executable to run.

    ![settings](https://poke-cli-s3-bucket.s3.us-west-2.amazonaws.com/macos_privacy_settings.png)


Example usage
  ```bash
  # Windows
  .\poke-cli.exe pokemon charizard --types --abilities
   
  # Unix
  .\poke-cli ability airlock --pokemon
  ```

---

### Linux Packages
[![Hosted By: Cloudsmith](https://img.shields.io/badge/OSS%20hosting%20by-cloudsmith-blue?logo=cloudsmith&style=flat-square)](https://cloudsmith.com)

_Coming Soon..._

---

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
    docker run --rm -it digitalghostdev/poke-cli:v1.5.2 <command> [subcommand] flag]
    ```
    * Enter the container and use its shell:
    ```bash
    docker run --rm -it --name poke-cli --entrypoint /bin/sh digitalghostdev/poke-cli:v1.5.2 -c "cd /app && exec sh"
   # placed into the /app directory, run the program with './poke-cli'
   # example: ./poke-cli ability swift-swim
    ```
   
---

### Homebrew
1. Install the Cask:
    ```bash
    brew install --cask digitalghost-dev/tap/poke-cli
    ````
2. Verify install:
    ```bash
    poke-cli -v
    ```
   
---

### Winget
1. Install the package:
    ```powershell
    winget install poke-cli
    ```

2. Verify install:
    ```bash
    poke-cli -v
    ```
   
---

### Source

1. Run the following command:
   ```bash
   go install github.com/digitalghost-dev/poke-cli@latest
   ```
2. The tool is ready to use!
