use crate::{
    error::{evaluation_error, ErrorKind, NmpolicyError},
    resolve::{
        path::{Path, Step},
        visitor,
    },
    types::NMState,
};

use serde_json::{Map, Value};

struct ReplaceVisitor {
    replace_value: Value,
}

pub(crate) fn visit_state(
    input_state: NMState,
    path: Path,
    replace_value: Value,
) -> Result<NMState, NmpolicyError> {
    match visitor::visit_state(
        path,
        Value::Object(input_state),
        &mut ReplaceVisitor { replace_value },
    ) {
        Ok(replaced) => match replaced {
            Value::Object(map) => Ok(map),
            Value::Null => Ok(NMState::new()),
            _ => Err(evaluation_error(
                "failed converting result to a map".to_string(),
            )),
        },
        Err(e) => Err(e.ctx("failed applying operation on the path".to_string())),
    }
}

impl visitor::StateVisitor for ReplaceVisitor {
    fn visit_last_map(
        &mut self,
        path: Path,
        mut map_to_visit: Map<String, Value>,
    ) -> Result<Value, NmpolicyError> {
        match path.current_step() {
            Step::Identity(_, step) => {
                map_to_visit.insert(step.clone(), self.replace_value.clone());
                Ok(Value::from(map_to_visit))
            }
            _ => Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
        }
    }

    fn visit_last_slice(
        &mut self,
        path: Path,
        slice_to_visit: Vec<Value>,
    ) -> Result<Value, NmpolicyError> {
        match path.current_step() {
            Step::Identity(_, _) => self.visit_slice(path, slice_to_visit),
            _ => Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
        }
    }

    fn visit_map(
        &mut self,
        mut path: Path,
        mut map_to_visit: Map<String, Value>,
    ) -> Result<Value, NmpolicyError> {
        match path.clone().current_step() {
            Step::Identity(_, step) => match map_to_visit.get(step) {
                Some(value_to_visit) => {
                    path.next_step();
                    let visit_result = visitor::visit_state(path, value_to_visit.clone(), self)?;
                    map_to_visit.insert(step.clone(), visit_result);
                    Ok(Value::from(map_to_visit))
                }
                None => Ok(Value::Null),
            },
            _ => Err(NmpolicyError::new(ErrorKind::NotSupportedError)),
        }
    }

    fn visit_slice(
        &mut self,
        path: Path,
        slice_to_visit: Vec<Value>,
    ) -> Result<Value, NmpolicyError> {
        match path.current_step() {
            Step::Identity(_, _) => {
                let result: Result<Vec<Value>, NmpolicyError> = slice_to_visit
                    .into_iter()
                    .map(|v| visitor::visit_state(path.clone(), v, self))
                    .collect();
                Ok(Value::from(result?))
            }
            _ => Err(NmpolicyError::new(ErrorKind::NotSupportedError)),
        }
    }
}
