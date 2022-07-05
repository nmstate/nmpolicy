use crate::{
    ast::node::{current_state, eqfilter, identity, number, path, replace, string, Node, NodeKind},
    error::{ErrorKind, NmpolicyError},
    lex::tokens::Token,
};

type TernaryOperator = (Option<Box<Node>>, Option<Box<Node>>, Option<Box<Node>>);

pub(crate) struct Parser<'a> {
    expression: String,
    tokens: &'a mut dyn Iterator<Item = Result<Token, NmpolicyError>>,
    current_token: Option<Token>,
    root_node: Option<Box<Node>>,
    piped_in_node: Option<Box<Node>>,
    token_consumed: bool,
}

impl<'a> Parser<'a> {
    pub(crate) fn new(
        expression: String,
        tokens: &'a mut dyn Iterator<Item = Result<Token, NmpolicyError>>,
    ) -> Self {
        Self {
            expression,
            tokens,
            current_token: None,
            root_node: None,
            piped_in_node: None,
            token_consumed: true,
        }
    }

    pub(crate) fn parse(&mut self) -> Result<Option<Box<Node>>, NmpolicyError> {
        match self.parse_tokens() {
            Ok(root_node) => Ok(root_node),
            Err(mut e) => {
                if let Some(current_token) = self.current_token.clone() {
                    e.decorate(self.expression.clone(), current_token.pos())
                }
                Err(e)
            }
        }
    }

    fn parse_tokens(&mut self) -> Result<Option<Box<Node>>, NmpolicyError> {
        loop {
            match self.next_token() {
                Some(Ok(_)) => match self.parse_token() {
                    Ok(()) => continue,
                    Err(e) => return Err(e),
                },
                Some(Err(e)) => {
                    return Err(e);
                }
                None => break,
            }
        }
        match self.piped_in_node.clone() {
            None => Ok(self.root_node.clone()),
            Some(_) => Err(NmpolicyError::new(
                ErrorKind::InvalidPipeMissingRightExpression,
            )),
        }
    }

    fn parse_token(&mut self) -> Result<(), NmpolicyError> {
        match self.current_token.clone() {
            Some(Token::Identity(pos, literal)) => self.parse_path(pos, literal),
            Some(Token::EqFilter(pos)) => self.parse_eqfilter(pos),
            Some(Token::Replace(pos)) => self.parse_replace(pos),
            Some(Token::Pipe(pos)) => self.parse_pipe(pos),
            Some(Token::Str(pos, literal)) => self.parse_string(pos, literal),
            Some(t) => Err(NmpolicyError::new(
                ErrorKind::InvalidExpressionUnexpectedToken(t.literal()),
            )),
            None => match self.piped_in_node.clone() {
                None => Ok(()),
                Some(_) => Err(NmpolicyError::new(
                    ErrorKind::InvalidPipeMissingRightExpression,
                )),
            },
        }
    }

    fn parse_path(&mut self, pos: usize, literal: String) -> Result<(), NmpolicyError> {
        let mut steps = vec![*self.identity(pos, literal)];
        loop {
            match self.next_token() {
                Some(Ok(Token::Dot(_))) => match self.next_token() {
                    Some(Ok(Token::Identity(npos, nliteral))) => {
                        steps.push(*self.identity(npos, nliteral))
                    }
                    Some(Ok(Token::Number(npos, nliteral))) => {
                        steps.push(*self.number(npos, nliteral))
                    }
                    Some(Ok(_)) | None => {
                        return Err(NmpolicyError::new(
                            ErrorKind::InvalidPathUnexpectedTokenAfterDot,
                        ));
                    }
                    Some(Err(e)) => return Err(e),
                },
                Some(Ok(
                    Token::EqFilter(_) | Token::Replace(_) | Token::Merge(_) | Token::Pipe(_),
                )) => {
                    self.token_consumed = false;
                    break;
                }
                None => break,
                Some(Ok(_)) => return Err(NmpolicyError::new(ErrorKind::InvalidPathMissingDot)),
                Some(Err(e)) => return Err(e),
            }
        }

        self.path(pos, steps);
        Ok(())
    }

    fn parse_string(&mut self, pos: usize, literal: String) -> Result<(), NmpolicyError> {
        self.root_node = Some(string(pos, literal));
        Ok(())
    }

    fn parse_eqfilter(&mut self, pos: usize) -> Result<(), NmpolicyError> {
        match self.fill_in_ternary_operator("equality filter") {
            Ok((value1, value2, value3)) => {
                self.eqfilter(pos, value1.unwrap(), value2.unwrap(), value3.unwrap());
                Ok(())
            }
            Err(e) => Err(e),
        }
    }

