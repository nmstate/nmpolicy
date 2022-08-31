use std::io::Read;

use nmpolicy::{
    error::NmpolicyError,
    types::{CapturedStates, NMState, PolicySpec},
};

fn main() -> Result<(), NmpolicyError> {
    let matches = clap::Command::new("nmpolicyctl")
        .subcommand(
            clap::command!("gen")
                .usage("gen policy.yaml")
                .about("Generates NMState by policy filename")
                .arg(
                    clap::arg!(<POLICY>)
                        .required(true)
                        .value_parser(clap::value_parser!(std::path::PathBuf)),
                )
                .arg(
                    clap::arg!(-s --"current-state" <PATH>)
                        .required(false)
                        .help(
                            "input file path to current NMState. If not specified, STDIN is used.",
                        )
                        .value_parser(clap::value_parser!(std::path::PathBuf)),
                )
                .arg(
                    clap::arg!(-i --"captured-states-input" <PATH>)
                        .required(false)
                        .help("input file path for already resolved captured states.")
                        .value_parser(clap::value_parser!(std::path::PathBuf)),
                )
                .arg(
                    clap::arg!(-o --"captured-states-output" <PATH>)
                        .required(false)
                        .help("output file path to the emitted captured states.")
                        .value_parser(clap::value_parser!(std::path::PathBuf)),
                ),
        )
        .get_matches();

    match matches.subcommand() {
        Some(("gen", gen_matches)) => {
            let policy_spec_path = gen_matches.get_one::<std::path::PathBuf>("POLICY");
            match generate_state(
                policy_spec_path.unwrap(),
                gen_matches.get_one::<std::path::PathBuf>("current-state"),
                gen_matches.get_one::<std::path::PathBuf>("captured-states-input"),
                gen_matches.get_one::<std::path::PathBuf>("captured-states-output"),
            ) {
                Ok(generated_state) => {
                    print!("{}", generated_state);
                    Ok(())
                }
                Err(e) => Err(e),
            }
        }
        _ => unreachable!(), // If all subcommands are defined above, anything else is unreachable
    }
}

fn generate_state(
    policy_spec_path: &std::path::Path,
    current_state_path: Option<&std::path::PathBuf>,
    captured_states_input_path: Option<&std::path::PathBuf>,
    captured_states_output_path: Option<&std::path::PathBuf>,
) -> Result<String, NmpolicyError> {
    let policy_spec = read_policy_spec(policy_spec_path)?;
    let current_state = read_current_state(current_state_path)?;
    let captured_states_input = read_captured_states_input(captured_states_input_path)?;
    if policy_spec == PolicySpec::new() {
        if let Some(o) = captured_states_output_path {
            std::fs::File::create(o)?;
        }
        return Ok(String::new());
    }
    let captured_states =
        nmpolicy::operations::generate_state(policy_spec, current_state, captured_states_input)?;
    if let Some(o) = captured_states_output_path {
        let captured_states_writer = std::fs::File::create(o)?;
        if !captured_states.cache.is_empty() {
            serde_yaml::to_writer(captured_states_writer, &captured_states.cache)?
        }
    }
    Ok(serde_yaml::to_string(&captured_states.desired_state)?)
}

fn read_policy_spec(policy_spec_path: &std::path::Path) -> Result<PolicySpec, NmpolicyError> {
    let policy_spec_string = std::fs::read_to_string(policy_spec_path)?;
    let policy_spec: PolicySpec = if !policy_spec_string.is_empty() && policy_spec_string != "\n" {
        serde_yaml::from_str(&policy_spec_string)?
    } else {
        PolicySpec::new()
    };
    Ok(policy_spec)
}

fn read_current_state(
    current_state_path: Option<&std::path::PathBuf>,
) -> Result<NMState, NmpolicyError> {
    let mut current_state_string = String::new();
    if let Some(path) = current_state_path {
        current_state_string = std::fs::read_to_string(path)?;
    } else {
        std::io::stdin().read_to_string(&mut current_state_string)?;
    };
    let current_state: NMState = if !current_state_string.is_empty() && current_state_string != "\n"
    {
        serde_yaml::from_str(&current_state_string)?
    } else {
        NMState::new()
    };
    Ok(current_state)
}

fn read_captured_states_input(
    captured_state_path_o: Option<&std::path::PathBuf>,
) -> Result<Option<CapturedStates>, NmpolicyError> {
    if let Some(captured_state_path) = captured_state_path_o {
        let captured_state_string = std::fs::read_to_string(captured_state_path)?;
        if !captured_state_string.is_empty() && captured_state_string != "\n" {
            let captured_states: CapturedStates = serde_yaml::from_str(&captured_state_string)?;
            return Ok(Some(captured_states));
        }
    }
    Ok(None)
}
