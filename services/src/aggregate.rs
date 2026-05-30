use serde_json::Value;

pub struct ProfileOptions {
    pub abilities: bool,
    pub defense: bool,
    pub image: Option<String>,
    pub moves: bool,
    pub stats: bool,
}

pub fn run(name: &str, opts: &ProfileOptions) -> anyhow::Result<Value> {
    let value = serde_json::json!({
        "name": name,
        "abilities": opts.abilities,
        "defense": opts.defense,
        "image": opts.image,
        "moves": opts.moves,
        "stats": opts.stats,
    });

    Ok(value)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn empty_opts() -> ProfileOptions {
        ProfileOptions {
            abilities: false,
            defense: false,
            image: None,
            moves: false,
            stats: false,
        }
    }

    #[test]
    fn run_includes_name() {
        let value: Value = run("charizard", &empty_opts()).unwrap();
        assert_eq!(value["name"], "charizard");
    }

    #[test]
    fn run_reflects_requested_flags() {
        let mut opts: ProfileOptions = empty_opts();
        opts.abilities = true;
        opts.stats = true;

        let value: Value = run("charizard", &opts).unwrap();

        assert_eq!(value["abilities"], true);
        assert_eq!(value["stats"], true);
        assert_eq!(value["defense"], false);
        assert_eq!(value["moves"], false);
    }

    #[test]
    fn run_image_none_serializes_as_null() {
        let value: Value = run("charizard", &empty_opts()).unwrap();
        assert!(value["image"].is_null());
    }

    #[test]
    fn run_image_some_serializes_as_string() {
        let mut opts: ProfileOptions = empty_opts();
        opts.image = Some("md".to_string());

        let value = run("charizard", &opts).unwrap();

        assert_eq!(value["image"], "md");
    }
}