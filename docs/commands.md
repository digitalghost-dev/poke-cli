# Commands

## main

**Available Flags**

* `--help | -h`
* `--latest | -l`
* `--version | -v`

---

## `ability`
* Retrieve information about a specific ability, including its flavor text, 
the generation in which it first appeared, and a list of Pokémon that possess it.

**Available Flags**

* `--help | -h` 
* `--pokemon | -p`

Example:
```console
$ poke-cli ability solar-power
$ poke-cli ability solar-power --pokemon    # list Pokémon that posses the ability
```

Output:

![ability_command](assets/ability.gif)

---

## `move`
* Retrieve information about a specific move, including its type, power, PP, accuracy, category, etc.,
and the move's effect.

Example:
```console
$ poke-cli move dazzling-gleam
```

Output:

![move_command](assets/move.gif)

---

## `natures`
* Retrieve a table of all natures and the stats they affect.

Example:
```console
$ poke-cli natures
```

Output:

![natures_gif](assets/natures.gif)

---

## `pokemon`
* Retrieve information about a specific Pokémon such as available abilities, learnable moves, typing, and base stats. All data is based on generation 9.

**Available Flags**

* `--help | -h`
* `--abilities | -a`
* `--image=xx | -i=xx`
* `--moves | -m`
* `--stats | -s`
* `--types | -t`

Example:
```console
$ poke-cli pokemon rockruff --abilities --moves
```

Output:

![pokemon_abilities_moves](assets/pokemon_abilities_moves.gif)

Example:
```shell
# choose between three sizes: 'sm', 'md', 'lg'
$ poke-cli pokemon tyranitar --image=sm
```

Output:

![pokemon_image](assets/pokemon_image.gif)

Example:
```console
$ poke-cli pokemon cacturne --stats --types
```

Output:

![pokemon_stats_types](assets/pokemon_stats_types.gif)

---

## `search`
* Search for resources from different endpoints. Searchable endpoints include `ability`, `pokemon`, and `move`.

Example:
```console
$ poke-cli search
```

Output:

![search_command](assets/search.gif)

---

## `types`
* Retrieve details about a specific type and a damage relation table.

Example:
```console
$ poke-cli types
```
Output:

![types_command](assets/types.gif)