#![allow(dead_code)]
mod ast;
mod error;
mod lex;
mod snippet;

pub use crate::error::NmpolicyError;

#[cfg(test)]
mod snippet_tests;
