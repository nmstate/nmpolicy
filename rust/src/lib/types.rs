use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::ast::node::Node;

pub(crate) type NMState = serde_json::Map<String, serde_json::Value>;

pub(crate) type Capture = HashMap<String, CaptureEntry>;

pub(crate) type CapturedStates = HashMap<String, CapturedState>;

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub(crate) struct CaptureEntry {
    pub expression: String,
    pub ast: Box<Node>,
}

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub(crate) struct MetaInfo {
    pub version: Option<String>,
    pub time_stamp: Option<DateTime<Utc>>,
}

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub(crate) struct CapturedState {
    pub state: NMState,
    pub meta_info: Option<MetaInfo>,
}
