use std::env;
use std::fs;

use anyhow::bail;
use services::cache::{cache_is_fresh, get_cache_dir, hash_url, write_atomic};

fn main() -> anyhow::Result<()> {
    let args: Vec<String> = env::args().collect();

    if args.len() != 3 || args[1] != "get" {
        bail!("Usage: poke-cache get <url>");
    }

    let url = &args[2];
    let filename = hash_url(url) + ".json";
    let cache_dir = get_cache_dir();

    fs::create_dir_all(&cache_dir)?;

    let cache_path = cache_dir.join(filename);

    if cache_is_fresh(&cache_path)? {
        eprintln!("Cache hit. Reading: {cache_path:?}");
        let body = fs::read_to_string(&cache_path)?;
        println!("{}", body);

        return Ok(());
    }

    if cache_path.exists() {
        eprintln!("Stale cache. Fetching: {url}");
    } else {
        eprintln!("Cache miss. Fetching: {url}");
    }

    let response_body = reqwest::blocking::get(url)?.error_for_status()?.text()?;

    write_atomic(&cache_path, &response_body)?;

    println!("{}", response_body);

    Ok(())
}
