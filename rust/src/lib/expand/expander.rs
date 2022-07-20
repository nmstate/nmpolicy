use crate::{
    error::{ErrorKind, NmpolicyError},
    resolve::resolver::CapturePathResolver,
    types::NMState,
};

use regex::Regex;

use serde_json::{Map, Value};

pub(crate) struct StateExpander {
    capture_path_resolver: Box<dyn CapturePathResolver>,
}

impl StateExpander {
    pub(crate) fn new(capture_path_resolver: Box<dyn CapturePathResolver>) -> StateExpander {
        StateExpander {
            capture_path_resolver,
        }
    }

    pub(crate) fn expand(&self, state: NMState) -> Result<NMState, NmpolicyError> {
        match self.expand_state(Value::from(state)) {
            Ok(value) => match value {
                Value::Object(map) => Ok(map),
                _ => Err(NmpolicyError::new(ErrorKind::Bug)),
            },
            Err(e) => Err(e),
        }
    }

    fn expand_state(&self, state: Value) -> Result<Value, NmpolicyError> {
        match state {
            Value::Null => Ok(state),
            Value::String(string) => self.expand_string(string),
            Value::Object(map) => self.expand_map(map),
            Value::Array(slice) => self.expand_slice(slice),
            _ => Ok(state),
        }
    }

    fn expand_slice(&self, slice: Vec<Value>) -> Result<Value, NmpolicyError> {
        slice.into_iter().map(|v| self.expand_state(v)).collect()
    }
    fn expand_map(&self, map: Map<String, Value>) -> Result<Value, NmpolicyError> {
        map.into_iter()
            .map(|(k, v)| match self.expand_state(v) {
                Ok(ev) => Ok((k, ev)),
                Err(e) => Err(e),
            })
            .collect()
    }
    fn expand_string(&self, string: String) -> Result<Value, NmpolicyError> {
        let re = Regex::new(r"^\{\{ (.*) \}\}$")?;
        match re.captures(&string) {
            Some(caps) => match caps.get(1) {
                Some(capture_path) => self
                    .capture_path_resolver
                    .resolve_capture_entry_path(capture_path.as_str().to_string()),
                None => Ok(Value::from(string)),
            },
            None => Ok(Value::from(string)),
        }
    }
}
