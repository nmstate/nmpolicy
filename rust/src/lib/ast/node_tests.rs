use crate::ast::node::{current_state, eqfilter, identity, number, path, replace, string, Node};

#[test]
fn eqfilter_display() {
    let obtained_string = format!(
        "{}",
        eqfilter(
            0,
            path(0, vec![*current_state(0)]),
            path(
                0,
                vec![
                    *identity(0, "routes"),
                    *identity(0, "running"),
                    *identity(0, "table-id")
                ]
            ),
            number(0, 254),
        )
    );
    assert_eq!(obtained_string, "EqFilter([Path=[Identity=currentState] Path=[Identity=routes Identity=running Identity=table-id] Number=254])");
}

#[test]
fn replace_display() {
    let obtained_string = format!(
        "{}",
        replace(
            0,
            current_state(0),
            path(
                0,
                vec![
                    *identity(0, "routes"),
                    *identity(0, "running"),
                    *identity(0, "next-hop-interface")
                ]
            ),
            string(0, "br1"),
        )
    );
    assert_eq!(obtained_string, "Replace([Identity=currentState Path=[Identity=routes Identity=running Identity=next-hop-interface] String=br1])");
}

#[test]
fn eqfilter_yaml() {
    let obtained_deserialized: Box<Node> = serde_yaml::from_str(
        r#"
pos: 1
eqfilter:
- pos: 2
  path:
  - pos: 3
    identity: currentState
- pos: 4
  path:
  - pos: 5
    identity: routes
  - pos: 6
    identity: running
  - pos: 7
    identity: table-id
- pos: 8
  number: 254
"#,
    )
    .unwrap();

    let expected_deserialized = eqfilter(
        1,
        path(2, vec![*current_state(3)]),
        path(
            4,
            vec![
                *identity(5, "routes"),
                *identity(6, "running"),
                *identity(7, "table-id"),
            ],
        ),
        number(8, 254),
    );

    assert_eq!(obtained_deserialized, expected_deserialized);
}

#[test]
fn replace_yaml() {
    let obtained_deserialized: Box<Node> = serde_yaml::from_str(
        r#"
pos: 1
replace:
- pos: 2
  identity: currentState
- pos: 3
  path:
  - pos: 4
    identity: routes
  - pos: 5
    identity: running
  - pos: 6
    identity: next-hop-interface
- pos: 7
  string: br1
"#,
    )
    .unwrap();

    let expected_deserialized = replace(
        1,
        current_state(2),
        path(
            3,
            vec![
                *identity(4, "routes"),
                *identity(5, "running"),
                *identity(6, "next-hop-interface"),
            ],
        ),
        string(7, "br1"),
    );

    assert_eq!(obtained_deserialized, expected_deserialized);
}
