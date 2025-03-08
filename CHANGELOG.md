# Change Log
This change log provides version history for tags prior to `v1.0.0`

Due to the lack of control over dates when creating a release, I've decided to move the change history per tag prior to `v1.0.0` to a `CHANGELOG.md` file.

Some previous releases had their dates screwed up and without a way to backdate a release like [GitLab](https://docs.gitlab.com/api/releases/#update-a-release), the release history is out of date.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

* **MAJOR:** Any changes to the backend infrastructure that requires new methods of moving data that won't work with the previous architecture, mainly with the addition of new databases or data sources.
* **MINOR:** Any changes to the Streamlit dashboard that adds a new interaction/feature or removal of one.
* **PATCH:** Any changes that fix bugs, typos or small edits.

# Version History

## v0.12.2
_March 7th, 2024_

### Changed
* Changed location of style definitions and structs to their own respective packages for scalability. [(#106)](https://github.com/digitalghost-dev/poke-cli/issues/106)

### Details
* **Commit Change Log:** [v0.12.1 > v0.12.2](https://github.com/digitalghost-dev/poke-cli/compare/v0.12.1...v0.12.2)
* **Associated PRs:** [#107](https://github.com/digitalghost-dev/poke-cli/pull/107)

---

## v0.12.1
_March 4th, 2025_

### Changed
* Changed help menus for better clarity. [(#103)](https://github.com/digitalghost-dev/poke-cli/issues/103)

### Fixed
* Fixed error handling issue with `ability` command. [(#104)](https://github.com/digitalghost-dev/poke-cli/issues/104)

### Details
* **Commit Change Log:** [v0.12.0 > v0.12.1](https://github.com/digitalghost-dev/poke-cli/compare/v0.12.0...v0.12.1)
* **Associated PRs:** [#105](https://github.com/digitalghost-dev/poke-cli/pull/105)

---

## v0.12.0
_February 2nd, 2025_

### Added
* Added new `ability` command. Get details about a specific ability. [(#101)](https://github.com/digitalghost-dev/poke-cli/issues/101)

    ```bash
    poke-cli ability unware

    # output
    Unaware
    Effect: Ignores other Pokémon's stat modifiers for damage and accuracy calculation.
    ```

* Added `-p | --pokemon` flag to the `ability` command. [(#99)](https://github.com/digitalghost-dev/poke-cli/issues/99)

    ```bash
    poke-cli ability stench --pokemon

    Stench
    Effect: Has a 10% chance of making target Pokémon flinch with each hit.

    Pokemon with Stench

    1. Gloom                       2. Grimer                      3. Muk
    4. Koffing                     5. Weezing                     6. Stunky
    7. Skuntank                    8. Trubbish                    9. Garbodor
    10. Garbodor-Gmax
    ```

* Added better error messaging to API call functions. [(#98)](https://github.com/digitalghost-dev/poke-cli/issues/98)
* Added a designated file for holding style variables to the `flags/` directory. [(#100)](https://github.com/digitalghost-dev/poke-cli/issues/100)

### Details
* **Commit Change Log:** [v0.11.1 > v0.12.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.11.1...v0.12.0)
* **Associated PRs:** [#102](https://github.com/digitalghost-dev/poke-cli/pull/102)

---

## v0.11.1
_January 13th, 2024_

### Fixed
* Fix issue of using the `--image | -i` flag without an argument but with a `=` at the end. [(#96)](https://github.com/digitalghost-dev/poke-cli/issues/96)
    * For example, running: `poke-cli pokemon cacturne --image=` would not error out the program or print an image. Running this in `v0.11.1` returns an error.

### Details
* **Commit Change Log:** [v0.11.0 > v0.11.1](https://github.com/digitalghost-dev/poke-cli/compare/v0.11.0...v0.11.1)
* **Associated PRs:** [#97](https://github.com/digitalghost-dev/poke-cli/pull/97)

---

## v0.11.0
_January 6th, 2024_

### Added
* Added a new flag for generating Pokémon sprites with options for different sizes. [(#93)](https://github.com/digitalghost-dev/poke-cli/issues/93)
    * _Example:_ `poke-cli pokemon gengar -i=sm`

### Details
* **Commit Change Log:** [v0.10.0 > v0.11.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.10.0...v0.11.0)
* **Associated PRs:** [#94](https://github.com/digitalghost-dev/poke-cli/pull/94)

**Contributions**
* @ancientcatz - replaced simpleicons image with local `.svg` file.

---

## v0.10.0
_December 26th, 2024_

### Added
* Added total base points to output of the `-s | --stats` flag. [(#87)](https://github.com/digitalghost-dev/poke-cli/issues/87)
* Added a navigation menu while selecting a type. [(#90)](https://github.com/digitalghost-dev/poke-cli/issues/90)
* Added a `natures` command that prints out a table of all natures and affected stats. [(#91)](https://github.com/digitalghost-dev/poke-cli/issues/91)

### Details
* **Commit Change Log:** [v0.9.3 > v0.10.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.9.3...v0.10.0)
* **Associated PRs:** [#92](https://github.com/digitalghost-dev/poke-cli/pull/92)

---

## v0.9.3
_December 16th, 2024_

### Fixed
* Fixed spacing in damage table when selecting type with the `types` command. This will help the chart fit on smaller terminal windows. [(#88)](https://github.com/digitalghost-dev/poke-cli/issues/88)

### Details
* **Commit Change Log:** [v0.9.2 > v0.9.3](https://github.com/digitalghost-dev/poke-cli/compare/v0.9.2...v0.9.3)
* **Associated PRs:** [#89]((https://github.com/digitalghost-dev/poke-cli/pull/89)


---

## v0.9.2
_December 14th, 2024_

### Fixed
* Fixed various string formatting, specifically with Pokémon names and abilities using a hyphen. [(#85)]((https://github.com/digitalghost-dev/poke-cli/issues/85))

### Details
* **Commit Change Log:** [v0.9.1 > v0.9.2](https://github.com/digitalghost-dev/poke-cli/compare/v0.9.1...v0.9.2)
* **Associated PRs:** [#86](https://github.com/digitalghost-dev/poke-cli/pull/86)

---

## v0.9.1
_December 8th, 2024_

### Fixed
* Fixed error message outputs to use a `lipgloss` border, specifically when using an empty flag (`poke-cli pokemon -s --`) and when mistyping a Pokémon's name. [(#83)](https://github.com/digitalghost-dev/poke-cli/issues/83)

### Details
* **Commit Change Log:** [v0.9.0 > v0.9.1](https://github.com/digitalghost-dev/poke-cli/compare/v0.9.0...v0.9.1)
* **Associated PRs:** [#84](https://github.com/digitalghost-dev/poke-cli/pull/84)

---

## v0.9.0
_December 6th, 2024_

### Added
* Added a new flag `-v | --version` to check the current version of the tool. [(#78)](https://github.com/digitalghost-dev/poke-cli/issues/78)
* Added the number of moves belonging to a specified type when using the `types` command. [(#79)](https://github.com/digitalghost-dev/poke-cli/issues/79)

### Security
* Secured against G107 (CWE-88). [(#80)](https://github.com/digitalghost-dev/poke-cli/issues/80)

### Details
* **Commit Change Log:** [v0.8.0 > v0.9.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.8.0...v0.9.0)
* **Associated PRs:** [#81](https://github.com/digitalghost-dev/poke-cli/pull/81)

---

## v0.8.0
_November 27th, 2024_

### Added
* Added a new flag `-s | --stats` to the `pokemon` command to view a Pokémon's base stats. [(#74)](https://github.com/digitalghost-dev/poke-cli/issues/74)
* Added a function to print out the header for any flag option under the `pokemon` command, reducing redundancy. [(#75)](https://github.com/digitalghost-dev/poke-cli/issues/75)
* Added metrics to the `pokemon` command that includes height and weight. [(#76)](https://github.com/digitalghost-dev/poke-cli/issues/76)

### Details
* **Commit Change Log:** [v0.7.2 > v0.8.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.7.2...v0.8.0)
* **Associated PRs:** [#77](https://github.com/digitalghost-dev/poke-cli/pull/77)

---

## v0.7.2
_November 21st, 2024_

### Changed
* Changed location of `getTypeColor()` function to `cmd/styles.go`. [(#71)]((https://github.com/digitalghost-dev/poke-cli/issues/71)

### Removed
* Removed `httpGet` variable from `connections/connection.go`. [(#72)](https://github.com/digitalghost-dev/poke-cli/issues/72)


### Details
* **Commit Change Log:** [v0.7.1 > v0.7.2](https://github.com/digitalghost-dev/poke-cli/compare/v0.7.1...v0.7.2)
* **Associated PRs:** [#73](https://github.com/digitalghost-dev/poke-cli/pull/73)

---

## v0.7.1
_November 12th, 2024_

### Added
* Added a helper method that helps reduce redundancy by passing it into each command's argument validator. [(#69)](https://github.com/digitalghost-dev/poke-cli/issues/69)

### Changed
* Changed the help menus for the current commands by simplifying the text and verbiage. [(#67)](https://github.com/digitalghost-dev/poke-cli/issues/67)
* Changed the `pokemon` and `types` help menu to include the `-h, --help` flag. [(#65)](https://github.com/digitalghost-dev/poke-cli/issues/65)

### Fixed
* Fixed a bug where misspelling a Pokémon's name would error out but still would run the rest of the program with blank data. [(#68)](https://github.com/digitalghost-dev/poke-cli/issues/68)
* Fixed a bug when using the help flag on either the `pokemon` or `types` command; the program defaults to using the `mainFlagSet.Usage()` menu instead of the help menu for each command. [(#64)](https://github.com/digitalghost-dev/poke-cli/issues/64)

### Details
* **Commit Change Log:** [v0.7.0 > v0.7.1](https://github.com/digitalghost-dev/poke-cli/compare/v0.7.0...v0.7.1)
* **Associated PRs:** [#70](https://github.com/digitalghost-dev/poke-cli/pull/70)

---

## v0.7.0
_November 2nd, 2024_

### Added
* Added a damage chart to `poke-cli types`. It shows a list of the selected type's weaknesses, resistances, immunities, and damage to other types using a [BubbleTea](https://github.com/charmbracelet/bubbletea/tree/main) list. [(#62)](https://github.com/digitalghost-dev/poke-cli/issues/62)

### Changed
* Changed the `ApiCallSetup()` function in `connections/connection.go` to return `fmt.Errorf()` instead of various `fmt.Println()` and `log.Fatalf()` uses which made testing the file difficult. Now, testing coverage for the `connections` package has increased to over 85%. [(#66)](https://github.com/digitalghost-dev/poke-cli/issues/66)

### Details
* **Commit Change Log:** [v0.6.5 > v0.7.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.6.5...v0.7.0)
* **Associated PRs:** [#63](https://github.com/digitalghost-dev/poke-cli/pull/63)

---

## v0.6.5
_October 29th, 2024_

### Changed
* Changed `cli.go` to provide test coverage opportunities and better maintainability. Moved logic out of `main()` and into a separate `runCLI` function. [(#58)]((https://github.com/digitalghost-dev/poke-cli/issues/58))

### Fixed
* Fixed expression always equating to `false`. [(#57)](https://github.com/digitalghost-dev/poke-cli/issues/57)

### Details
* **Commit Change Log:** [v0.6.4 > v0.6.5](https://github.com/digitalghost-dev/poke-cli/compare/v0.6.4...v0.6.5)
* **Associated PRs:** [#59](https://github.com/digitalghost-dev/poke-cli/pull/59)

---

## v0.6.4
_October 26th, 2024_

### Changed
* Changed output to include color when printing out a type for the `poke-cli types` command. [(#55)](https://github.com/digitalghost-dev/poke-cli/issues/55)

### Details
* **Commit Change Log:** [v0.6.3 > v0.6.4](https://github.com/digitalghost-dev/poke-cli/compare/v0.6.3...v0.6.4)
* **Associated PRs:** [#56](https://github.com/digitalghost-dev/poke-cli/pull/55)

---

## v0.6.3
_October 20th, 2024_

### Changed
* Changed code in `types.go` for better organization. [(#51)](https://github.com/digitalghost-dev/poke-cli/issues/51)
* Changed message when using a non-available command. [(#52)](https://github.com/digitalghost-dev/poke-cli/issues/52)
    * Previously, only a generic `Unknown command` error would show when not using an available commanded.
* Changed location of `lipgloss` variables to their own file since they are used multiple times throughout the `cmd` package. [(#53)](https://github.com/digitalghost-dev/poke-cli/issues/53)

### Details
* **Commit Change Log:** [v0.6.2 > v0.6.3](https://github.com/digitalghost-dev/poke-cli/compare/v0.6.2...v0.6.3)
* **Associated PRs:** [#54](https://github.com/digitalghost-dev/poke-cli/pull/54)

---

## v0.6.2
_October 7th, 2024_

### Changed
* Changed the location of checking for an `-h` or `--help` flag from `cmd/types.go` to its intended place in `cmd/validateargs.go` where the rest of the validation takes place. (#48)

### Fixed
* Fixed not directly handling when the 3rd argument is `-h` or `--help` after using `poke-cli types`. [(#49)](https://github.com/digitalghost-dev/poke-cli/issues/49)
* Fixed `os.Exit(0)` being called unconditionally, which interrupted tests. Now, the program can exit and not affect tests. [(#49)](https://github.com/digitalghost-dev/poke-cli/issues/49)

### Details
* **Commit Change Log:** [v0.6.1 > v0.6.2](https://github.com/digitalghost-dev/poke-cli/compare/v0.6.1...v0.6.2)
* **Associated PRs:** [#50](https://github.com/digitalghost-dev/poke-cli/pull/50)

---

## v0.6.1
_October 5th, 2024_

### Changed
*  Changed code for better organization by moving the results of table selection when using `poke-cli types` to a separate helper method. [(#42)](https://github.com/digitalghost-dev/poke-cli/issues/42)
* Changed the usage of `fmt.Errorf()` by inputting a styled error message to improve readability and easier future changes. [(#45)](https://github.com/digitalghost-dev/poke-cli/issues/45)
    * For example:
    ```go
    // new
    if len(args) < 3 {
        errMessage := errorBorder.Render(errorColor.Render("Error!"), "some text...")
        return fmt.Errorf("%s", errMessage)
    }

   // old
    if len(args) < 3 {
        return fmt.Errorf(errorBorder.Render(errorColor.Render("Error!"), "some text..."))
    }
    ```

### Fixed
* Fixed an issue where the `-h` or `--help` flag for `poke-cli types` would incorrectly return an error (`fmt.Errorf`) instead of exiting gracefully. Now, using the `-h` or `--help` flags correctly stops the program and displays the help menu. [(#44)](https://github.com/digitalghost-dev/poke-cli/issues/44)

### Details
* **Commit Change Log:** [v0.6.0 > v0.6.1](https://github.com/digitalghost-dev/poke-cli/compare/v0.6.0...v0.6.1)
* **Associated PRs:** [#47](https://github.com/digitalghost-dev/poke-cli/pull/47)

---

## v0.6.0
_September 30th, 2024_

### Added
* Added a new `types` command. (#41)
    * This command can be used to get details about a specific typing.
    * For example:
    ```
    poke-cli types
    ```
* Added a new `cmd/validateargs.go` file that will validate each commands given arguments. [(#39)](https://github.com/digitalghost-dev/poke-cli/issues/39)

### Changed
*  Changed the program's help menu to include the new `types` command. [(#37)](https://github.com/digitalghost-dev/poke-cli/issues/37)

### Fixed
* Fixed the argument check when running only the program name to not display the program's help menu. [(#38)](https://github.com/digitalghost-dev/poke-cli/issues/38)

### Details
* **Commit Change Log:** [v0.5.2 > v0.6.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.5.2...v0.6.0)
* **Associated PRs:** [#40](https://github.com/digitalghost-dev/poke-cli/pull/40)

---

## v0.5.2
_September 21st, 2024_

### Changed
*  Changed the main help menu to use a `lipgloss` border. [(#34)](https://github.com/digitalghost-dev/poke-cli/issues/34)

### Fixed
* Fixed spacing/tabbing issues with `poke-cli pokemon -h` help menu. [(#35)](https://github.com/digitalghost-dev/poke-cli/issues/35)

### Details
* **Commit Change Log:** [v0.5.1 > v0.5.2](https://github.com/digitalghost-dev/poke-cli/compare/v0.5.1...v0.5.2)
* **Associated PRs:** [#36](https://github.com/digitalghost-dev/poke-cli/pull/36)

___

## v0.5.1
_September 9th, 2024_

### Changed
* Changed code for better organization and output. [(#31)](https://github.com/digitalghost-dev/poke-cli/issues/31)
* Changed help menu output when using `-h` or `--help` after declaring a Pokémon's name. [(#32)](https://github.com/digitalghost-dev/poke-cli/issues/32)

Example:
`poke-cli pokemon lycanroc --help`

Output:
```
 poke-cli pokemon <pokemon-name> [flags]                
                                                                                            
 FLAGS:                                                                              
     -a, --abilities      Prints out the Pokémon's abilities.
     -t, --types          Prints out the Pokémon's typing.    
```

### Details
* **Commit Change Log:** [v0.5.0 > v0.5.1](https://github.com/digitalghost-dev/poke-cli/compare/v0.5.0...v0.5.1)
* **Associated PRs:** [#33](https://github.com/digitalghost-dev/poke-cli/pull/33)

---

## v0.5.0
_August 27th, 2024_

### Changed
* Changed the way the tool is used by requiring the name of the API endpoint in `os.Args`. This change will allow the tool to use different endpoint in future updates. [(#26)](https://github.com/digitalghost-dev/poke-cli/issues/26)
* Changed the `subcommands/` directory to `cmd/` to better fit the tool's architecture. [(#29)](https://github.com/digitalghost-dev/poke-cli/issues/29)

**Example:**
```bash
$ poke-cli pokemon cacturne -t -a

# future example
$ poke-cli types
```

### Details
* **Commit Change Log:** [v0.4.0 > v0.5.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.4.0...v0.5.0)
* **Associated PRs:** [#30](https://github.com/digitalghost-dev/poke-cli/pull/30)

---

## v0.4.0
_August 1st, 2024_

### Added
* Added a flag to check the latest Docker image tag on DockerHub and the latest release tag on GitHub. [(#24)](https://github.com/digitalghost-dev/poke-cli/issues/24)
    *  `-l` or `--latest`

**Example:**
```bash
[21:15:16] ~ $ poke-cli -l
Latest Docker image version: v0.4.0
Latest release tag: v0.4.0
```

### Details
* **Commit Change Log:** [v0.3.2 > v0.4.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.3.2...v0.4.0)
* **Associated PRs:** [#25](https://github.com/digitalghost-dev/poke-cli/pull/25)

---

## v0.3.2
_July 18th, 2024_

### Added
* Added short aliases for current flags and will include in future flags. [(#20)](https://github.com/digitalghost-dev/poke-cli/issues/20)
    *  `--abilities` = `-a`
    * `--types` = `-t`

### Changed
* Changed the way a Pokémon's hidden ability in the output. [(#21)](https://github.com/digitalghost-dev/poke-cli/issues/21)

### Details
* **Commit Change Log:** [v0.3.1 > v0.3.2](https://github.com/digitalghost-dev/poke-cli/compare/v0.3.1...v0.3.2)
* **Associated PRs:** [#22](https://github.com/digitalghost-dev/poke-cli/pull/22)

---

## v0.3.1
_July 12th, 2024_

### Changed
* Changed logic behind validating arguments given to the command. [(#16)](https://github.com/digitalghost-dev/poke-cli/issues/16)

### Fixed
* Fixed tests to verify output with new `--abilities` flag added in `v0.3.0`. [(#15)](https://github.com/digitalghost-dev/poke-cli/issues/15)
* Fixed `--help` message to display to information about the `--abilities` flag. [(#17)](https://github.com/digitalghost-dev/poke-cli/issues/17)

### Details
* **Commit Change Log:** [v0.3.0 > v0.3.1](https://github.com/digitalghost-dev/poke-cli/compare/v0.3.0...v0.3.1)
* **Associated PRs:** [#18](https://github.com/digitalghost-dev/poke-cli/pull/18)

---

## v0.3.0
_July 7th, 2024_

### Added
* Added a new flag `--abilities` that prints out the Pokémon's ability or abilities. [(#13)](https://github.com/digitalghost-dev/poke-cli/issues/13)

Example:
```
poke-cli ambipom --abilities
```

output:

![image](https://github.com/digitalghost-dev/poke-cli/assets/86637723/71a79d1f-9b6c-4245-873f-b04cebb52e0d)

### Details
* **Commit Change Log:** [v0.2.0 > v0.3.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.2.0...v0.3.0)
* **Associated PRs:** [#14](https://github.com/digitalghost-dev/poke-cli/pull/14)

---

## v0.2.0
_June 8th, 2024_

### Added
- Added the Pokémon's national dex number to output. [(#8)](https://github.com/digitalghost-dev/poke-cli/issues/8)
- Added a `colorMap` with all colors based on the 18 types. When declaring the `--types` flag after a Pokémon's name, the printed types will be presented in their respective color. Colors are based on the type chart from [pokemondb.net](https://pokemondb.net/type). [(#9)](https://github.com/digitalghost-dev/poke-cli/issues/9)
- Added `strings.ToLower()` to allow mixed-case typing of a Pokémon's name. [(#10)](https://github.com/digitalghost-dev/poke-cli/issues/10)
- Added a header when declaring the `--types` flag for a better organized output. [(#11)](https://github.com/digitalghost-dev/poke-cli/issues/11)

### Details
* **Commit Change Log:** [v0.1.1 > v0.2.0](https://github.com/digitalghost-dev/poke-cli/compare/v0.1.1...v0.2.0)
* **Associated PRs:** [#12](https://github.com/digitalghost-dev/poke-cli/pull/12)

---

## v0.1.1
_June 4th, 2024_

### Changed
- Moved the logic for the `--types` flag under the `TypesFlag()` function instead of it being handing under the `connections/connection.go` file. [(#5)](https://github.com/digitalghost-dev/poke-cli/issues/5)
- Moved the `type Pokemon struct` outside of functions and created a single `struct` instead of breaking it up per function. [(#6)](https://github.com/digitalghost-dev/poke-cli/issues/6)

### Details
* **Commit Change Log:** [v0.1.0. > v0.1.1](https://github.com/digitalghost-dev/poke-cli/compare/v0.1.0...v0.1.1)
* **Associated PRs:** [#7](https://github.com/digitalghost-dev/poke-cli/pull/7)

---

## v0.1.0
_June 1st, 2024_

### Added
- Added a connections file that'll hold all the logic for calling the API. [(#1)](https://github.com/digitalghost-dev/poke-cli/issues/1)
- Added a flagset to be used with `pokemon` subcommand. [(#2)](https://github.com/digitalghost-dev/poke-cli/issues/2)
- Added a subcommand that takes `os.Args[1]` for use after calling the tool. [(#3)](https://github.com/digitalghost-dev/poke-cli/issues/3)
    - For example:
```bash
~ $ poke-cli bulbasaur --types

# output:
> Selected Pokémon: Bulbasaur
> Type 1: grass
> Type 2: poison
```

### Details
* [Full Changelog](https://github.com/digitalghost-dev/poke-cli/commits/v0.1.0)
* **Associated PRs:** [#4](https://github.com/digitalghost-dev/poke-cli/pull/4)