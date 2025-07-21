use anyhow::{anyhow, Result};
use std::collections::HashMap;
use std::fs;
use std::path::PathBuf;
use regex::Regex;
use std::io::{self, Write};

const LINKS_JSON_URL: &str = 
    "https://raw.githubusercontent.com/you/yourrepo/main/licenses/links.json";

/// Download the JSON manifest of license → URL mappings.
fn fetch_links_manifest() -> Result<HashMap<String,String>> {
    let resp = reqwest::blocking::get(LINKS_JSON_URL)
        .map_err(|e| anyhow!("Failed to fetch links.json: {}", e))?;
    if !resp.status().is_success() {
        return Err(anyhow!("links.json not found at {}", LINKS_JSON_URL));
    }
    let text = resp.text()?;
    let map: HashMap<String,String> = serde_json::from_str(&text)
        .map_err(|e| anyhow!("Invalid JSON in links.json: {}", e))?;
    Ok(map)
}

pub fn download_license(name: &str) -> Result<PathBuf> {
    let manifest = fetch_links_manifest()?;

    // 1) direct key match?
    if let Some(url) = manifest.get(name) {
        return do_download(name, url);
    }

    // 2) regex‐based fuzzy search of keys
    let pat = Regex::new(&format!(r"(?i){}", regex::escape(name)))?;
    let candidates = manifest
        .keys()
        .filter(|k| pat.is_match(k))
        .collect::<Vec<_>>();

    if candidates.is_empty() {
        // nothing matched
        return Err(anyhow!(
           "No entry for `{}` in links.json. \
            Please add it or update links.json on GitHub.",
           name
        ));
    }

    // 3) if multiple, let the user choose
    println!("Found these similar licenses:");
    for (i, key) in candidates.iter().enumerate() {
        println!("  {}) {}", i+1, key);
    }
    print!("Download which one? [1-{}]: ", candidates.len());
    io::stdout().flush()?;

    let mut line = String::new();
    io::stdin().read_line(&mut line)?;
    let idx: usize = line.trim().parse()
        .map_err(|_| anyhow!("Invalid choice"))?;
    if idx == 0 || idx > candidates.len() {
        return Err(anyhow!("Choice out of range"));
    }
    let chosen = candidates[idx-1];
    let url = &manifest[chosen];
    println!("Downloading `{}` from:\n  {}", chosen, url);

    do_download(chosen, url)
}

/// Perform the actual HTTP GET + write to cache
fn do_download(name: &str, url: &str) -> Result<PathBuf> {
    let resp = reqwest::blocking::get(url)
        .map_err(|e| anyhow!("HTTP GET failed: {}", e))?;
    if !resp.status().is_success() {
        return Err(anyhow!("{} returned {}", url, resp.status()));
    }
    let body = resp.text()?;

    let store = crate::store::get_store_dir()?;
    let path = store.join(name);
    fs::write(&path, body)?;
    Ok(path)
}