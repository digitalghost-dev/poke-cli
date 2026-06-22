---
weight: 2
---

# Overview

This section is a behind-the-scenes reference for how individual parts of the CLI work at runtime.

Where the [Infrastructure Guide](../Infrastructure_Guide/index.md) covers how the backend is **built and provisioned**, these pages cover how the CLI **reads and assembles data** when a command runs: 

* The request paths
* The services involved
* Which tables or APIs answer each call 
    
Each page focuses on one command or service and leads with a sequence diagram.

## Commands

How individual CLI commands fetch and assemble their data at runtime.

| Command | Covers |
|------|--------|
| [`comp`](Commands/comp.md) | How `poke-cli comp` loads TCG and VGC tournament standings and other data from Supabase |

## Services

Standalone services the CLI relies on.

| Service | Covers |
|------|--------|
| [Rust Aggregation Service](Services/rust-aggregation-service.md) | How the Go CLI delegates Pokémon data assembly to the `poke-aggregate` Rust binary |
| [Rust Caching Service](Services/rust-caching-service.md) | How the program calls the Rust binary to fetch cached data |
