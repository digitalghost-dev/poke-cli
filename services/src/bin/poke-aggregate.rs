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

impl From<PokemonArgs> for ProfileOptions {
    fn from(args: PokemonArgs) -> Self {
        Self {
            abilities: args.abilities,
            defense: args.defense,
            image: args.image.map(|size| match size {
                ImageSize::Sm => "sm",
                ImageSize::Md => "md",
                ImageSize::Lg => "lg",
            }.to_string()),
            moves: args.moves,
            stats: args.stats,
        }
    }
}

fn main() -> anyhow::Result<()> {
    let cli: Cli = Cli::parse();

    match cli.command{
        SubCommands::Pokemon(args) => {
            let name: String = args.name.clone();
            let options: ProfileOptions = args.into();

            let value: serde_json::Value = run(&name, &options)?;

            serde_json::to_writer_pretty(std::io::stdout(), &value)?;
            println!();
        }
    }

    Ok(())
}
