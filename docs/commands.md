# Commands

## main
* Print the help menu or view information about the program

**Available Flags**

* `--config | -c`
* `--latest | -l`
* `--version | -v`

Exampe:

```bash
# print help menu
poke-cli
# or
poke-cli --help

# edit program settings
poke-cli --config

# check latest release vesion
poke-cli --latest

# check current installed version
poke-cli --version
```

---

## `ability`
* Retrieve information about a specific ability, including its flavor text, 
the generation in which it first appeared, and a list of Pokémon that possess it.

**Available Flags**

* `--pokemon | -p`

Example:
```bash
poke-cli ability solar-power
poke-cli ability solar-power --pokemon    # list Pokémon that posses the ability
```

Output:

![ability_command](assets/command_gifs/ability.gif)

---

## `berry`
* Retrieve information about a specific berry.

Example:
```bash
# specific berry
poke-cli berry oran

# TUI screen
poke-cli berry
```

Output:

![berry_command](assets/command_gifs/berry.gif)

---

## `card`
* Browse Pokémon TCG card data through an interactive TUI.

The command opens a multi-step browser:

1. Select a series.
2. Select a set from that series.
3. Browse the cards in the selected set.
4. Option to open the selected card in the image viewer with `?`.

Card images use your terminal's graphics protocol. Image rendering support depends on the terminal. 

The following terminals are confirmed to have protocol support and render card images correctly:
* Kitty
* WezTerm
* iTerm2
* Ghostty
* Konsole
* Rio
* Tabby
* Windows Terminal

Basic terminal emulators may show card details without images or may not render images correctly.

Example:
```bash
poke-cli card
```

Output:

![card_command](assets/command_gifs/card.gif)

---

## `comp`
* Browse current competitive Pokémon standings through an interactive TUI.

The command opens a competition picker:

1. Select `TCG` or `VGC`.
2. Select a tournament.
3. Browse the tournament dashboard.

The dashboard supports Overview / Standings / Decks / Countries tabs for TCG and Overview / Standings / Usage / Countries tabs for VGC. Press `w` inside the TUI to open the web dashboard.

Example:
```bash
poke-cli comp
```

Output:

![comp_command](assets/command_gifs/comp.gif)

---

## `item`
* Retrieve information about a specific item, including its cost, category and description.

Example:
```bash
poke-cli item poke-ball
```

Output:

![item_command](assets/command_gifs/item.gif)

---

## `mechanics`
* Retrieve data about video game mechanics.

**Available Flags**

* `--natures | -n`

Example:
```bash
poke-cli mechanics --natures
```

Output:

![mechanics_natures_gif](assets/command_gifs/natures.gif)

---

## `move`
* Retrieve information about a specific move, including its type, power, PP, accuracy, category, etc.,
and the move's effect.

Example:
```bash
poke-cli move dazzling-gleam
```

Output:

![move_command](assets/command_gifs/move.gif)

---

## `pokemon`
* Retrieve information about a specific Pokémon such as available abilities, learnable moves, typing, and base stats. All data is based on generation 9.

**Available Flags**

* `-a | --abilities`
* `-d | --defenses`
* `-i=xx | --image=xx`
* `-m | --moves`
* `-s | --stats`

The Pokémon's typing is included in the base `pokemon` command output.

Example:
```bash
poke-cli pokemon rockruff --abilities --moves
```

Output:

![pokemon_abilities_moves](assets/command_gifs/pokemon-abilities-moves.gif)

Example:
```bash
poke-cli pokemon gastrodon --defenses
```

Output:

![pokemon_defense](assets/command_gifs/pokemon-defense.gif)

Example:
```bash
# choose between three sizes: 'sm', 'md', 'lg'
poke-cli pokemon tyranitar --image=sm
```

Output:

![pokemon_image](assets/command_gifs/pokemon-image.gif)

Example:
```bash
poke-cli pokemon cacturne --stats
```

Output:

![pokemon_types](assets/command_gifs/pokemon-stats.gif)

---

## `search`
* Search for resources from different endpoints. Searchable endpoints include `ability`, `pokemon`, and `move`.

Example:
```bash
poke-cli search
```

Output:

![search_command](assets/command_gifs/search.gif)

---

## `speed`
* Calculate the speed of a Pokémon in battle.

The command opens an interactive form and asks for the following values:

* Pokémon name
* Level: `1-100`
* Speed EVs: `0-252`
* Speed IVs: `0-31`
* Modifiers: `Choice Scarf`, `Tailwind`
* Ability: `None`, `Swift Swim`, `Chlorophyll`, `Sand Rush`, `Slush Rush`, `Unburden`, `Quick Feet`, `Surge Surfer`
* Nature multiplier: `+10%`, `0%`, `-10%`
* Speed stage: `-6` to `+6`

The final speed is calculated with the standard stat formula and rounded down.

Example:
```bash
poke-cli speed
```
Output:

![speed_command](assets/command_gifs/speed.gif)

---

## `types`
* Retrieve details about a specific type and a damage relation table.

Example:
```bash
poke-cli types
```
Output:

![types_command](assets/command_gifs/types.gif)
