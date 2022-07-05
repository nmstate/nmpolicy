use crate::{
    error::{ErrorKind, NmpolicyError},
    resolve::path::{Path, Step},
};

use serde_json::{Map, Value};

pub(crate) trait StateVisitor {
    fn visit_last_map(
        &mut self,
        path: Path,
        state: Map<String, Value>,
    ) -> Result<Value, NmpolicyError>;
    fn visit_last_slice(&mut self, path: Path, state: Vec<Value>) -> Result<Value, NmpolicyError>;
    fn visit_map(&mut self, path: Path, state: Map<String, Value>) -> Result<Value, NmpolicyError>;
    fn visit_slice(&mut self, path: Path, state: Vec<Value>) -> Result<Value, NmpolicyError>;
}

pub(crate) fn visit_state(
    path: Path,
    input_state: Value,
    state_visitor: &mut dyn StateVisitor,
) -> Result<Value, NmpolicyError> {
    match input_state {
        Value::Object(original_map) => {
            if path.has_more_steps() {
                match path.current_step() {
                    Step::Identity(_) => state_visitor.visit_map(path, original_map),
                    _ => Err(NmpolicyError::new(ErrorKind::PathErrorUnexpectedMapStep(
                        original_map,
                    ))),
                }
            } else {
                state_visitor.visit_last_map(path, original_map)
            }
        }
        Value::Array(original_slice) => {
            if path.has_more_steps() {
                state_visitor.visit_slice(path, original_slice)
            } else {
                state_visitor.visit_last_slice(path, original_slice)
            }
        }
        _ => Err(NmpolicyError::new(ErrorKind::PathErrorStepInvalidType(
            input_state,
            format!("{}", path.current_step()),
        ))),
    }
}
