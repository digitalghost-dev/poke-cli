# Installation Methods
There are a variety of ways to install the tool and support for different operating systems.

* [Homebrew](#homebrew)
* [Winget](#winget)
* [Linux Packages](#linux-packages)
* [Docker Image](#docker-image)
* [Binary](#binary)
* [Source](#source)


### Homebrew
1. Install the Cask:
    ```console
    brew install --cask digitalghost-dev/tap/poke-cli
    ````
2. Verify install:
    ```console
    poke-cli -v
    ```

### Winget
1. Install the package:
    ```powershell
    winget install poke-cli
    ```

2. Verify install:
    ```console
    poke-cli -v
    ```

### Linux Packages
[![Hosted By: Cloudsmith](https://img.shields.io/badge/OSS%20hosting%20by-cloudsmith-blue?logo=cloudsmith&style=flat-square)](https://cloudsmith.com)

This package repository is generously hosted by Cloudsmith.
Cloudsmith is a fully cloud-based service that lets you easily create, store, and share packages in any format, anywhere.

1. Run the **Repository Setup** script first for the correct Linux distribution.
2. Run the corresponding **Installation Command** afterwards.

| Package Type | Distributions                     | Repository Setup                                                                                                                        | Installation Command                         |
|:------------:|-----------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------|
|    `apk`     | Alpine                            | `sudo apk add --no-cache bash && curl -1sLf 'https://dl.cloudsmith.io/basic/digitalghost-dev/poke-cli/setup.alpine.sh' \| sudo -E bash` | `sudo apk add poke-cli=1.6.0 --update-cache` |
|    `deb`     | Ubuntu, Debian                    | `curl -1sLf 'https://dl.cloudsmith.io/public/digitalghost-dev/poke-cli/setup.deb.sh' \| sudo -E bash`                                   | `sudo apt-get install poke-cli=1.6.0`        |
|    `rpm`     | Fedora, CentOS, Red Hat, openSUSE | `curl -1sLf 'https://dl.cloudsmith.io/public/digitalghost-dev/poke-cli/setup.rpm.sh' \| sudo -E bash`                                   | `sudo yum install poke-cli-1.6.0-1`          |

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
    ```console
    docker run --rm -it digitalghostdev/poke-cli:v1.6.0 <command> [subcommand] flag]
    ```
    * Enter the container and use its shell:
    ```console
    docker run --rm -it --name poke-cli --entrypoint /bin/sh digitalghostdev/poke-cli:v1.6.0 -c "cd /app && exec sh"
   # placed into the /app directory, run the program with './poke-cli'
   # example: ./poke-cli ability swift-swim
    ```

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
  ```console
  # Windows
  .\poke-cli.exe pokemon charizard --types --abilities
   
  # Unix
  .\poke-cli ability airlock --pokemon
  ```

### Source

1. Run the following command:
   ```console
   go install github.com/digitalghost-dev/poke-cli@latest
   ```
2. The tool should be ready to use if `$PATH` is set up.
