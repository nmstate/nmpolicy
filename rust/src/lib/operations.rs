use crate::{
    capture,
    capture::CaptureEntryResolver,
    error::NmpolicyError,
    expand::expander::StateExpander,
    types::{CapturedStates, GeneratedState, NMState, PolicySpec},
};

pub fn generate_state(
    policy_spec: PolicySpec,
    current_state: NMState,
    cache: Option<CapturedStates>,
) -> Result<GeneratedState, NmpolicyError> {
    let captured_states = capture::resolve_entries(policy_spec.capture, current_state, cache)?;
    let capture_entry_resolver = CaptureEntryResolver::new(captured_states.clone());
    let state_expander = StateExpander::new(Box::new(capture_entry_resolver));
    let expanded_desird_state = state_expander.expand(policy_spec.desired_state)?;
    Ok(GeneratedState {
        cache: captured_states,
        desired_state: expanded_desird_state,
    })
}
