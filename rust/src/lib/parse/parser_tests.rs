use crate::{
    ast::node::Node,
    error::NmpolicyError,
    lex::tokens::{
        Token, Token::Dot, Token::EqFilter, Token::Identity, Token::Number, Token::Pipe,
        Token::Replace, Token::Str,
    },
    parse::parser::Parser,
};

macro_rules! parse_tokens_tests{
		($($name:ident: $value:expr,)*) => {
		$(
			#[test]
			fn $name() {
                let (from_tokens, expected_ast_yaml) = $value;
                let expected_ast: Option<Box<Node>> = if expected_ast_yaml.is_empty() {
                    None
                } else {
                    Some(serde_yaml::from_str(expected_ast_yaml).unwrap())
                };
            	println!("{:?}", from_tokens);
                let tokens = from_tokens.tokens.clone();
                let expression = from_tokens.expression.clone();
                let tokens_iterator = &mut tokens.into_iter();
                let mut parser = Parser::new(expression, tokens_iterator);
                let obtained_result = parser.parse();
                let expected_result = Ok(expected_ast);
            	assert_eq!(expected_result, obtained_result);
			}
		)*
		}
	}

macro_rules! parse_errors_tests{
		($($name:ident: $value:expr,)*) => {
		$(
			#[test]
			fn $name() {
                let (from_tokens, expected_error) = $value;
            	println!("{:?}", from_tokens);
                let tokens = from_tokens.tokens.clone();
                let expression = from_tokens.expression.clone();
                let tokens_iterator = &mut tokens.into_iter();
                let mut parser = Parser::new(expression, tokens_iterator);
                let obtained_result = parser.parse();
                assert_eq!(true, obtained_result.is_err());
                assert_eq!(expected_error, obtained_result.unwrap_err().to_string());
			}
		)*
		}
	}

#[derive(Clone, Eq, PartialEq, Debug)]
struct FromTokens {
    expression: String,
    tokens: Vec<Result<Token, NmpolicyError>>,
}

impl FromTokens {
    fn new() -> Self {
        Self {
            expression: "".to_string(),
            tokens: vec![],
        }
    }
    fn identity(&mut self, literal: &str) -> &mut Self {
        self.tokens
            .push(Ok(Identity(self.expression.len(), literal.to_string())));
        self.expression.push_str(literal);
        self
    }
    fn string(&mut self, literal: &str) -> &mut Self {
        self.tokens
            .push(Ok(Str(self.expression.len(), literal.to_string())));
        self.expression.push_str(literal);
        self
    }
    fn number(&mut self, literal: i32) -> &mut Self {
        self.tokens.push(Ok(Number(self.expression.len(), literal)));
        self.expression.push_str(format!("{literal}").as_str());
        self
    }
    fn dot(&mut self) -> &mut Self {
        let dot = Dot(self.expression.len());
        self.tokens.push(Ok(dot.clone()));
        self.expression = format!("{}{}", self.expression, dot.literal());
        self
    }

    fn eqfilter(&mut self) -> &mut Self {
        let eqfilter = EqFilter(self.expression.len());
        self.tokens.push(Ok(eqfilter.clone()));
        self.expression = format!("{}{}", self.expression, eqfilter.literal());
        self
    }

    fn replace(&mut self) -> &mut Self {
        let eqfilter = Replace(self.expression.len());
        self.tokens.push(Ok(eqfilter.clone()));
        self.expression = format!("{}{}", self.expression, eqfilter.literal());
        self
    }

    fn pipe(&mut self) -> &mut Self {
        let eqfilter = Pipe(self.expression.len());
        self.tokens.push(Ok(eqfilter.clone()));
        self.expression = format!("{}{}", self.expression, eqfilter.literal());
        self
    }

    fn eof(&self) -> Self {
        self.clone()
    }
}

fn from_tokens() -> FromTokens {
    FromTokens::new()
}

