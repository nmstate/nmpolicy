use std::collections::HashMap;

use crate::{
    error::NmpolicyError, expand::expander::StateExpander, resolve::resolver::CapturePathResolver,
    types::NMState,
};

use serde_json::Value;

struct CapturePathResolverStub {
    captured_states: HashMap<&'static str, Value>,
    force_failure: bool,
}

impl CapturePathResolver for CapturePathResolverStub {
    fn resolve_capture_entry_path(&self, capture_path: String) -> Result<Value, NmpolicyError> {
        if self.force_failure {
            return Err(NmpolicyError::from(
                "unit test forced resolver error".to_string(),
            ));
        }
        match self.captured_states.get(capture_path.as_str()) {
            Some(captured) => Ok(captured.clone()),
            None => Err(NmpolicyError::from("not found".to_string())),
        }
    }
}

#[test]
fn captures_with_map_values() {
    let desired_state: NMState = serde_yaml::from_str(
        r#"
interfaces:
- name: br1
  description: Linux bridge with base interface as a port
  type: linux-bridge
  state: up
  ipv4: "{{ capture.base-iface.interfaces.0.ipv4 }}"
  bridge:
    options:
      stp:
        enabled: false
    port:
    - name: "{{ capture.base-iface.interfaces.0.name }}"
routes:
  config: "{{ capture.bridge-routes-takeover.running }}"
"#,
    )
    .unwrap();

    let expected_expanded_desired_state: NMState = serde_yaml::from_str(
        r#"
interfaces:
- bridge:
    options:
      stp:
        enabled: false
    port:
    - name: eth1
  description: Linux bridge with base interface as a port
  ipv4: 1.2.3.4
  name: br1
  state: up
  type: linux-bridge
routes:
  config:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
"#,
    )
    .unwrap();

    let routes: Value = serde_yaml::from_str(
        r#"
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
"#,
    )
    .unwrap();

    let captured_states = HashMap::<&str, Value>::from([
        (
            "capture.base-iface.interfaces.0.ipv4",
            Value::from("1.2.3.4"),
        ),
        ("capture.base-iface.interfaces.0.name", Value::from("eth1")),
        ("capture.bridge-routes-takeover.running", routes),
    ]);

    let capture_path_resolver = CapturePathResolverStub {
        captured_states,
        force_failure: false,
    };

    let expander = StateExpander::new(Box::new(capture_path_resolver));
    match expander.expand(desired_state) {
        Ok(obtained_expanded_desired_state) => assert_eq!(
            expected_expanded_desired_state,
            obtained_expanded_desired_state
        ),
        Err(e) => panic!("{}", e),
    }
}

#[test]
fn capture_is_top_level() {
    let desired_state: NMState = serde_yaml::from_str(
        r#"
interfaces: "{{ capture.base-iface }}"
"#,
    )
    .unwrap();

    let interfaces: Value = serde_yaml::from_str(
        r#"
- name: eth1
  type: ethernet
  state: up
  ipv4:
    address:
    - ip: 10.244.0.1
      prefix-length: 24
    - ip: 169.254.1.0
      prefix-length: 16
    dhcp: false
    enabled: true
"#,
    )
    .unwrap();

    let mut expected_expanded_desired_state = NMState::new();
    expected_expanded_desired_state.insert("interfaces".to_string(), interfaces.clone());

    let captured_states = HashMap::<&str, Value>::from([("capture.base-iface", interfaces)]);

    let capture_path_resolver = CapturePathResolverStub {
        captured_states,
        force_failure: false,
    };

    let expander = StateExpander::new(Box::new(capture_path_resolver));
    match expander.expand(desired_state) {
        Ok(obtained_expanded_desired_state) => assert_eq!(
            expected_expanded_desired_state,
            obtained_expanded_desired_state
        ),
        Err(e) => panic!("{}", e),
    }
}

#[test]
fn resolve_capture_fails() {
    let desired_state: NMState = serde_yaml::from_str(
        r#"
interfaces: "{{ capture.enabled-iface }}"
"#,
    )
    .unwrap();

    let capture_path_resolver = CapturePathResolverStub {
        captured_states: HashMap::<&str, Value>::new(),
        force_failure: true,
    };

    let expander = StateExpander::new(Box::new(capture_path_resolver));

    let expand_result = expander.expand(desired_state);
    assert_eq!(true, expand_result.is_err());
}