    fn parse_replace(&mut self, pos: usize) -> Result<(), NmpolicyError> {
        match self.fill_in_ternary_operator("replace") {
            Ok((value1, value2, value3)) => {
                self.replace(pos, value1.unwrap(), value2.unwrap(), value3.unwrap());
                Ok(())
            }
            Err(e) => Err(e),
        }
    }

    fn parse_pipe(&mut self, _: usize) -> Result<(), NmpolicyError> {
        match self.root_node.clone() {
            Some(root_node) => match root_node.kind {
                NodeKind::Path(_) => {
                    self.piped_in_node = self.root_node.clone();
                    Ok(())
                }
                _ => Err(NmpolicyError::new(ErrorKind::InvalidPipeMissingLeftPath)),
            },
            None => Err(NmpolicyError::new(
                ErrorKind::InvalidPipeMissingLeftExpression,
            )),
        }
    }

    fn fill_in_ternary_operator(
        &mut self,
        operator_name: &'static str,
    ) -> Result<TernaryOperator, NmpolicyError> {
        let mut values: TernaryOperator =
            (Default::default(), Default::default(), Default::default());
        match self.root_node.clone() {
            Some(root_node) => match root_node.kind {
                NodeKind::Path(steps) => {
                    values.0 = match self.piped_in_node.clone() {
                        Some(piped_in_node) => {
                            self.piped_in_node = None;
                            Some(piped_in_node)
                        }
                        None => Some(current_state(0)),
                    };
                    values.1 = Some(self.path(root_node.pos, steps));
                    match self.next_token() {
                        Some(Ok(Token::Str(pos, literal))) => {
                            values.2 = Some(self.string(pos, literal));
                            Ok(values)
                        }
                        Some(Ok(Token::Identity(pos, literal))) => {
                            match self.parse_path(pos, literal) {
                                Ok(()) => {
                                    values.2 = self.root_node.clone();
                                    Ok(values)
                                }
                                Err(e) => Err(e),
                            }
                        }
                        Some(Ok(_)) => Err(NmpolicyError::new(
                            ErrorKind::InvalidTernaryUnexpectedRightHand(operator_name),
                        )),
                        Some(Err(e)) => Err(e),
                        None => Err(NmpolicyError::new(
                            ErrorKind::InvalidTernaryMissingRightHand(operator_name),
                        )),
                    }
                }
                _ => Err(NmpolicyError::new(
                    ErrorKind::InvalidTernaryUnexpectedLeftHand(operator_name),
                )),
            },
            None => Err(NmpolicyError::new(
                ErrorKind::InvalidTernaryMissingLeftHand(operator_name),
            )),
        }
    }
    fn next_token(&mut self) -> Option<Result<Token, NmpolicyError>> {
        if !self.token_consumed {
            self.token_consumed = true;
            return Some(Ok(self.current_token.clone().unwrap()));
        }
        self.token_consumed = true;
        let next_result = self.tokens.next();
        match next_result {
            Some(Ok(current_token)) => {
                self.current_token = Some(current_token.clone());
                Some(Ok(current_token))
            }
            Some(Err(e)) => Some(Err(e)),
            None => None,
        }
    }
    fn identity(&mut self, pos: usize, literal: String) -> Box<Node> {
        self.root_node = Some(identity(pos, literal));
        self.root_node.clone().unwrap()
    }

    fn string(&mut self, pos: usize, literal: String) -> Box<Node> {
        self.root_node = Some(string(pos, literal));
        self.root_node.clone().unwrap()
    }

    fn number(&mut self, pos: usize, literal: i32) -> Box<Node> {
        self.root_node = Some(number(pos, literal));
        self.root_node.clone().unwrap()
    }

    fn path(&mut self, pos: usize, steps: Vec<Node>) -> Box<Node> {
        self.root_node = Some(path(pos, steps));
        self.root_node.clone().unwrap()
    }

    fn eqfilter(
        &mut self,
        pos: usize,
        value1: Box<Node>,
        value2: Box<Node>,
        value3: Box<Node>,
    ) -> Box<Node> {
        self.root_node = Some(eqfilter(pos, value1, value2, value3));
        self.root_node.clone().unwrap()
    }

    fn replace(
        &mut self,
        pos: usize,
        value1: Box<Node>,
        value2: Box<Node>,
        value3: Box<Node>,
    ) -> Box<Node> {
        self.root_node = Some(replace(pos, value1, value2, value3));
        self.root_node.clone().unwrap()
    }
}
