use clap::{Parser, Subcommand, ValueEnum};
use services::aggregate::{run, ProfileOptions};

#[derive(Parser)]
#[command(name = "poke-aggregate", version)]
struct Cli {
    #[command(subcommand)]
    command: SubCommands,
}

#[derive(Subcommand)]
enum SubCommands {
    Pokemon(PokemonArgs),
}

#[derive(Parser)]
struct PokemonArgs {
    name: String,

    #[arg(short = 'a', long)]
    abilities: bool,

    #[arg(short = 'd', long)]
    defense: bool,

    #[arg(short = 'i', long, value_enum)]
    image: Option<ImageSize>,

    #[arg(short = 'm', long)]
    moves: bool,

    #[arg(short = 's', long)]
    stats: bool,
}

#[derive(Clone, Copy, ValueEnum)]
enum ImageSize {
    Sm,
    Md,
    Lg,
}

fn main() -> anyhow::Result<()> {
    let cli: Cli = Cli::parse();

    match cli.command{
        SubCommands::Pokemon(args) => {
            println!("Pokémon name: {}", args.name);
            println!("Abilities: {}", args.abilities);
        }
    }
    
    Ok(())
}