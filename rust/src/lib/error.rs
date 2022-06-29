use std::error::Error;

use crate::snippet::snippet;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
#[non_exhaustive]
pub enum ErrorKind {
    InvalidExpression,
    InvalidPath,
    InvalidEqFilter,
    InvalidReplace,
    InvalidPipe,
    Bug,
    NotImplementedError,
    NotSupportedError,
}

impl std::fmt::Display for NmpolicyError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}: {}", self.kind, self.msg,)
    }
}

impl std::fmt::Display for ErrorKind {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ErrorKind::InvalidExpression => write!(f, "invalid expression"),
            ErrorKind::InvalidPath => write!(f, "invalid path"),
            ErrorKind::InvalidEqFilter => write!(f, "invalid equality filter"),
            ErrorKind::InvalidReplace => write!(f, "invalid replace"),
            ErrorKind::InvalidPipe => write!(f, "invalid pipe"),
            ErrorKind::Bug => write!(f, "bug"),
            ErrorKind::NotImplementedError => write!(f, "not implemented"),
            ErrorKind::NotSupportedError => write!(f, "not supported"),
        }
    }
}

impl Error for NmpolicyError {}

#[derive(Clone, Debug, Eq, PartialEq)]
#[non_exhaustive]
pub struct NmpolicyError {
    kind: ErrorKind,
    msg: String,
}

impl NmpolicyError {
    pub fn new(kind: ErrorKind, msg: String) -> Self {
        Self { kind, msg }
    }

    pub fn kind(&self) -> ErrorKind {
        self.kind
    }

    pub fn msg(&self) -> &str {
        self.msg.as_str()
    }

    pub fn decorate(&mut self, expression: String, position: usize) {
        self.msg = format!("{}\n{}", self.msg, snippet(expression, position).as_str())
    }
}
