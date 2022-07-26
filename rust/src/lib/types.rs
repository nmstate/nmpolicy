use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::capture::Capture;

pub type NMState = serde_json::Map<String, serde_json::Value>;

pub type CapturedStates = HashMap<String, CapturedState>;

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub struct MetaInfo {
    pub version: Option<String>,
    pub time_stamp: Option<DateTime<Utc>>,
}

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub struct CapturedState {
    pub state: NMState,
    pub meta_info: Option<MetaInfo>,
}

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub struct PolicySpec {
    pub capture: Capture,
    pub desired_state: NMState,
}

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub struct GeneratedState {
    pub cache: CapturedStates,
    pub desired_state: NMState,
}
