use crate::{
    error::{evaluation_error, NmpolicyError},
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
                    Step::Identity(_, _) => state_visitor.visit_map(path, original_map),
                    Step::Number(pos, _) => Err(evaluation_error(format!(
                        "unexpected non identity step for map state '{}'",
                        serde_json::to_string(&original_map)?,
                    ))
                    .path(*pos)),
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
        _ => Err(evaluation_error(format!(
            "invalid type {:?} for identity step '{}'",
            input_state,
            path.current_step(),
        ))
        .path(path.current_step().pos())),
    }
}
