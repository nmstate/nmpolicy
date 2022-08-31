#![allow(dead_code)]
mod ast;
mod capture;
pub mod error;
mod expand;
mod lex;
pub mod operations;
mod parse;
mod resolve;
mod snippet;
pub mod types;

pub use crate::error::NmpolicyError;

#[cfg(test)]
mod snippet_tests;
