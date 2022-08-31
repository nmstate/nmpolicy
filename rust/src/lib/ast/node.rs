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

pub(crate) type TernaryOperator = (Box<Node>, Box<Node>, Box<Node>);

#[derive(Default, Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub(crate) struct Node {
    pub pos: usize,
    #[serde(flatten)]
    pub kind: NodeKind,
}

#[derive(Clone, Eq, PartialEq, Debug, Serialize, Deserialize)]
pub(crate) enum NodeKind {
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

impl Default for NodeKind {
    fn default() -> Self {
        NodeKind::Identity("nil".to_string())
    }
}

pub(crate) fn string(pos: usize, literal: String) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Str(literal),
    })
}

pub(crate) fn identity(pos: usize, literal: String) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Identity(literal),
    })
}

pub(crate) fn number(pos: usize, literal: i32) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Number(literal),
    })
}

pub(crate) fn eqfilter(
    pos: usize,
    value1: Box<Node>,
    value2: Box<Node>,
    value3: Box<Node>,
) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::EqFilter(value1, value2, value3),
    })
}

pub(crate) fn replace(
    pos: usize,
    value1: Box<Node>,
    value2: Box<Node>,
    value3: Box<Node>,
) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Replace(value1, value2, value3),
    })
}

pub(crate) fn path(pos: usize, nodes: Vec<Node>) -> Box<Node> {
    Box::new(Node {
        pos,
        kind: NodeKind::Path(nodes),
    })
}

pub(crate) fn current_state_identity() -> NodeKind {
    NodeKind::Identity(String::from("currentState"))
}

pub(crate) fn current_state(pos: usize) -> Box<Node> {
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
