use crate::{
    ast::node::{current_state_identity, Node, NodeKind, TernaryOperator},
    error::{evaluation_error, NmpolicyError},
    resolve::{filter, path::Path, replace, walk},
    types::{Capture, CapturedState, CapturedStates, NMState},
};

use serde_json::Value;

pub(crate) struct Resolver {
    capture: Capture,
    current_state: Option<NMState>,
    captured_states: CapturedStates,
    current_node: Option<Box<Node>>,
    current_expression: Option<String>,
}

impl Resolver {
    pub(crate) fn new(capture: Capture) -> Self {
        Self {
            capture,
            current_state: None,
            captured_states: CapturedStates::new(),
            current_node: None,
            current_expression: None,
        }
    }

    pub(crate) fn resolve(
        &mut self,
        current_state: NMState,
        cache: Option<CapturedStates>,
    ) -> Result<CapturedStates, NmpolicyError> {
        self.current_state = Some(current_state);
        if let Some(c) = cache {
            self.captured_states = c
        };
        for capture_entry_name in self.capture.clone().keys() {
            if let Err(mut e) = self.resolve_capture_entry_by_name(capture_entry_name) {
                e = e.expression(self.current_expression.clone()).resolver();
                match self.current_node.clone() {
                    Some(current_node) => return Err(e.pos(current_node.pos)),
                    None => return Err(e),
                }
            }
        }

        Ok(self.captured_states.clone())
    }

    fn resolve_capture_entry_by_name(
        &mut self,
        capture_entry_name: &String,
    ) -> Result<NMState, NmpolicyError> {
        if let Some(captured_state_entry) = self.captured_states.get(capture_entry_name) {
            return Ok(captured_state_entry.clone().state);
        }
        match self.capture.get(capture_entry_name) {
            Some(capture_entry) => {
                self.current_expression = Some(capture_entry.expression.clone());
                self.current_node = Some(capture_entry.ast.clone());
            }
            _ => {
                return Err(evaluation_error(format!(
                    "capture entry '{}' not found",
                    capture_entry_name.clone()
                )));
            }
        }
        let resolved_state = self.resolve_current_capture_entry()?;
        self.captured_states.insert(
            capture_entry_name.to_string(),
            CapturedState {
                state: resolved_state.clone(),
                meta_info: None,
            },
        );
        Ok(resolved_state)
    }

    fn resolve_current_capture_entry(&mut self) -> Result<NMState, NmpolicyError> {
        match self.current_node.clone() {
            Some(current_node) => match current_node.kind {
                NodeKind::EqFilter(lhs, ms, rhs) => match self.resolve_eqfilter((lhs, ms, rhs)) {
                    Ok(state) => Ok(state),
                    Err(e) => Err(e.eqfilter()),
                },
                NodeKind::Replace(lhs, ms, rhs) => match self.resolve_replace((lhs, ms, rhs)) {
                    Ok(state) => Ok(state),
                    Err(e) => Err(e.replace()),
                },
                NodeKind::Path(_) => self.resolve_path_filter(*current_node),
                _ => Err(evaluation_error(format!(
                    "root node has unsupported operation : {current_node}"
                ))),
            },
            None => Ok(NMState::new()),
        }
    }

    fn resolve_eqfilter(&mut self, operator: TernaryOperator) -> Result<NMState, NmpolicyError> {
        let (input_source, path, value) = self.resolve_ternary_operator(operator)?;
        filter::visit_state(input_source, path, value)
    }

    fn resolve_replace(&mut self, operator: TernaryOperator) -> Result<NMState, NmpolicyError> {
        let (input_source, path, value) = self.resolve_ternary_operator(operator)?;
        replace::visit_state(input_source, path, value)
    }

    fn resolve_path_filter(&mut self, node: Node) -> Result<NMState, NmpolicyError> {
        let path = Path::compose_from_node(node)?;
        let current_state = self.current_state.clone().unwrap();
        filter::visit_state(current_state, path, Value::Null)
    }
    fn resolve_ternary_operator(
        &mut self,
        operator: TernaryOperator,
    ) -> Result<(NMState, Path, Value), NmpolicyError> {
        let operator_node = self.current_node.clone();
        let (lhs, ms, rhs) = operator;

        self.current_node = Some(lhs);
        let input_source = self.resolve_input_source()?;

        self.current_node = Some(ms.clone());
        let path = Path::compose_from_node(*ms)?;

        self.current_node = Some(rhs.clone());
        let value: Value = match rhs.kind {
            NodeKind::Str(string) => Value::String(string),
            NodeKind::Path(_) => self.resolve_capture_entry_path()?,
            _ => {
                return Err(evaluation_error(
                    "not supported value. Only string or capture entry path are supported"
                        .to_string(),
                ));
            }
        };
        self.current_node = operator_node;
        Ok((input_source, path, value))
    }

    fn resolve_capture_entry_path(&mut self) -> Result<Value, NmpolicyError> {
        let current_node = self.current_node.clone().unwrap();
        let path = Path::compose_from_node(*current_node)?;
        match path.clone().capture_entry_name {
            Some(capture_entry_name) => {
                let captured_state_entry = self.resolve_capture_entry_by_name(&capture_entry_name)?;
                walk::visit_state(captured_state_entry, path)
            }
            None => Err(evaluation_error("not supported filtered value path. Only paths with a capture entry reference are supported".to_string()))
        }
    }

    fn resolve_input_source(&mut self) -> Result<NMState, NmpolicyError> {
        let current_node = self.current_node.clone().unwrap();
        match current_node.kind {
            n if n == current_state_identity() => Ok(self.current_state.clone().unwrap()),
            _ => match Path::compose_from_node(*current_node.clone()) {
                Ok(path) => match path.capture_entry_name {
                    Some(capture_entry_name) => {
                        self.resolve_capture_entry_by_name(&capture_entry_name)
                    }
                    None => {
                        return Err(evaluation_error(format!(
                            "invalid path input source ({}), only capture reference is supported",
                            current_node
                        )))
                    }
                },
                Err(_) => {
                    return Err(evaluation_error(format!("invalid input source ({}), only current state or capture reference is supported", current_node)));
                }
            },
        }
    }
}
