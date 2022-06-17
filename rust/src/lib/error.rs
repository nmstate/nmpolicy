use std::error::Error;

use crate::snippet::snippet;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
#[non_exhaustive]
pub enum ErrorKind {
    InvalidArgument,
    Bug,
    VerificationError,
    NotImplementedError,
    NotSupportedError,
}

impl std::fmt::Display for ErrorKind {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{:?}", self)
    }
}

impl std::fmt::Display for NmpolicyError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}: {}", self.kind, self.msg,)
    }
}

impl Error for NmpolicyError {}

#[derive(Debug, Eq, PartialEq)]
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
