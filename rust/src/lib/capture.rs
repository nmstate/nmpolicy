use serde::Deserialize;

use crate::{
    ast::node::Node,
    error::{ErrorKind, NmpolicyError},
    lex::tokens::Tokens,
    parse::parser::Parser,
    resolve::resolver::{CapturePathResolver, Resolver},
    types::{CapturedStates, NMState},
};
use std::collections::HashMap;

use serde_json::Value;

pub type Capture = HashMap<String, CaptureEntry>;

#[derive(Clone, Eq, PartialEq, Debug, Deserialize)]
#[serde(transparent)]
pub struct CaptureEntry {
    #[serde(flatten)]
    pub expression: String,
    #[serde(skip)]
    pub(crate) ast: Box<Node>,
}

pub(crate) fn resolve_entries(
    capture: Capture,
    current_state: NMState,
    cache: Option<CapturedStates>,
) -> Result<CapturedStates, NmpolicyError> {
    if capture.is_empty() || current_state.is_empty() && cache.is_none() {
        return Ok(CapturedStates::new());
    }
    let filtered_cache = filter_cache_by_capture(cache.clone(), capture.clone());
    let filtered_capture = filter_capture_by_cache(capture, cache);
    let parsed_capture = filtered_capture
        .into_iter()
        .map(|(k, v)| match parse(v.expression) {
            Ok(pv) => Ok((k, pv)),
            Err(e) => Err(e),
        })
        .collect::<Result<Capture, NmpolicyError>>()?;

    let mut resolver = Resolver::new(parsed_capture);
    let captured_states = resolver.resolve(current_state, filtered_cache)?;
    Ok(captured_states)
}

pub(crate) fn resolve_entry_path(
    expression: String,
    captured_states: CapturedStates,
) -> Result<Value, NmpolicyError> {
    let capture_entry = parse(expression)?;
    let mut resolver =
        Resolver::from_captured(capture_entry.expression, capture_entry.ast, captured_states);
    resolver.resolve_capture_entry_path()
}

fn parse(expression: String) -> Result<CaptureEntry, NmpolicyError> {
    let tokens = &mut Tokens::new(expression.as_str());
    match Parser::new(expression.clone(), tokens).parse() {
        Ok(Some(ast_root)) => Ok(CaptureEntry {
            expression,
            ast: ast_root,
        }),
        Ok(None) => Err(NmpolicyError::new(ErrorKind::Bug)),
        Err(e) => Err(e),
    }
}

fn filter_cache_by_capture(
    cache_op: Option<CapturedStates>,
    capture: Capture,
) -> Option<CapturedStates> {
    cache_op.map(|cache| {
        cache
            .into_iter()
            .filter(|(k, _)| capture.get(k).is_some())
            .collect()
    })
}

fn filter_capture_by_cache(capture: Capture, cache_op: Option<CapturedStates>) -> Capture {
    match cache_op {
        Some(cache) => capture
            .into_iter()
            .filter(|(k, _)| cache.get(k).is_some())
            .collect(),
        None => capture,
    }
}

pub(crate) struct CaptureEntryResolver {
    captured_states: CapturedStates,
}

impl CaptureEntryResolver {
    pub(crate) fn new(captured_states: CapturedStates) -> Self {
        Self { captured_states }
    }
}

impl CapturePathResolver for CaptureEntryResolver {
    fn resolve_capture_entry_path(&self, capture_path: String) -> Result<Value, NmpolicyError> {
        resolve_entry_path(capture_path, self.captured_states.clone())
    }
}
