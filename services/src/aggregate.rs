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