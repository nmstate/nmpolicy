use crate::{
    ast::node::{Node, NodeKind},
    error::{validation_error, ErrorKind, NmpolicyError},
};

#[derive(Clone, Eq, PartialEq, Debug)]
pub(crate) enum Step {
    Identity(usize, String),
    Number(usize, usize),
}

impl Step {
    pub(crate) fn pos(&self) -> usize {
        match self {
            Step::Identity(pos, _) => *pos,
            Step::Number(pos, _) => *pos,
        }
    }
}

#[derive(Clone, Eq, PartialEq, Debug)]
pub(crate) struct Path {
    pub capture_entry_name: Option<String>,
    pub steps: Vec<Step>,
    current_step_idx: usize,
}

impl Path {
    pub(crate) fn compose_from_node(node: Node) -> Result<Path, NmpolicyError> {
        match node.kind {
            NodeKind::Path(node_steps) => match node_steps {
                node_steps if !node_steps.is_empty() => {
                    let mut steps = Vec::<Step>::new();
                    for node_step in node_steps {
                        match node_step.kind {
                            NodeKind::Identity(step) => {
                                steps.push(Step::Identity(node_step.pos, step))
                            }
                            NodeKind::Number(step) => {
                                steps.push(Step::Number(node_step.pos, step as usize))
                            }
                            _ => return Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
                        }
                    }
                    let mut path = Path {
                        capture_entry_name: None,
                        steps,
                        current_step_idx: 0,
                    };
                    match &path.steps[0] {
                        Step::Identity(_, first_step) => {
                            if first_step == "capture" {
                                let capture_ref_size = 2;
                                if path.steps.len() < capture_ref_size {
                                    return Err(validation_error(
                                        "path capture ref is missing capture entry name"
                                            .to_string(),
                                    ));
                                } else {
                                    match &path.steps[1] {
                                        Step::Identity(_, capture_entry_name) => {
                                            path.capture_entry_name =
                                                Some(capture_entry_name.clone());
                                            path.steps = path.steps[2..].to_vec();
                                        }
                                        _ => {
                                            return Err(NmpolicyError::new(
                                                ErrorKind::NotImplementedError,
                                            ))
                                        }
                                    }
                                }
                            }
                        }
                        _ => return Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
                    }
                    Ok(path)
                }
                _ => Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
            },
            _ => Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
        }
    }

    pub(crate) fn has_more_steps(&self) -> bool {
        self.current_step_idx + 1 < self.steps.len()
    }

    pub(crate) fn current_step(&self) -> &Step {
        &self.steps[self.current_step_idx]
    }
    pub(crate) fn next_step(&mut self) {
        if self.has_more_steps() {
            self.current_step_idx += 1;
        }
    }
}

impl std::fmt::Display for Step {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Step::Identity(_, identity) => write!(f, "identity({identity})"),
            Step::Number(_, number) => write!(f, "number({number})"),
        }
    }
}
