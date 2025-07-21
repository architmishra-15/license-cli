use clap::{Parser, Subcommand};

/// CLI definition
#[derive(Parser)]
#[command(disable_help_flag = true, author = "Archit Mishra", version, about = "A CLI tool for managing project licenses")]
pub struct Cli {

    #[arg(short = 'h', long = "help", action)]
    pub help: bool,

    #[arg(short = 'v', long = "version", action)]
    pub version: bool,


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
    Version

}