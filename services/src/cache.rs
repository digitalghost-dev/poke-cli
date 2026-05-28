use std::fs;
use std::path::{Path, PathBuf};
use std::time::Duration;

use anyhow::Result;
use sha2::{Digest, Sha256};

pub const TTL: Duration = Duration::from_secs(60 * 60 * 24); // 24 hours

pub fn hash_url(url: &str) -> String {
    let mut hasher = Sha256::new();
    hasher.update(url.as_bytes());

    hex::encode(hasher.finalize())
}

pub fn get_cache_dir() -> PathBuf {
    dirs::cache_dir()
        .unwrap_or_else(|| PathBuf::from("poke-cache"))
        .join("poke-cache")
}

pub fn cache_is_fresh(path: &PathBuf) -> anyhow::Result<bool> {
    if !path.exists() {
        return Ok(false);
    }

    let metadata = fs::metadata(path)?;
    let age = metadata.modified()?.elapsed()?;

    Ok(age < TTL)
}

pub fn write_atomic(path: &Path, data: &str) -> Result<()> {
    let temp_path = path.with_extension("tmp");
    fs::write(&temp_path, data)?;
    fs::rename(&temp_path, path)?;

    Ok(())
}
