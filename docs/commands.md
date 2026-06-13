# Commands

## main

**Available Flags**

* `--latest | -l`
* `--version | -v`

---

## `ability`
* Retrieve information about a specific ability, including its flavor text, 
the generation in which it first appeared, and a list of Pokémon that possess it.

**Available Flags**

* `--pokemon | -p`

Example:
```console
poke-cli ability solar-power
poke-cli ability solar-power --pokemon    # list Pokémon that posses the ability
```

Output:

![ability_command](assets/command_gifs/ability.gif)

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
```console
poke-cli card
```

Output:

![card_command](assets/command_gifs/card.gif)

---

## `item`
* Retrieve information about a specific item, including its cost, category and description.

Example:
```console
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
```console
poke-cli mechanics --natures
```

Output:

![mechanics_gif](assets/command_gifs/mechanics.gif)

---

## `move`
* Retrieve information about a specific move, including its type, power, PP, accuracy, category, etc.,
and the move's effect.

Example:
```console
poke-cli move dazzling-gleam
```

Output:

![move_command](assets/command_gifs/move.gif)

---

## `pokemon`
* Retrieve information about a specific Pokémon such as available abilities, learnable moves, typing, and base stats. All data is based on generation 9.

**Available Flags**

* `-a | --abilities`
* `-d | --defense`
* `-i=xx | --image=xx`
* `-m | --moves`
* `-s | --stats`
* `-t | --types`

!!! warning

    The `-t | --types` flag is deprecated will be removed in v2.
    The Pokémon's typing is now included in the base `pokemon` command.

Example:
```console
poke-cli pokemon rockruff --abilities --moves
```

Output:

![pokemon_abilities_moves](assets/command_gifs/pokemon-abilities-moves.gif)

Example:
```console
poke-cli pokemon gastrodon --defense
```

Output:

![pokemon_defense](assets/command_gifs/pokemon-defense.gif)

Example:
```console
# choose between three sizes: 'sm', 'md', 'lg'
poke-cli pokemon tyranitar --image=sm
```

Output:

![pokemon_image](assets/command_gifs/pokemon-image.gif)

Example:
```console
poke-cli pokemon cacturne --stats
```

Output:

![pokemon_types](assets/command_gifs/pokemon-stats.gif)

---

## `search`
* Search for resources from different endpoints. Searchable endpoints include `ability`, `pokemon`, and `move`.

Example:
```console
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
```console
poke-cli speed
```
Output:

![speed_command](assets/command_gifs/speed.gif)

---

## `tcg`
* Retrieve details about all competitive TCG tournaments for the current season.

**Available Flags**

* `--web | -w` - Open the tournament's website in the default browser.

Example:
```console
poke-cli tcg
```

Output:

![tcg_command](assets/command_gifs/tcg.gif)

---

## `types`
* Retrieve details about a specific type and a damage relation table.

Example:
```console
poke-cli types
```
Output:

![types_command](assets/command_gifs/types.gif)
