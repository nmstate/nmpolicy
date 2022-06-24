/*
 * Copyright 2021 NMPolicy Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

use serde::{Deserialize, Serialize};
use std::fmt;

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub struct Node {
    pub pos: usize,
    #[serde(flatten)]
    pub kind: NodeKind,
}

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub enum NodeKind {
    #[serde(rename = "string")]
    Str(String),
    #[serde(rename = "identity")]
    Identity(String),
    #[serde(rename = "number")]
    Number(i32),
    #[serde(rename = "eqfilter")]
    EqFilter(Box<Node>, Box<Node>, Box<Node>),
    #[serde(rename = "replace")]
    Replace(Box<Node>, Box<Node>, Box<Node>),
    #[serde(rename = "path")]
    Path(Vec<Node>),
}

pub fn string(pos: usize, value: &str) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Str(value.to_string()),
    })
}

pub fn identity(pos: usize, value: &str) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Identity(value.to_string()),
    })
}

pub fn number(pos: usize, value: i32) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Number(value),
    })
}

pub fn eqfilter(pos: usize, value1: Box<Node>, value2: Box<Node>, value3: Box<Node>) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::EqFilter(value1, value2, value3),
    })
}

pub fn replace(pos: usize, value1: Box<Node>, value2: Box<Node>, value3: Box<Node>) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Replace(value1, value2, value3),
    })
}

pub fn path(pos: usize, nodes: Vec<Node>) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Path(nodes),
    })
}

pub fn current_state_identity() -> NodeKind {
    NodeKind::Identity(String::from("currentState"))
}

pub fn current_state(pos: usize) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: current_state_identity(),
    })
}

impl fmt::Display for Node {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match &self.kind {
            NodeKind::Str(value) => write!(f, "String={}", value),
            NodeKind::Identity(value) => write!(f, "Identity={}", value),
            NodeKind::Number(value) => write!(f, "Number={}", value),
            NodeKind::EqFilter(value1, value2, value3) => {
                write!(f, "EqFilter([{} {} {}])", value1, value2, value3)
            }
            NodeKind::Replace(value1, value2, value3) => {
                write!(f, "Replace([{} {} {}])", value1, value2, value3)
            }
            NodeKind::Path(value) => write!(
                f,
                "Path=[{}]",
                value
                    .iter()
                    .map(|item| format!("{item}"))
                    .collect::<Vec<String>>()
                    .join(" ")
            ),
        }
    }
}