parse_tokens_tests! {
   parse_empty: (from_tokens().eof(), ""),
   parse_path: (from_tokens()
                    .identity("routes")
                    .dot()
                    .identity("running")
                    .dot()
                    .identity("destination")
                    .eof(), r#"
pos: 0
path:
- pos: 0
  identity: routes
- pos: 7
  identity: running
- pos: 15
  identity: destination
"#),
  parse_eqfilter_with_string:(from_tokens()
                        .identity("routes")
                        .dot()
                        .identity("running")
                        .dot()
                        .identity("destination")
                        .eqfilter()
                        .string("0.0.0.0/0")
                        .eof(), r#"
pos: 26
eqfilter: 
- pos: 0
  identity: currentState
- pos: 0
  path: 
  - pos: 0
    identity: routes
  - pos: 7
    identity: running
  - pos: 15
    identity: destination
- pos: 28 
  string: 0.0.0.0/0
"#),
parse_eqfilter_with_path:(from_tokens()
                    .identity("routes")
                    .dot()
                    .identity("running")
                    .dot()
                    .identity("next-hop-interface")
                    .eqfilter()
                    .identity("capture")
                    .dot()
                    .identity("default-gw")
                    .dot()
                    .identity("routes")
                    .dot()
                    .number(0)
                    .dot()
                    .identity("next-hop-interface")
                    .eof(),r#"
pos: 33
eqfilter:
- pos: 0
  identity: currentState
- pos: 0 
  path: 
  - pos: 0 
    identity: routes
  - pos: 7
    identity: running
  - pos: 15
    identity: next-hop-interface
- pos: 35
  path:
  - pos: 35 
    identity: capture
  - pos: 43
    identity: default-gw
  - pos: 54
    identity: routes
  - pos: 61
    number: 0
  - pos: 63
    identity: next-hop-interface
"#),
parse_replace_with_string:(from_tokens()
                .identity("routes")
                .dot()
                .identity("running")
                .dot()
                .identity("next-hop-interface")
                .replace()
                .string("br1")
                .eof(),r#"
pos: 33
replace:
- pos: 0
  identity: currentState
- pos: 0 
  path: 
  - pos: 0 
    identity: routes
  - pos: 7
    identity: running
  - pos: 15
    identity: next-hop-interface
- pos: 35
  string: br1

"#),
parse_replace_with_path:(from_tokens()
                .identity("routes")
                .dot()
                .identity("running")
                .dot()
                .identity("next-hop-interface")
                .replace()
                .identity("capture")
                .dot()
                .identity("primary-nic")
                .dot()
                .identity("interfaces")
                .dot()
                .number(0)
                .dot()
                .identity("name")
                .eof(),r#"
pos: 33
replace:
- pos: 0
  identity: currentState
- pos: 0 
  path: 
  - pos: 0 
    identity: routes
  - pos: 7
    identity: running
  - pos: 15
    identity: next-hop-interface
- pos: 35
  path:
  - pos: 35
    identity: capture
  - pos: 43
    identity: primary-nic
  - pos: 55
    identity: interfaces
  - pos: 66
    number: 0
  - pos: 68
    identity: name
"#),
parse_capture_pipe_replace: (from_tokens()
            .identity("capture")
            .dot()
            .identity("default-gw")
            .pipe()
            .identity("routes")
            .dot()
            .identity("running")
            .dot()
            .identity("next-hop-interface")
            .replace()
            .string("br1")
            .eof(),r#"
pos: 52
replace:
- pos: 0
  path:
  - pos: 0
    identity: capture
  - pos: 8
    identity: default-gw
- pos: 19 
  path: 
  - pos: 19
    identity: routes
  - pos: 26
    identity: running
  - pos: 34
    identity: next-hop-interface
- pos: 54
  string: br1
"#),
}

parse_errors_tests! {
    parse_basic_failures: (from_tokens().dot().eof(),
    r#"invalid expression: unexpected token `.`
| .
| ^"#),
    parse_path_failure_0: (from_tokens().identity("routes").dot().eof(),
    r#"invalid path: missing identity or number after dot
| routes.
| ......^"#),
    parse_path_failure_1: (from_tokens().identity("routes").identity("destination").eof(),
    r#"invalid path: missing dot
| routesdestination
| ......^"#),
    parse_path_failure_2: (from_tokens().identity("routes").dot().dot().identity("destination").eof(),
    r#"invalid path: missing identity or number after dot
| routes..destination
| .......^"#),
    parse_eqfilter_failure_0: (from_tokens().eqfilter().string("0.0.0.0/0").eof(),
    r#"invalid equality filter: missing left hand argument
| ==0.0.0.0/0
| ^"#),
    parse_eqfilter_failure_1: (from_tokens().string("foo").eqfilter().string("0.0.0.0/0").eof(),
    r#"invalid equality filter: left hand argument is not a path
| foo==0.0.0.0/0
| ...^"#),
    parse_eqfilter_failure_2: (from_tokens().identity("routes").dot().identity("running").dot().identity("destination").eqfilter().eof(),
    r#"invalid equality filter: missing right hand argument
| routes.running.destination==
| ..........................^"#),
    parse_eqfilter_failure_3: (from_tokens().identity("routes").dot().identity("running").dot().identity("destination").eqfilter().eqfilter().eof(),
    r#"invalid equality filter: right hand argument is not a string or identity
| routes.running.destination====
| ............................^"#),
    parse_pipe_failure_0: (from_tokens().pipe().identity("routes").dot().identity("running").dot().identity("next-hop-interface").replace().string("br1").eof(),
    r#"invalid pipe: missing pipe in expression
| |routes.running.next-hop-interface:=br1
| ^"#),
    parse_pipe_failure_1: (from_tokens().identity("capture").dot().identity("default-gw").pipe().eof(),
    r#"invalid pipe: missing pipe out expression
| capture.default-gw|
| ..................^"#),
    parse_pipe_failure_2: (from_tokens().string("foo").pipe().identity("routes").dot().identity("running").dot().identity("next-hop-interface").replace().string("br1").eof(),
    r#"invalid pipe: only paths can be piped in
| foo|routes.running.next-hop-interface:=br1
| ...^"#),
    parse_replace_failure_0: (from_tokens().replace().string("0.0.0.0/0").eof(),
    r#"invalid replace: missing left hand argument
| :=0.0.0.0/0
| ^"#),
    parse_replace_failure_1: (from_tokens().string("foo").replace().string("0.0.0.0/0").eof(),
    r#"invalid replace: left hand argument is not a path
| foo:=0.0.0.0/0
| ...^"#),
    parse_replace_failure_2: (from_tokens().identity("routes").dot().identity("running").dot().identity("destination").replace().eof(),
    r#"invalid replace: missing right hand argument
| routes.running.destination:=
| ..........................^"#),
    parse_replace_failure_3: (from_tokens().identity("routes").dot().identity("running").dot().identity("destination").replace().replace().eof(),
    r#"invalid replace: right hand argument is not a string or identity
| routes.running.destination:=:=
| ............................^"#),
}
