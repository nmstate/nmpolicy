use serde::{Deserialize, Serialize};

use crate::{
    ast::node::Node,
    error::{ErrorKind, NmpolicyError},
    lex::tokens::Tokens,
    parse::parser::Parser,
    resolve::resolver::Resolver,
    types::{CapturedStates, NMState},
};
use std::collections::HashMap;

use serde_json::Value;

pub(crate) type Capture = HashMap<String, CaptureEntry>;

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub(crate) struct CaptureEntry {
    pub expression: String,
    pub ast: Box<Node>,
}

fn resolve_entries(
    capture: Capture,
    current_state: NMState,
    cache: Option<CapturedStates>,
) -> Result<CapturedStates, NmpolicyError> {
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

fn resolve_entry_path(
    expression: String,
    captured_states: CapturedStates,
) -> Result<Value, NmpolicyError> {
    let capture_entry = parse(expression)?;
    let mut resolver = Resolver::from_capture_entry_and_captured(
        capture_entry.expression,
        capture_entry.ast,
        captured_states,
    );
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
