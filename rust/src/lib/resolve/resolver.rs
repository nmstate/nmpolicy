use crate::{
    ast::node::{current_state_identity, Node, NodeKind, TernaryOperator},
    error::{ErrorKind, NmpolicyError},
    resolve::{filter, path::Path},
    types::{Capture, CapturedState, CapturedStates, NMState},
};

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
            if let Err(e) = self.resolve_capture_entry_by_name(capture_entry_name) {
                return Err(e);
            }
        }

        Ok(self.captured_states.clone())
    }

    fn resolve_capture_entry_by_name(
        &mut self,
        capture_entry_name: &String,
    ) -> Result<NMState, NmpolicyError> {
        match self.capture.get(capture_entry_name) {
            Some(capture_entry) => {
                self.current_expression = Some(capture_entry.expression.clone());
                self.current_node = Some(capture_entry.ast.clone());
            }
            _ => {
                return Err(NmpolicyError::new(
                    ErrorKind::ResolveErrorCaptureEntryNotFound(capture_entry_name.clone()),
                ));
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
                NodeKind::EqFilter(lhs, ms, rhs) => self.resolve_eqfilter((lhs, ms, rhs)),
                _ => Err(NmpolicyError::new(
                    ErrorKind::ResolveErrorUnsupportedOperation(format!("{:?}", current_node)),
                )),
            },
            None => Ok(NMState::new()),
        }
    }

    fn resolve_eqfilter(&mut self, operator: TernaryOperator) -> Result<NMState, NmpolicyError> {
        self.resolve_ternary_operator(operator, filter::visit_state)
    }

    fn resolve_ternary_operator(
        &mut self,
        operator: TernaryOperator,
        resolver_fn: fn(NMState, Path, serde_json::Value) -> Result<NMState, NmpolicyError>,
    ) -> Result<NMState, NmpolicyError> {
        let operator_node = self.current_node.clone();
        let (lhs, ms, rhs) = operator;

        self.current_node = Some(lhs);
        let input_source = self.resolve_input_source()?;

        self.current_node = Some(ms.clone());
        let path = Path::compose_from_node(ms.kind)?;

        self.current_node = Some(rhs.clone());
        let value: serde_json::Value = match rhs.kind {
            NodeKind::Str(string) => serde_json::Value::String(string),
            _ => {
                return Err(NmpolicyError::new(ErrorKind::ResolveErrorNotSupportedValue));
            }
        };
        self.current_node = operator_node;
        resolver_fn(input_source, path, value)
    }

    fn resolve_input_source(&mut self) -> Result<NMState, NmpolicyError> {
        let current_node = self.current_node.clone().unwrap();
        match current_node.kind {
            n if n == current_state_identity() => Ok(self.current_state.clone().unwrap()),
            _ => match Path::compose_from_node(current_node.kind) {
                Ok(path) => match path.capture_entry_name {
                    Some(capture_entry_name) => {
                        self.resolve_capture_entry_by_name(&capture_entry_name)
                    }
                    None => {
                        return Err(NmpolicyError::new(
                            ErrorKind::ResolveErrorInvalidPathInputSource(format!(
                                "{:?}",
                                self.current_node
                            )),
                        ));
                    }
                },
                Err(_) => {
                    return Err(NmpolicyError::new(
                        ErrorKind::ResolveErrorInvalidInputSource(format!(
                            "{:?}",
                            self.current_node
                        )),
                    ));
                }
            },
        }
    }
}
