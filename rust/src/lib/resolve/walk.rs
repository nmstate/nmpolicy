use crate::{
    error::{evaluation_error, NmpolicyError},
    resolve::{
        path::{Path, Step},
        visitor,
    },
    types::NMState,
};

use serde_json::{Map, Value};

struct WalkVisitor {}

pub(crate) fn visit_state(input_state: NMState, path: Path) -> Result<Value, NmpolicyError> {
    match visitor::visit_state(path, Value::Object(input_state), &mut WalkVisitor {}) {
        Ok(visit_result) => Ok(visit_result),
        Err(e) => Err(e.ctx("failed walking path".to_string())),
    }
}

impl visitor::StateVisitor for WalkVisitor {
    fn visit_last_map(
        &mut self,
        path: Path,
        map_to_access: Map<String, Value>,
    ) -> Result<Value, NmpolicyError> {
        access_map_with_current_state(path, map_to_access)
    }

    fn visit_last_slice(
        &mut self,
        path: Path,
        slice_to_access: Vec<Value>,
    ) -> Result<Value, NmpolicyError> {
        access_slice_with_current_state(path, slice_to_access)
    }

    fn visit_map(
        &mut self,
        mut path: Path,
        map_to_access: Map<String, Value>,
    ) -> Result<Value, NmpolicyError> {
        let value_to_visit = access_map_with_current_state(path.clone(), map_to_access)?;
        path.next_step();
        visitor::visit_state(path, value_to_visit, self)
    }
    fn visit_slice(
        &mut self,
        mut path: Path,
        slice_to_visit: Vec<Value>,
    ) -> Result<Value, NmpolicyError> {
        let value_to_visit = access_slice_with_current_state(path.clone(), slice_to_visit)?;
        path.next_step();
        visitor::visit_state(path, value_to_visit, self)
    }
}

fn access_map_with_current_state(
    path: Path,
    map_to_access: Map<String, Value>,
) -> Result<Value, NmpolicyError> {
    match path.current_step() {
        Step::Identity(pos, step) => match map_to_access.get(step) {
            Some(value) => Ok(value.clone()),
            None => Err(evaluation_error(format!(
                "step not found at map state '{}'",
                serde_json::to_string(&map_to_access)?,
            ))
            .path(*pos)),
        },
        Step::Number(pos, _) => Err(evaluation_error(format!(
            "unexpected non identity step for map state '{}'",
            serde_json::to_string(&map_to_access)?,
        ))
        .path(*pos)),
    }
}

fn access_slice_with_current_state(
    path: Path,
    slice_to_access: Vec<Value>,
) -> Result<Value, NmpolicyError> {
    match path.current_step() {
        Step::Number(pos, step) => {
            if slice_to_access.len() > *step {
                Ok(slice_to_access[*step].clone())
            } else {
                Err(evaluation_error(format!(
                    "step not found at slice state '{}'",
                    serde_json::to_string(&slice_to_access)?,
                ))
                .path(*pos))
            }
        }
        Step::Identity(pos, _) => Err(evaluation_error(format!(
            "unexpected non numeric step for slice state '{}'",
            serde_json::to_string(&slice_to_access)?,
        ))
        .path(*pos)),
    }
}
