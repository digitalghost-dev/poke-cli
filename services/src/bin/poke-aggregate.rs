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

            let profile: services::domain::Pokemon = run(&name, &options)?;

            serde_json::to_writer_pretty(std::io::stdout(), &profile)?;
            println!();
        }
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_chained_flags() {
        let cli = Cli::try_parse_from([
            "poke-aggregate", "pokemon", "charizard", "-a", "-s", "--image=md",
        ])
        .unwrap();

        match cli.command {
            SubCommands::Pokemon(args) => {
                assert_eq!(args.name, "charizard");
                assert!(args.abilities);
                assert!(args.stats);
                assert!(!args.defense);
                assert!(!args.moves);
                assert!(matches!(args.image, Some(ImageSize::Md)));
            }
        }
    }

    #[test]
    fn image_equals_and_space_forms_are_equivalent() {
        let equals = Cli::try_parse_from([
            "poke-aggregate", "pokemon", "charizard", "--image=lg",
        ])
        .unwrap();
        let spaced = Cli::try_parse_from([
            "poke-aggregate", "pokemon", "charizard", "--image", "lg",
        ])
        .unwrap();

        for cli in [equals, spaced] {
            match cli.command {
                SubCommands::Pokemon(args) => {
                    assert!(matches!(args.image, Some(ImageSize::Lg)));
                }
            }
        }
    }

    #[test]
    fn rejects_bad_image_size() {
        let result = Cli::try_parse_from([
            "poke-aggregate", "pokemon", "charizard", "--image", "xl",
        ]);

        assert!(result.is_err());
    }

    #[test]
    fn missing_name_is_rejected() {
        let result = Cli::try_parse_from(["poke-aggregate", "pokemon"]);

        assert!(result.is_err());
    }

    #[test]
    fn from_pokemon_args_maps_image_enum_to_string() {
        let cli = Cli::try_parse_from([
            "poke-aggregate", "pokemon", "charizard", "--image=md", "-a",
        ])
        .unwrap();

        match cli.command {
            SubCommands::Pokemon(args) => {
                let opts: ProfileOptions = args.into();
                assert!(opts.abilities);
                assert_eq!(opts.image, Some("md".to_string()));
            }
        }
    }
}