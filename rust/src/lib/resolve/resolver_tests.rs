use crate::{
    lex::tokens::Tokens,
    parse::parser::Parser,
    resolve::resolver::Resolver,
    types::{Capture, CaptureEntry, CapturedStates, NMState},
};
use std::collections::HashMap;
const CURRENT_STATE_YAML: &str = r#"
routes:
  running:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
  - destination: 2.2.2.0/24
    next-hop-address: 192.168.200.1
    next-hop-interface: eth2
    table-id: 254
  config:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
  - destination: 1.1.1.0/24
    next-hop-address: 192.168.100.1
    next-hop-interface: eth1
    table-id: 254
interfaces:
  - name: eth1
    description: "1st ethernet interface"
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
  - name: eth2
    type: ethernet
    state: down
    ipv4:
      address:
      - ip: 1.2.3.4
        prefix-length: 24
      dhcp: false
      enabled: false
"#;

struct Test<'a> {
    capture: HashMap<&'a str, &'a str>,
    cache: &'a str,
    captured: &'a str,
    error: &'a str,
}

macro_rules! resolve_capture_tests{
		($($name:ident: $value:expr,)*) => {
		$(
			#[test]
			fn $name() {
                let current_state: NMState = serde_yaml::from_str(CURRENT_STATE_YAML).unwrap();
                let test  = $value;
                let mut capture = Capture::new();
                for (k, expression) in test.capture.iter() {
                    let tokens: &mut Tokens = &mut Tokens::new(expression);
                    let mut parser = Parser::new(expression.to_string(), tokens);
                    let ast = parser.parse().unwrap().unwrap();
                    capture.insert(k.to_string(), CaptureEntry{expression: expression.to_string(), ast: ast});
                }
                let cache: Option<CapturedStates> = if !test.cache.is_empty(){
                    Some(serde_yaml::from_str(test.cache).unwrap())
                }else{
                    None
                };

                let expected_captured_states: Option<CapturedStates> = if !test.captured.is_empty(){
                    Some(serde_yaml::from_str(test.captured).unwrap())
                }else{
                    None
                };
                let mut resolver = Resolver::new(capture);
                match resolver.resolve(current_state, cache){
                    Ok(obtained_captured_states) => assert_eq!(expected_captured_states, Some(obtained_captured_states)),
                    Err(e) => assert_eq!(test.error, e.to_string()),
                };
			}
		)*
		}
	}

resolve_capture_tests! {
    filter_map_list_on_second_path_identity: Test{
        capture: HashMap::from([
            ("default-gw", "routes.running.destination == '0.0.0.0/0'")
        ]),
        cache: "",
        captured: r#"
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
"#,
        error: "",
    },
    filter_map_list_on_first_path_identity: Test{
        capture: HashMap::from([
            ("up-interfaces", "interfaces.state=='down'")
        ]),
        cache: "",
        captured: r#"
up-interfaces:
  state: 
    interfaces:
      - name: eth2
        type: ethernet
        state: down
        ipv4:
          address:
          - ip: 1.2.3.4
            prefix-length: 24
          dhcp: false
          enabled: false
"#,
        error: "",
    },
    filter_list: Test{
        capture: HashMap::from([
            ("specific-ipv4", "interfaces.ipv4.address.ip=='10.244.0.1'")
        ]),
        cache: "",
        captured: r#"
specific-ipv4:
  state: 
    interfaces:
    - name: eth1
      description: "1st ethernet interface"
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
        error: "",
    },

    filter_capture_ref: Test{
        capture: HashMap::from([
            ("base-iface-route", "routes.running.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface")
        ]),
        cache: r#"
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
"#,
        captured: r#"
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
base-iface-routes:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
      - destination: 1.1.1.0/24
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
"#,
        error: "",
    },
}
