#![allow(dead_code)]
mod ast;
mod capture;
mod error;
mod expand;
mod lex;
mod operations;
mod parse;
mod resolve;
mod snippet;
mod types;

pub use crate::error::NmpolicyError;

#[cfg(test)]
mod snippet_tests;
