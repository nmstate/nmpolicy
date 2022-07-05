use crate::{
    ast::node::NodeKind,
    error::{ErrorKind, NmpolicyError},
};

#[derive(Clone, Eq, PartialEq, Debug)]
pub(crate) enum Step {
    Identity(String),
    Number(usize),
}

impl Step {
    fn is_identity(&self) -> bool {
        matches!(self, Step::Identity(_))
    }
    fn is_number(&self) -> bool {
        matches!(self, Step::Number(_))
    }
}

#[derive(Clone, Eq, PartialEq, Debug)]
pub(crate) struct Path {
    pub capture_entry_name: Option<String>,
    pub steps: Vec<Step>,
    current_step_idx: usize,
}

impl Path {
    pub(crate) fn compose_from_node(node_kind: NodeKind) -> Result<Path, NmpolicyError> {
        match node_kind {
            NodeKind::Path(node_steps) => match node_steps {
                node_steps if !node_steps.is_empty() => {
                    let mut steps = Vec::<Step>::new();
                    for node_step in node_steps {
                        match node_step.kind {
                            NodeKind::Identity(step) => steps.push(Step::Identity(step)),
                            NodeKind::Number(step) => steps.push(Step::Number(step as usize)),
                            _ => return Err(NmpolicyError::new(ErrorKind::NotImplementedError)),
                        }
                    }
                    let mut path = Path {
                        capture_entry_name: None,
                        steps,
                        current_step_idx: 0,
                    };
                    match &path.steps[0] {
                        Step::Identity(first_step) => {
                            if first_step == "capture" {
                                let capture_ref_size = 2;
                                if path.steps.len() < capture_ref_size {
                                    return Err(NmpolicyError::new(ErrorKind::NotImplementedError));
                                } else {
                                    match &path.steps[1] {
                                        Step::Identity(capture_entry_name) => {
                                            path.capture_entry_name =
                                                Some(capture_entry_name.clone());
                                            path.steps = path.steps[1..].to_vec();
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
            Step::Identity(identity) => write!(f, "identity({identity})"),
            Step::Number(number) => write!(f, "number({number})"),
        }
    }
}
