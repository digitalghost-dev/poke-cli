# Refactoring TODO List for v1.8.x

> Last reviewed: December 17, 2025

This document tracks DRY improvements, code quality fixes, and refactoring opportunities for the poke-cli project during the v1.8.x maintenance phase.

---

## Completed Items

These items from the previous review have been addressed:

- [x] **Error surfacing in image viewer** (`cmd/card/imageviewer.go`)
  - `imageReadyMsg` now carries a proper `err error` field
  - `Update()` properly checks `msg.err` and sets `m.Error`
  - `View()` uses `styling.Red.Render()` for error display

- [x] **Avoid printing from library functions** (`connections/connection.go`)
  - No more `fmt.Println(err)` in connection functions
  - All errors are properly returned to callers
  - `FetchEndpoint` returns styled error messages

---

## High Priority

### 1. Consolidate argument validators (`cmd/utils/validateargs.go`)

**Current state:** 9 separate `ValidateXArgs` functions with significant duplication.

**Issue:**
- `ValidateBerryArgs`, `ValidateCardArgs`, `ValidateNaturesArgs`, `ValidateSearchArgs`, `ValidateSpeedArgs`, `ValidateTypesArgs` are nearly identical (just different command names)
- All call `checkLength()` and `checkNoOtherOptions()` with the same pattern

**Suggestion:** Create a parameterized validator:

```go
type ValidatorConfig struct {
    MaxArgs       int
    CommandName   string  // for error message
    RequireName   bool    // some commands require a resource name (ability, item, move)
}

func ValidateArgs(args []string, cfg ValidatorConfig) error {
    if err := checkLength(args, cfg.MaxArgs); err != nil {
        return err
    }
    if cfg.RequireName && len(args) == 2 {
        return fmt.Errorf("Please specify a %s", cfg.CommandName)
    }
    if err := checkNoOtherOptions(args, cfg.MaxArgs, cfg.CommandName); err != nil {
        return err
    }
    return nil
}
```

**Impact:** Reduce from 9 functions to ~3 (generic + ability + pokemon special cases)

---

### 2. Deduplicate command list in cli.go

**Current state:** Command names and descriptions appear twice:
- Lines 73-82: Help menu under "COMMANDS:"
- Lines 151-160: Error message when invalid command provided

**Suggestion:**

```go
var commandDescriptions = []struct {
    name string
    desc string
}{
    {"ability", "Get details about an ability"},
    {"berry", "Get details about a berry"},
    // ...
}

func renderCommandList() string {
    var sb strings.Builder
    for _, cmd := range commandDescriptions {
        sb.WriteString(fmt.Sprintf("\n\t%-15s %s", cmd.name, cmd.desc))
    }
    return sb.String()
}
```

**Impact:** Single source of truth for command list; easier to add/remove commands

---

### 3. Return struct from SetupPokemonFlagSet (`flags/pokemonflagset.go:43-82`)

**Current state:** Returns 13 separate values:
```go
return pokeFlags, abilitiesFlag, shortAbilitiesFlag, defenseFlag, shortDefenseFlag,
       imageFlag, shortImageFlag, moveFlag, shortMoveFlag, statsFlag, shortStatsFlag,
       typesFlag, shortTypesFlag
```

**Suggestion:**

```go
type PokemonFlags struct {
    FlagSet        *flag.FlagSet
    Abilities      *bool
    ShortAbilities *bool
    Defense        *bool
    ShortDefense   *bool
    Image          *string
    ShortImage     *string
    Move           *bool
    ShortMove      *bool
    Stats          *bool
    ShortStats     *bool
    Types          *bool
    ShortTypes     *bool
}

func SetupPokemonFlagSet() *PokemonFlags {
    pf := &PokemonFlags{}
    pf.FlagSet = flag.NewFlagSet("pokeFlags", flag.ExitOnError)
    pf.Abilities = pf.FlagSet.Bool("abilities", false, "...")
    // ...
    return pf
}

// Helper method
func (pf *PokemonFlags) IsAbilitiesSet() bool {
    return *pf.Abilities || *pf.ShortAbilities
}
```

**Impact:** Improves readability at call site in `pokemon.go:48`

---

## Medium Priority

### 4. Remove manual help flag checking

**Current state:** Same pattern repeated in 7 command files:
```go
if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
    flag.Usage()
    return output.String(), nil
}
```

**Files affected:**
- `cmd/ability/ability.go:40-43`
- `cmd/berry/berry.go:36-39`
- `cmd/item/item.go:38-41`
- `cmd/move/move.go:33-36`
- `cmd/pokemon/pokemon.go:54-57`
- `cmd/natures/natures.go` (check)
- `cmd/speed/speed.go` (check)

**Suggestion:** Create helper in `cmd/utils/`:

```go
func CheckHelpFlag(output *strings.Builder, usageFunc func()) bool {
    if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
        usageFunc()
        return true
    }
    return false
}
```

Then in commands:
```go
if utils.CheckHelpFlag(&output, flag.Usage) {
    return output.String(), nil
}
```

---

### 5. Centralize header styling (`flags/pokemonflagset.go:28-41`)

**Current state:** `header()` function defined locally:
```go
func header(header string) string {
    HeaderBold := lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("#FFCC00")).
        BorderTop(true).
        Bold(true).
        Render(header)
    // ...
}
```

**Suggestion:** Move to `styling/styling.go`:

```go
var HeaderStyle = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("#FFCC00")).
    BorderTop(true).
    Bold(true)

func RenderHeader(text string) string {
    return HeaderStyle.Render(text)
}
```

