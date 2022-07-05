use std::iter::FromIterator;

use crate::{
    error::{ErrorKind, NmpolicyError},
    resolve::{
        path::{Path, Step},
        visitor,
    },
    types::NMState,
};

use serde_json::{Map, Value};

struct FilterVisitor {
    merge_visit_result: bool,
    expected_value: Value,
}

pub(crate) fn visit_state(
    input_state: NMState,
    path: Path,
    expected_value: Value,
) -> Result<NMState, NmpolicyError> {
    match visitor::visit_state(
        path,
        Value::Object(input_state),
        &mut FilterVisitor {
            merge_visit_result: false,
            expected_value,
        },
    ) {
        Ok(filtered) => match filtered {
            Value::Object(map) => Ok(map),
            Value::Null => Ok(NMState::new()),
            _ => Err(NmpolicyError::new(
                ErrorKind::ResolveErrorFailedFilterResultConvertion(filtered),
            )),
        },
        Err(e) => Err(NmpolicyError::new(
            ErrorKind::ResolveErrorFailingFilteringPath(e.to_string()),
        )),
    }
}

impl visitor::StateVisitor for FilterVisitor {
    fn visit_last_map(
        &mut self,
        path: Path,
        map_to_filter: Map<String, Value>,
    ) -> Result<Value, NmpolicyError> {
        match path.current_step() {
            Step::Identity(step) => match map_to_filter.get(step) {
                Some(obtained_value) => match self.expected_value {
                    Value::Null => Ok(Value::from(Map::<String, Value>::from_iter([(
                        step.clone(),
                        obtained_value.clone(),
                    )]))),
                    _ => {
                        if !value_has_same_type(&self.expected_value, obtained_value) {
                            Err(NmpolicyError::new(ErrorKind::NotImplementedError))
                        } else if *obtained_value == self.expected_value {
                            Ok(Value::from(map_to_filter))
                        } else {
                            Ok(Value::Null)
                        }
                    }
                },
                None => Ok(Value::Null),
            },
            _ => Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
        }
    }
    fn visit_last_slice(
        &mut self,
        path: Path,
        slice_to_visit: Vec<Value>,
    ) -> Result<Value, NmpolicyError> {
        match path.current_step() {
            Step::Identity(_) => self.visit_slice(path, slice_to_visit),
            _ => Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
        }
    }
    fn visit_map(
        &mut self,
        mut path: Path,
        map_to_visit: Map<String, Value>,
    ) -> Result<Value, NmpolicyError> {
        match path.clone().current_step() {
            Step::Identity(step) => match map_to_visit.get(step) {
                Some(value_to_visit) => {
                    path.next_step();
                    match visitor::visit_state(path.clone(), value_to_visit.clone(), self) {
                        Ok(visit_result) => match visit_result {
                            Value::Null => Ok(Value::Null),
                            _ => {
                                let mut filtered_map = if self.merge_visit_result {
                                    map_to_visit.clone()
                                } else {
                                    Map::<String, Value>::new()
                                };
                                filtered_map.insert(step.clone(), visit_result);
                                Ok(Value::from(filtered_map))
                            }
                        },
                        Err(e) => Err(e),
                    }
                }
                None => Ok(Value::Null),
            },
            Step::Number(_) => Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
        }
    }
    fn visit_slice(
        &mut self,
        path: Path,
        slice_to_visit: Vec<Value>,
    ) -> Result<Value, NmpolicyError> {
        match path.current_step() {
            Step::Identity(_) => {
                let mut filtered_slice = Vec::<Value>::new();
                let mut has_visit_result = false;
                for value_to_visit in slice_to_visit {
                    let visit_result = visitor::visit_state(
                        path.clone(),
                        value_to_visit.clone(),
                        &mut FilterVisitor {
                            merge_visit_result: true,
                            expected_value: self.expected_value.clone(),
                        },
                    )?;
                    if !visit_result.is_null() {
                        has_visit_result = true;
                        filtered_slice.push(visit_result);
                    } else if self.merge_visit_result {
                        filtered_slice.push(value_to_visit);
                    }
                }
                if !has_visit_result {
                    Ok(Value::Null)
                } else {
                    Ok(Value::from(filtered_slice))
                }
            }
            Step::Number(_) => Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
        }
    }
}

fn value_has_same_type(lhs: &Value, rhs: &Value) -> bool {
    lhs.is_boolean() == rhs.is_boolean()
        && lhs.is_number() == rhs.is_number()
        && lhs.is_string() == rhs.is_string()
        && lhs.is_array() == rhs.is_array()
        && lhs.is_object() == rhs.is_object()
}
