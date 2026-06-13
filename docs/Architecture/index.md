---
weight: 2
---

# Overview

This section is a behind-the-scenes reference for how individual parts of the CLI work at runtime.

Where the [Infrastructure Guide](../Infrastructure_Guide/index.md) covers how the backend is **built and provisioned**, these pages cover how the CLI **reads and assembles data** when a command runs: the request paths, the services involved, and which tables or APIs answer each call. Each page focuses on one command or service and leads with a sequence diagram.

## Commands

How individual CLI commands fetch and assemble their data at runtime.

| Page | Covers |
|------|--------|
| [`comp` TCG Standings Flow](Commands/comp-standings-flow.md) | How `poke-cli comp` loads tournament standings from Supabase via `comp_standings_view` |

## Services

Standalone services the CLI relies on.

| Page | Covers |
|------|--------|
| [Rust Aggregation Service](Services/rust-aggregation-service.md) | How the Go CLI delegates Pokémon data assembly to the `poke-aggregate` Rust binary |
