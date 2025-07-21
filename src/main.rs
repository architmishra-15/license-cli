mod cli;
mod matcher;
mod store;
mod downloader;


use anyhow::Result;
use chrono::Datelike;
use clap::Parser;
use colored::{Colorize, control::{set_override}};
use std::{fs, io::{self, Write}, process::Command};

use crate::cli::{Cli, Commands};
use crate::matcher::LicenseMatcher;
use crate::store::{discover_licenses, ensure_license};

fn get_git_author() -> Result<String> {
    let out = Command::new("git")
        .args(["config", "user.name"])
        .output()?;
    if out.status.success() {
        let s = String::from_utf8_lossy(&out.stdout).trim().to_string();
        if !s.is_empty() {
            return Ok(s);
        }
    }
    Err(anyhow::anyhow!("Cannot determine git user.name"))
}

fn get_current_year() -> i32 {
    chrono::Utc::now().year()
}

fn project_root() -> Result<std::path::PathBuf> {
    match Command::new("git").args(["rev-parse", "--show-toplevel"]).output() {
        Ok(o) if o.status.success() => {
            let p = String::from_utf8_lossy(&o.stdout).trim().to_string();
            Ok(std::path::PathBuf::from(p))
        }
        _ => std::env::current_dir().map_err(Into::into),
    }
}

fn process_license_text(template: &str, author: &str, year: i32) -> String {
    template
        .replace("[yyyy]", &year.to_string())
        .replace("[year]", &year.to_string())
        .replace("[name of the author]", author)
        .replace("[fullname]", author)
        .replace("[author]", author)
}

fn add_license(name: &str, author: Option<String>) -> Result<()> {
    // ensure we have the template locally
    let path = ensure_license(name)?;
    let content = fs::read_to_string(&path)?;
    let author = author.unwrap_or_else(|| get_git_author().unwrap());
    let processed = process_license_text(&content, &author, get_current_year());

    let root = project_root()?;
    let dest = root.join("LICENSE");
    if dest.is_file() {
        eprint!("LICENSE exists at {}. Overwrite? [y/N]: ", dest.display());
        io::stdout().flush()?;
        let mut ans = String::new();
        io::stdin().read_line(&mut ans)?;
        if !ans.trim().eq_ignore_ascii_case("y") {
            println!("Cancelled.");
            return Ok(());
        }
    }

    fs::write(&dest, processed)?;
    println!("✓ Installed `{}` license", name.green());
    Ok(())
}

fn print_help() {
    // Header
    println!();
    println!("{}", "License CLI Tool".bold().underline().blue());
    println!();

    // Synopsis
    println!(
        "{}\n  {}",
        "Usage:".bold().bright_yellow(),
        "license [OPTIONS] <SUBCOMMAND>"
    );
    println!();

    // Global options
    println!("{}", "Global Options:".bold().bright_green());
    println!(
        "  {}  {}    {}",
        "-h, --help".cyan(),
        "Print this help message".white(),
        ""
    );
    println!(
        "  {}  {}    {}",
        "-v, --version".cyan(),
        "Show version and author info".white(),
        ""
    );
    println!();

    // Subcommands
    println!("{}", "Subcommands:".bold().bright_green());

    // Show
    println!(
        "  {}        {}",
        "show".magenta().bold(),
        "List all locally cached licenses".white()
    );

    // Add
    println!(
        "  {} <name> {}",
        "add".magenta().bold(),
        "Install a license into your project".white()
    );
    println!(
        " {} {}",
        "--author <you>".cyan(),
        "(override author field)".white()
    );
    println!();

    // Examples
    println!("{}", "Examples:".bold().bright_green());
    println!(
        "  {}",
        "license show".yellow()
    );
    println!(
        "  {}",
        "license add mit".yellow()
    );
    println!(
        "  {}",
        "license add apache-2.0 --author \"Alice\"".yellow()
    );
    println!();

    // Footer
    println!(
        "{}  {}",
        "Repository:".bold(),
        "https://github.com/architmishra-15/license-cli".underline().blue()
    );
    println!(
        "{}  {}",
        "Author:".bold(),
        "Archit Mishra".cyan()
    );
    println!(
        "{}  {}",
        "Version:".bold(),
        env!("CARGO_PKG_VERSION").green()
    );
    println!();
}

fn main() -> Result<()> {
    set_override(true);

    let cli = Cli::parse();

    if cli.help {
        print_help();
        return Ok(());
    }

    if cli.version {
        println!("license_cli {}", env!("CARGO_PKG_VERSION").green());
        println!("Author: {}", "Archit Mishra".bold().italic().green());
        println!("License: {}", "GPL-3.0".green());
        println!("GitHub: {}", "https://github.com/architmishra-15/license-cli".bright_blue());
        return Ok(());
    }

    match cli.command {
        Some(Commands::List) => {
            println!("{}", "Available licenses:".bold());
            for lic in discover_licenses()? {
                println!("  {}", lic);
            }
            Ok(())
        }

        Some(Commands::Add { license, author }) => {
            let matcher = LicenseMatcher::new(discover_licenses()?)?;
            if let Some(found) = matcher.find_license(&license) {
                add_license(&found, author)
            } else {
                eprintln!("{}", format!("License '{}' not found.", license).red());
                let sugg = matcher.suggest_licenses(&license);
                if !sugg.is_empty() {
                    eprintln!("\n{}", "Did you mean:".yellow());
                    for s in sugg {
                        eprintln!("  {}", s.green());
                    }
                }
                std::process::exit(1);
            }
        }

        Some(Commands::Version) => {
            println!("license_cli {}", env!("CARGO_PKG_VERSION").green());
            println!("Author: {}", "Archit Mishra".bold().italic().green());
            println!("License: {}", "GPL-3.0".green());
            println!("GitHub: {}", "https://github.com/architmishra-15/license-cli".bright_blue());
        return Ok(());
        }

        None => {
            println!("{}", "License CLI Tool".bold());
            println!();
            println!("Usage:");
            println!("  license list                       Show all available licenses");
            println!("  license add <name>                 Add a license to current directory");
            println!("  license add <name> --author=<you>  Add license with custom author");
            println!();
            println!("Examples:");
            println!("  license list");
            println!("  license add mit");
            println!("  license add apache-2.0 --author=\"Alice\"");
            println!();
            Ok(())
        }
    }
}
