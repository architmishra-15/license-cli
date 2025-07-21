use anyhow::Result;
use regex::Regex;
use std::collections::HashMap;

/// Holds license-names + matching regexes
pub struct LicenseMatcher {
    pub licenses: Vec<String>,
    patterns: HashMap<String, Vec<Regex>>,
}

impl LicenseMatcher {
    pub fn new(licenses: Vec<String>) -> Result<Self> {
        let mut patterns = HashMap::new();
        for lic in &licenses {
            let mut regs = Vec::new();
            match lic.as_str() {
                "apache-2.0" => {
                    regs.push(Regex::new(r"(?i)^apache[-_]?2(\.0)?$")?);
                    regs.push(Regex::new(r"(?i)^ap(ache)?[-_]?2(\.0)?$")?);
                }
                "mit" => {
                    regs.push(Regex::new(r"(?i)^mit$")?);
                }
                _ => {
                    regs.push(Regex::new(&format!(r"(?i)^{}$", regex::escape(lic)))?);
                }
            }
            patterns.insert(lic.clone(), regs);
        }
        Ok(Self { licenses, patterns })
    }

    /// Exact or regex-based match
    pub fn find_license(&self, input: &str) -> Option<String> {
        // exact
        for lic in &self.licenses {
            if lic.eq_ignore_ascii_case(input) {
                return Some(lic.clone());
            }
        }
        // regex
        for (lic, regs) in &self.patterns {
            for r in regs {
                if r.is_match(input) {
                    return Some(lic.clone());
                }
            }
        }
        None
    }

    /// Simple substring suggestions
    pub fn suggest_licenses(&self, input: &str) -> Vec<String> {
        let mut out = Vec::new();
        let lower = input.to_lowercase();
        for lic in &self.licenses {
            if lic.to_lowercase().contains(&lower) {
                out.push(lic.clone());
            }
        }
        out
    }
}
