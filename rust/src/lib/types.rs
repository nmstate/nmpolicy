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

#[derive(Default, Clone, Eq, PartialEq, Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PolicySpec {
    #[serde(default)]
    pub capture: Capture,
    pub desired_state: NMState,
}

#[derive(Default, Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub struct GeneratedState {
    pub cache: CapturedStates,
    pub desired_state: NMState,
}

impl GeneratedState {
    pub fn new() -> Self {
        GeneratedState {
            cache: CapturedStates::new(),
            desired_state: NMState::new(),
        }
    }
}

impl PolicySpec {
    pub fn new() -> Self {
        PolicySpec {
            capture: Capture::new(),
            desired_state: NMState::new(),
        }
    }
}