---

### 6. Add text formatting utilities to styling package

**Current state:** Title-casing with hyphen replacement repeated throughout:
```go
cases.Title(language.English).String(strings.ReplaceAll(name, "-", " "))
```

**Locations:**
- `cmd/ability/ability.go:84`
- `cmd/move/move.go:60`
- `cmd/item/item.go:63`
- `cmd/pokemon/pokemon.go:86, 116, 140, 176`
- `flags/pokemonflagset.go:104, 215, 253, 452-453, 664`

**Suggestion:** Add to `styling/styling.go`:

```go
var titleCaser = cases.Title(language.English)

// CapitalizeResourceName converts "strong-jaw" to "Strong Jaw"
func CapitalizeResourceName(name string) string {
    return titleCaser.String(strings.ReplaceAll(name, "-", " "))
}
```

---

### 7. Use HandleFlagError consistently

**Current state:** `cmd/ability/ability.go:54-55` uses inline error handling:
```go
output.WriteString(fmt.Sprintf("error parsing flags: %v\n", err))
```

But `HandleFlagError` exists in `cmd/utils/output.go:24-27` and is used in `pokemon.go:230`.

**Action:** Replace inline flag error handling with `utils.HandleFlagError()` in ability.go and any other commands.

---

## Low Priority

### 8. Expand CLI test coverage (`cli_test.go:131`)

**Current state:** TODO comment with commented-out test cases:
```go
// TODO: finish testing different commands?
func TestRunCLI_VariousCommands(t *testing.T) {
    tests := []struct{...}{
        //{"Invalid command", []string{"foobar"}, 1},
        //{"Missing Pokemon name", []string{"pokemon"}, 1},
        //{"Another invalid command", []string{"invalid"}, 1},
    }
}
```

**Action:** Uncomment and implement the commented test cases.

---

### 9. Consider removing API call wrapper functions (`connections/connection.go:82-104`)

**Current state:** 6 thin wrappers that just call `FetchEndpoint`:
```go
func AbilityApiCall(endpoint, abilityName, baseURL string) (structs.AbilityJSONStruct, string, error) {
    return FetchEndpoint[structs.AbilityJSONStruct](endpoint, abilityName, baseURL, "Ability")
}
```

**Trade-off:**
- Keep: Slightly cleaner call sites, enforces type safety
- Remove: Less code, callers use `FetchEndpoint` directly

**Recommendation:** Keep for now - provides good API boundary and clear naming.

---

### 10. Normalize hyphen hint text

**Current state:** Similar hints appear with slightly different wording:
- `cli.go:83-85`: "hint: when calling a resource with a space, use a hyphen"
- `ability.go:25`: "Use a hyphen when typing a name with a space."
- `pokemon.go:34`: Same as ability
- `move.go:26`: Same

**Suggestion:** Define a constant:
```go
// In styling or utils
const HyphenHint = "Use a hyphen when typing a name with a space."
```

---

### 11. Consistent flag parsing pattern

**Current state:** `pokemon.go:68-72` uses `os.Exit(1)`:
```go
if err := pokeFlags.Parse(args[3:]); err != nil {
    fmt.Printf("error parsing flags: %v\n", err)
    pokeFlags.Usage()
    os.Exit(1)
}
```

But `ability.go:53-58` returns an error properly:
```go
if err := abilityFlags.Parse(args[3:]); err != nil {
    output.WriteString(fmt.Sprintf("error parsing flags: %v\n", err))
    abilityFlags.Usage()
    return output.String(), err
}
```

**Action:** Standardize on returning errors (avoid `os.Exit` in library code).

---

### 12. Split pokemonflagset.go into smaller files

**Current state:** 680 lines with 6 major functions:
- `SetupPokemonFlagSet()` (43-82)
- `AbilitiesFlag()` (84-135)
- `DefenseFlag()` (137-319)
- `ImageFlag()` (321-408)
- `MovesFlag()` (410-553)
- `StatsFlag()` (555-648)
- `TypesFlag()` (650-680) - marked deprecated

**Suggestion:** Split into:
- `flags/pokemonflagset.go` - just `SetupPokemonFlagSet()`
- `flags/pokemon_abilities.go`
- `flags/pokemon_defense.go`
- `flags/pokemon_image.go`
- `flags/pokemon_moves.go`
- `flags/pokemon_stats.go`

**Trade-off:** More files vs. easier navigation. Lower priority since file is well-organized.

---

## Suggested Refactor Order

For incremental improvement without breaking changes:

1. **Quick wins (< 1 hour each):**
   - [ ] Add `CapitalizeResourceName` to styling package (#6)
   - [ ] Normalize hyphen hint constant (#10)
   - [ ] Use `HandleFlagError` consistently (#7)

2. **Medium effort (1-2 hours each):**
   - [ ] Consolidate simple validators (#1)
   - [ ] Deduplicate command list in cli.go (#2)
   - [ ] Create help flag checker utility (#4)

3. **Larger refactors (2+ hours):**
   - [ ] Return struct from SetupPokemonFlagSet (#3)
   - [ ] Move header styling to styling package (#5)
   - [ ] Standardize flag parsing pattern (#11)

4. **When time permits:**
   - [ ] Expand CLI test coverage (#8)
   - [ ] Consider splitting pokemonflagset.go (#12)

---

## Notes

- All refactors should maintain backward compatibility with existing CLI usage
- Run `go test ./...` after each change
- Consider adding deprecation warnings before removing any public functions
