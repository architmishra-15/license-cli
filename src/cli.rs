use clap::{Parser, Subcommand};

/// CLI definition
#[derive(Parser)]
#[command(author, version, about = "A CLI tool for managing project licenses")]
pub struct Cli {
    #[command(subcommand)]
    pub command: Option<Commands>,
}

#[derive(Subcommand)]
pub enum Commands {
    List,
    Add {
        /// License name or pattern to match
        license: String,
        /// Author name (overrides git config)
        #[arg(short, long)]
        author: Option<String>,
    },
}
