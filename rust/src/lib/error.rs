use std::error::Error;

use crate::snippet::snippet;

#[derive(Debug, Clone, PartialEq, Eq)]
#[non_exhaustive]
pub enum ErrorKind {
    ValidationError,
    EvaluationError,
    Bug,
    NotImplementedError,
    NotSupportedError,
}

impl std::fmt::Display for NmpolicyError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self.expression.clone() {
            Some(expression) => {
                write!(f, "{}\n{}", self.msg, snippet(expression, self.pos))
            }
            None => write!(f, "{}", self.msg),
        }
    }
}

impl Error for NmpolicyError {}

#[derive(Clone, Debug, Eq, PartialEq)]
pub struct NmpolicyError {
    kind: ErrorKind,
    expression: Option<String>,
    pos: usize,
    msg: String,
}

impl NmpolicyError {
    pub fn new(kind: ErrorKind) -> Self {
        Self {
            kind,
            expression: None,
            pos: 0,
            msg: String::from(""),
        }
    }

    pub fn kind(&self) -> ErrorKind {
        self.kind.clone()
    }

    pub fn decorate(self, expression: String, pos: usize) -> Self {
        self.expression(Some(expression)).pos(pos)
    }

    pub fn ctx(mut self, ctx: String) -> Self {
        if !self.msg.contains(ctx.as_str()) {
            self.msg = format!("{}: {}", ctx, self.msg);
        }
        self
    }

    pub fn pos(mut self, pos: usize) -> Self {
        if self.pos == 0 {
            self.pos = pos;
        }
        self
    }

    pub fn expression(mut self, expression: Option<String>) -> Self {
        self.expression = expression;
        self
    }

    pub fn resolver(self) -> Self {
        self.ctx("resolve error".to_string())
    }
    pub fn eqfilter(self) -> Self {
        self.ctx("eqfilter error".to_string())
    }
    pub fn replace(self) -> Self {
        self.ctx("replace error".to_string())
    }
    pub fn path(self, pos: usize) -> Self {
        self.pos(pos).ctx("invalid path".to_string())
    }
}

impl From<String> for NmpolicyError {
    fn from(error: String) -> NmpolicyError {
        evaluation_error(error)
    }
}

impl From<serde_json::Error> for NmpolicyError {
    fn from(error: serde_json::Error) -> NmpolicyError {
        evaluation_error(error.to_string())
    }
}

impl From<regex::Error> for NmpolicyError {
    fn from(error: regex::Error) -> NmpolicyError {
        evaluation_error(error.to_string())
    }
}

impl From<std::io::Error> for NmpolicyError {
    fn from(error: std::io::Error) -> NmpolicyError {
        evaluation_error(error.to_string())
    }
}

impl From<serde_yaml::Error> for NmpolicyError {
    fn from(error: serde_yaml::Error) -> NmpolicyError {
        evaluation_error(error.to_string())
    }
}

pub(crate) fn error(kind: ErrorKind, msg: String) -> NmpolicyError {
    NmpolicyError {
        kind,
        msg,
        expression: None,
        pos: 0,
    }
}

pub(crate) fn validation_error(msg: String) -> NmpolicyError {
    error(ErrorKind::ValidationError, msg)
}

pub(crate) fn evaluation_error(msg: String) -> NmpolicyError {
    error(ErrorKind::EvaluationError, msg)
}
