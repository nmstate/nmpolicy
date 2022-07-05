use std::error::Error;

use crate::snippet::snippet;

#[derive(Debug, Clone, PartialEq, Eq)]
#[non_exhaustive]
pub enum ErrorKind {
    // Tokenizer
    IllegalChar(char),
    InvalidEqFilterMissingEqual(char),
    InvalidEqFilterEOF,
    InvalidReplaceMissingEqual(char),
    InvalidReplaceEOF,
    InvalidStringMissingDelimiter(char),
    InvalidNumberFormat(char),
    InvalidIdentityFormat(char),

    // Parser
    InvalidExpressionUnexpectedToken(String),
    InvalidPipeMissingLeftPath,
    InvalidPipeMissingLeftExpression,
    InvalidPipeMissingRightExpression,
    InvalidPathUnexpectedTokenAfterDot,
    InvalidPathMissingDot,
    InvalidTernaryUnexpectedRightHand(&'static str),
    InvalidTernaryMissingRightHand(&'static str),
    InvalidTernaryUnexpectedLeftHand(&'static str),
    InvalidTernaryMissingLeftHand(&'static str),

    // Resolver
    ResolveErrorCaptureEntryNotFound(String),
    ResolveErrorNotSupportedValue,
    ResolveErrorInvalidPathInputSource(String),
    ResolveErrorInvalidInputSource(String),
    ResolveErrorUnsupportedOperation(String),
    ResolveErrorFailingFilteringPath(String),
    ResolveErrorFailedFilterResultConvertion(serde_json::Value),
    PathErrorUnexpectedMapStep(serde_json::Map<String, serde_json::Value>),
    PathErrorStepInvalidType(serde_json::Value, String),

    Bug,
    NotImplementedError,
    NotSupportedError,
}

impl std::fmt::Display for NmpolicyError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        if self.expression.is_empty() {
            write!(f, "{}", self.kind)
        } else {
            write!(
                f,
                "{}\n{}",
                self.kind,
                snippet(self.expression.clone(), self.pos)
            )
        }
    }
}

impl std::fmt::Display for ErrorKind {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ErrorKind::IllegalChar(ch) => write!(f, "illegal char {ch}"),
            ErrorKind::InvalidEqFilterMissingEqual(ch) => write!(
                f,
                "invalid EQFILTER operation format ({ch} is not equal char)"
            ),
            ErrorKind::InvalidEqFilterEOF => write!(f, "invalid EQFILTER operation format (EOF)"),
            ErrorKind::InvalidReplaceMissingEqual(ch) => write!(
                f,
                "invalid REPLACE operation format ({ch} is not equal char)"
            ),
            ErrorKind::InvalidReplaceEOF => write!(f, "invalid REPLACE operation format (EOF)"),
            ErrorKind::InvalidStringMissingDelimiter(delimiter) => {
                write!(f, "invalid string format (missing {delimiter} terminator)")
            }
            ErrorKind::InvalidNumberFormat(ch) => {
                write!(f, "invalid number format ({ch} is not a digit)")
            }
            ErrorKind::InvalidIdentityFormat(ch) => write!(
                f,
                "invalid identity format ({ch} is not a digit, letter or -)"
            ),
            ErrorKind::InvalidExpressionUnexpectedToken(literal) => {
                write!(f, "invalid expression: unexpected token `{}`", literal)
            }
            ErrorKind::InvalidPipeMissingLeftPath => {
                write!(f, "invalid pipe: only paths can be piped in",)
            }
            ErrorKind::InvalidPipeMissingLeftExpression => {
                write!(f, "invalid pipe: missing pipe in expression",)
            }
            ErrorKind::InvalidPipeMissingRightExpression => {
                write!(f, "invalid pipe: missing pipe out expression",)
            }
            ErrorKind::InvalidPathUnexpectedTokenAfterDot => {
                write!(f, "invalid path: missing identity or number after dot",)
            }
            ErrorKind::InvalidPathMissingDot => {
                write!(f, "invalid path: missing dot",)
            }
            ErrorKind::InvalidTernaryUnexpectedRightHand(operator_name) => {
                write!(
                    f,
                    "invalid {operator_name}: right hand argument is not a string or identity"
                )
            }
            ErrorKind::InvalidTernaryUnexpectedLeftHand(operator_name) => {
                write!(
                    f,
                    "invalid {operator_name}: left hand argument is not a path"
                )
            }
            ErrorKind::InvalidTernaryMissingRightHand(operator_name) => {
                write!(f, "invalid {operator_name}: missing right hand argument")
            }
            ErrorKind::InvalidTernaryMissingLeftHand(operator_name) => {
                write!(f, "invalid {operator_name}: missing left hand argument")
            }
            ErrorKind::ResolveErrorCaptureEntryNotFound(capture_entry_name) => {
                write!(f, "capture entry '{capture_entry_name}' not found")
            }
            ErrorKind::ResolveErrorNotSupportedValue => {
                write!(
                    f,
                    "not supported value. Only string or capture entry path are supported"
                )
            }
            ErrorKind::ResolveErrorInvalidPathInputSource(node) => {
                write!(
                    f,
                    "invalid path input source ({node}), only capture reference is supported"
                )
            }
            ErrorKind::ResolveErrorInvalidInputSource(node) => {
                write!(f, "invalid input source ({node}), only current state or capture reference is supported")
            }
            ErrorKind::ResolveErrorUnsupportedOperation(node) => {
                write!(f, "root node has unsupported operation : {node}")
            }
            ErrorKind::ResolveErrorFailingFilteringPath(error) => {
                write!(f, "failed applying operation on the path: {error}")
            }
            ErrorKind::ResolveErrorFailedFilterResultConvertion(value) => {
                write!(f, "failed converting filtering result `{value}` to a map")
            }
            ErrorKind::PathErrorUnexpectedMapStep(map) => {
                write!(
                    f,
                    "path error: unexpected non identity step for map state '{map:?}'"
                )
            }
            ErrorKind::PathErrorStepInvalidType(input_state, step) => {
                write!(f, "invalid type {input_state:?} for identity step '{step}'")
            }

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
    expression: String,
    pos: usize,
}

impl NmpolicyError {
    pub fn new(kind: ErrorKind) -> Self {
        Self {
            kind,
            expression: "".to_string(),
            pos: 0,
        }
    }

    pub fn kind(&self) -> ErrorKind {
        self.kind.clone()
    }

    pub fn decorate(&mut self, expression: String, pos: usize) {
        self.expression = expression;
        self.pos = pos;
    }
}
