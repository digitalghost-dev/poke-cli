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

#[cfg(test)]
mod tests {
    use super::*;
    use std::env;
    use std::fs::{File, FileTimes};
    use std::time::SystemTime;

    // Unique-per-test path inside the OS temp dir so parallel tests don't collide.
    fn temp_path(label: &str) -> PathBuf {
        env::temp_dir().join(format!("poke-cache-test-{}-{}.json", std::process::id(), label))
    }

    #[test]
    fn hash_url_is_deterministic_and_distinct() {
        let a: String = hash_url("https://pokeapi.co/api/v2/pokemon/charizard");
        let b: String = hash_url("https://pokeapi.co/api/v2/pokemon/charizard");
        let c: String = hash_url("https://pokeapi.co/api/v2/pokemon/pikachu");

        assert_eq!(a, b); // same input -> same hash
        assert_ne!(a, c); // different input -> different hash
        assert_eq!(a.len(), 64); // sha256 hex is 64 chars
    }

    #[test]
    fn get_cache_dir_targets_poke_cache() {
        assert!(get_cache_dir().ends_with("poke-cache"));
    }

    #[test]
    fn cache_is_fresh_is_false_for_missing_file() {
        let path: PathBuf = temp_path("missing");
        let _ = fs::remove_file(&path); // ensure it doesn't exist

        assert!(!cache_is_fresh(&path).unwrap());
    }

    #[test]
    fn cache_is_fresh_is_false_for_stale_file() {
        let path: PathBuf = temp_path("stale");
        let _ = fs::remove_file(&path);

        fs::write(&path, "hello").unwrap();
        let stale_time: SystemTime = SystemTime::now() - (TTL + Duration::from_secs(1));
        let file: File = File::options().write(true).open(&path).unwrap();
        file.set_times(FileTimes::new().set_modified(stale_time))
            .unwrap();

        assert!(!cache_is_fresh(&path).unwrap());

        fs::remove_file(&path).unwrap();
    }

    #[test]
    fn write_atomic_writes_content_and_is_fresh() {
        let path: PathBuf = temp_path("roundtrip");
        let _ = fs::remove_file(&path);

        write_atomic(&path, "hello").unwrap();

        assert_eq!(fs::read_to_string(&path).unwrap(), "hello");
        assert!(cache_is_fresh(&path).unwrap()); 

        fs::remove_file(&path).unwrap();
    }
}
