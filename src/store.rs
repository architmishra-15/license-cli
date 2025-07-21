use anyhow::{anyhow, Result};
use std::{env, fs, path::PathBuf};
use crate::downloader;

/// Where we keep cached license templates.
pub fn get_store_dir() -> Result<PathBuf> {
    let base = if cfg!(windows) {
        // e.g. C:\Users\You\AppData\Roaming\license_cli\licenses
        let local = env::var_os("LOCALAPPDATA")
            .map(PathBuf::from)
            .ok_or_else(|| anyhow!("%LOCALAPPDATA% is not set"))?;
            PathBuf::from(local).join("license_cli").join("licenses")
    } else {
        // e.g. /home/you/.license_cli/licenses
        dirs::home_dir()
            .ok_or_else(|| anyhow!("HOME not set"))?
            .join(".license_cli")
            .join("licenses")
    };

    if !base.exists() {
        fs::create_dir_all(&base)?;
    }
    Ok(base)
}

/// List all license templates we have locally (after any downloads)
pub fn discover_licenses() -> Result<Vec<String>> {
    let dir = get_store_dir()?;
    let mut names = Vec::new();

    for entry in fs::read_dir(dir)? {
        let entry = entry?;
        if entry.path().is_file() {
            if let Some(n) = entry.file_name().to_str() {
                names.push(n.to_string());
            }
        }
    }
    Ok(names)
}

/// Ensure that `name` exists locally—if not, try to download it.
/// Returns the full path to the file in the local store.
pub fn ensure_license(name: &str) -> Result<PathBuf> {
    let store = get_store_dir()?;
    let path = store.join(name);

    if path.exists() {
        return Ok(path);
    }

    downloader::download_license(name)
        .map_err(|e| anyhow!("Failed to download license '{}': {}", name, e))?;
    Ok(path)
}
