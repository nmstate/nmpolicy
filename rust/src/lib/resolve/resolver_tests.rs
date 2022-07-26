use crate::{
    capture::{Capture, CaptureEntry},
    lex::tokens::Tokens,
    parse::parser::Parser,
    resolve::resolver::Resolver,
    types::{CapturedStates, NMState},
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

struct Test {
    capture: HashMap<&'static str, &'static str>,
    cache: &'static str,
    captured: &'static str,
    error: &'static str,
    current: &'static str,
}

fn test_with() -> Test {
    Test {
        capture: HashMap::from([]),
        cache: "",
        captured: "",
        error: "",
        current: CURRENT_STATE_YAML,
    }
}

impl Test {
    fn current(mut self, current: &'static str) -> Self {
        self.current = current;
        self
    }
    fn capture(mut self, capture: HashMap<&'static str, &'static str>) -> Self {
        self.capture = capture;
        self
    }
    fn cache(mut self, cache: &'static str) -> Self {
        self.cache = cache;
        self
    }
    fn captured(mut self, captured: &'static str) -> Self {
        self.captured = captured;
        self
    }
    fn error(mut self, error: &'static str) -> Self {
        self.error = error;
        self
    }
}

macro_rules! resolve_capture_tests{
		($($name:ident: $value:expr,)*) => {
		$(
			#[test]
			fn $name() {
                let test  = $value;
                let current_state: NMState = serde_yaml::from_str(test.current).unwrap();
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
    filter_map_list_on_second_path_identity: test_with()
        .capture(HashMap::from([
            ("default-gw", "routes.running.destination == '0.0.0.0/0'")
        ]))
        .captured(r#"
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
"#),
    filter_map_list_on_first_path_identity: test_with()
        .capture(HashMap::from([
            ("up-interfaces", "interfaces.state=='down'")
        ]))
        .cache("")
        .captured(r#"
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
"#)
        .error("")
    ,
    filter_list: test_with()
        .capture(HashMap::from([
            ("specific-ipv4", "interfaces.ipv4.address.ip=='10.244.0.1'")
        ]))
        .cache("")
        .captured(r#"
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
"#)
        .error("")
    ,

    filter_capture_ref: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface")
        ]))
        .cache(r#"
default-gw:
  state: 
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
"#)
        .captured(r#"
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
"#)
        .error("")
    ,
    filter_capture_ref_without_captured_state: test_with()
        .capture(HashMap::from([
            ("default-gw", "routes.running.destination=='0.0.0.0/0'"),
            ("base-iface-routes", "routes.running.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface"),
        ]))
        .cache("")
        .captured(r#"
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
"#)
        .error("")
    ,
    filter_capture_ref_not_found: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface==capture.default-gw.routes"),
        ]))
        .cache("")
        .captured("")
        .error(r#"resolve error: eqfilter error: capture entry 'default-gw' not found
| routes.running.next-hop-interface==capture.default-gw.routes
| ...................................^"#)
    ,
    filter_bad_capture_ref: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface==capture"),
        ]))
        .cache("")
        .captured("")
        .error(r#"resolve error: eqfilter error: path capture ref is missing capture entry name
| routes.running.next-hop-interface==capture
| ...................................^"#)
    ,
    filter_capture_ref_invalid_state_for_path_map: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface==capture.default-gw.routes.running.badfield.next-hop-interface"),
        ]))
        .cache(r#"
default-gw:
  state:
    routes:
       running:
       - destination: 0.0.0.0/0
         next-hop-address: 192.168.100.1
         next-hop-interface: eth1
         table-id: 254
"#)
        .captured("")
        .error(r#"resolve error: eqfilter error: failed walking path: invalid path: unexpected non numeric step for slice state '[{"destination":"0.0.0.0/0","next-hop-address":"192.168.100.1","next-hop-interface":"eth1","table-id":254}]'
| routes.running.next-hop-interface==capture.default-gw.routes.running.badfield.next-hop-interface
| .....................................................................^"#)
    ,
    filter_capture_ref_invalid_state_for_path_slice: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface==capture.default-gw.routes.1.0.next-hop-interface"),
        ]))
        .cache(r#"
default-gw:
  state: 
    routes:
      running:
       - destination: 0.0.0.0/0
         next-hop-address: 192.168.100.1
         next-hop-interface: eth1
         table-id: 254
"#)
        .captured("")
        .error(r#"resolve error: eqfilter error: failed walking path: invalid path: unexpected non identity step for map state '{"running":[{"destination":"0.0.0.0/0","next-hop-address":"192.168.100.1","next-hop-interface":"eth1","table-id":254}]}'
| routes.running.next-hop-interface==capture.default-gw.routes.1.0.next-hop-interface
| .............................................................^"#)
    ,
    filter_capture_ref_path_not_found_map: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface==capture.default-gw.routes.badfield"),
        ]))
        .cache(r#"
default-gw:
  state: 
    routes:
      running:
       - destination: 0.0.0.0/0
         next-hop-address: 192.168.100.1
         next-hop-interface: eth1
         table-id: 254
"#)
        .captured("")
        .error(r#"resolve error: eqfilter error: failed walking path: invalid path: step not found at map state '{"running":[{"destination":"0.0.0.0/0","next-hop-address":"192.168.100.1","next-hop-interface":"eth1","table-id":254}]}'
| routes.running.next-hop-interface==capture.default-gw.routes.badfield
| .............................................................^"#)
    ,
    filter_capture_ref_path_not_found_slice: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface==capture.default-gw.routes.running.6"),
        ]))
        .cache(r#"
default-gw:
  state: 
    routes:
      running:
       - destination: 0.0.0.0/0
         next-hop-address: 192.168.100.1
         next-hop-interface: eth1
         table-id: 254
"#)
        .captured("")
        .error(r#"resolve error: eqfilter error: failed walking path: invalid path: step not found at slice state '[{"destination":"0.0.0.0/0","next-hop-address":"192.168.100.1","next-hop-interface":"eth1","table-id":254}]'
| routes.running.next-hop-interface==capture.default-gw.routes.running.6
| .....................................................................^"#)
    ,
    filter_different_type_on_path: test_with()
        .capture(HashMap::from([
            ("invalid-path-type", "interfaces.ipv4.address=='10.244.0.1'"),
        ]))
        .cache("")
        .captured("")
        .error(r#"resolve error: eqfilter error: failed applying operation on the path: invalid path: type missmatch: the value in the path doesn't match the value to filter. [{"ip":"10.244.0.1","prefix-length":24},{"ip":"169.254.1.0","prefix-length":16}] != "10.244.0.1"
| interfaces.ipv4.address=='10.244.0.1'
| ................^"#)
    ,
    filter_optional_field: test_with()
        .capture(HashMap::from([
            ("description-eth1", "interfaces.description=='1st ethernet interface'"),
        ]))
        .captured(r#"
description-eth1: 
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
"#)
    .cache("")
    .error("")
    ,
    filter_non_capture_ref_path_at_third_arg: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface==routes.running"),
        ]))
        .error(r#"resolve error: eqfilter error: not supported filtered value path. Only paths with a capture entry reference are supported
| routes.running.next-hop-interface==routes.running
| ...................................^"#)
        .cache("")
        .captured("")
    ,
    filter_by_path: test_with()
        .capture(HashMap::from([
            ("running-routes", "routes.running")
        ]))
        .captured(r#"
running-routes:
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
      - destination: 2.2.2.0/24
        next-hop-address: 192.168.200.1
        next-hop-interface: eth2
        table-id: 254
"#)
    .cache("")
    .error("")
    ,
    filter_with_invalid_input_source: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "invalidInputSource | routes.running.next-hop-interface=='eth1'")
        ]))
        .error(r#"resolve error: eqfilter error: invalid path input source (Path=[Identity=invalidInputSource]), only capture reference is supported
| invalidInputSource | routes.running.next-hop-interface=='eth1'
| ^"#)
        .cache("")
        .captured("")
    ,
    filter_with_invalid_type_in_source: test_with()
        .current(r#"
routes:
   running:
"#)
        .capture(HashMap::from([
            ("base-iface-routes", "routes.running.next-hop-interface=='eth1'")
        ]))
        .error(r#"resolve error: eqfilter error: failed applying operation on the path: invalid path: invalid type Null for identity step 'identity(next-hop-interface)'
| routes.running.next-hop-interface=='eth1'
| ...............^"#)
        .cache("")
        .captured("")
    ,
    filter_bad_path: test_with()
        .capture(HashMap::from([
            ("base-iface-routes", "routes.badfield.next-hop-interface=='eth1'")
        ]))
        .captured(r#"
base-iface-routes:
  state: {}
        "#)
        .cache("")
        .error("")
    ,
    replace_current_state: test_with()
        .capture(HashMap::from([
            ("bridge-routes", "routes.running.next-hop-interface := 'br1'")
        ]))
        .captured(r#"
bridge-routes:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: br1
        table-id: 254
      - destination: 1.1.1.0/24
        next-hop-address: 192.168.100.1
        next-hop-interface: br1
        table-id: 254
      - destination: 2.2.2.0/24
        next-hop-address: 192.168.200.1
        next-hop-interface: br1
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
"#),
    replace_captured_state: test_with()
    .cache(r#"
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
    "#)
    .capture(HashMap::from([
        ("bridge-routes", "capture.default-gw | routes.running.next-hop-interface := 'br1'")
    ]))
    .captured(r#"
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth1
        table-id: 254
bridge-routes:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: br1
        table-id: 254
    "#),
    replace_with_capture_ref: test_with()
    .cache(r#"
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: br1
        table-id: 254

br1-bridge:
  state:
    interfaces:
    - name: br1
      type: linux-bridge
      bridge:
        port:
        - name: eth3
"#)
    .capture(HashMap::from([
        ("default-gw-br1-first-port", "capture.default-gw | routes.running.next-hop-interface := capture.br1-bridge.interfaces.0.bridge.port.0.name")
    ]))
    .captured(r#"
default-gw:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: br1
        table-id: 254
br1-bridge:
  state:
    interfaces:
    - name: br1
      type: linux-bridge
      bridge:
        port:
        - name: eth3
default-gw-br1-first-port:
  state:
    routes:
      running:
      - destination: 0.0.0.0/0
        next-hop-address: 192.168.100.1
        next-hop-interface: eth3
        table-id: 254

"#), 
    replace_optional_field: test_with()
    .cache(r#"
eth2-interface:
  state:
    interfaces:
    - name: eth2
      type: ethernet
      state: down
"#)
    .capture(HashMap::from([
        ("description-eth2", "capture.eth2-interface | interfaces.description := '2nd ethernet interface'")
    ]))
    .captured(r#"
eth2-interface:
  state:
    interfaces:
    - name: eth2
      type: ethernet
      state: down
description-eth2:
  state:
    interfaces:
    - name: eth2
      description: "2nd ethernet interface"
      type: ethernet
      state: down
"#),
}
