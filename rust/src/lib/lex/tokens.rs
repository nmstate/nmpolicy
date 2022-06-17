/*
 * Copyright 2021 NMPolicy Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *	  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

use crate::error::{ErrorKind, NmpolicyError};

use std::{iter::Peekable, str::CharIndices};

#[derive(Clone, Eq, PartialEq, Debug)]
pub(crate) enum Token {
    Identity(usize, String),
    Number(usize, u32),
    Str(usize, String),
    Dot(usize),      // .
    Pipe(usize),     // |
    Replace(usize),  // !=
    EqFilter(usize), // ==
    Merge(usize),    // +
    True(usize),     // true
    False(usize),    // false
}

pub(crate) struct Tokens<'a> {
    input: &'a str,
    char_indices: Peekable<CharIndices<'a>>,
    tokens: Vec<Token>,
    pos: usize,
    has_error: bool,
}

impl<'a> Tokens<'a> {
    pub(crate) fn new(input: &'a str) -> Self {
        Self {
            input,
            char_indices: input.char_indices().peekable(),
            tokens: Vec::new(),
            pos: 0,
            has_error: false,
        }
    }

    fn next_token(&mut self) -> Option<Result<Token, NmpolicyError>> {
        if self.has_error {
            return None;
        }
        while self.char_indices.next_if(|(_, ch)| *ch == ' ') != None {}
        match self.char_indices.next() {
            Some((pos, ch)) => {
                self.pos = pos;
                let result = match ch {
                    '|' => self.tokenize_pipe(ch),
                    '.' => self.tokenize_dot(ch),
                    '+' => self.tokenize_plus(ch),
                    '=' => self.tokenize_equal(ch),
                    ':' => self.tokenize_colon(ch),
                    '"' | '\'' => self.tokenize_string(ch),
                    ch if ch.is_numeric() => self.tokenize_number(ch),
                    ch if ch.is_alphabetic() => self.tokenize_identity(ch),
                    ch => Err(NmpolicyError::new(
                        ErrorKind::InvalidArgument,
                        format!("illegal char {ch}"),
                    )),
                };
                match result {
                    Ok(token) => Some(Ok(token)),
                    Err(mut e) => {
                        self.has_error = true;
                        e.decorate(self.input.to_string(), self.pos);
                        Some(Err(e))
                    }
                }
            }
            None => None,
        }
    }

    fn tokenize_plus(&self, _: char) -> Result<Token, NmpolicyError> {
        Ok(Token::Merge(self.pos))
    }
    fn tokenize_dot(&self, _: char) -> Result<Token, NmpolicyError> {
        Ok(Token::Dot(self.pos))
    }
    fn tokenize_pipe(&self, _: char) -> Result<Token, NmpolicyError> {
        Ok(Token::Pipe(self.pos))
    }
    fn tokenize_equal(&mut self, _: char) -> Result<Token, NmpolicyError> {
        match self.char_indices.next() {
            Some((_, '=')) => Ok(Token::EqFilter(self.pos)),
            Some((pos, nch)) => {
                self.pos = pos;
                Err(NmpolicyError::new(
                    ErrorKind::InvalidArgument,
                    format!("invalid EQFILTER operation format ({nch} is not equal char)"),
                ))
            }
            None => Err(NmpolicyError::new(
                ErrorKind::InvalidArgument,
                "invalid EQFILTER operation format (EOF)".to_string(),
            )),
        }
    }
    fn tokenize_colon(&mut self, _: char) -> Result<Token, NmpolicyError> {
        match self.char_indices.next() {
            Some((_, '=')) => Ok(Token::Replace(self.pos)),
            Some((pos, nch)) => {
                self.pos = pos;
                Err(NmpolicyError::new(
                    ErrorKind::InvalidArgument,
                    format!("invalid REPLACE operation format ({nch} is not equal char)"),
                ))
            }
            None => Err(NmpolicyError::new(
                ErrorKind::InvalidArgument,
                "invalid REPLACE operation format (EOF)".to_string(),
            )),
        }
    }
    fn tokenize_string(&mut self, delimiter: char) -> Result<Token, NmpolicyError> {
        let mut last_matched: char = '\0';
        let begin_pos = self.pos;
        let mut last_pos = self.pos;
        let s: String = self
            .char_indices
            .by_ref()
            .take_while(|(pos, c)| {
                last_matched = *c;
                last_pos = *pos;
                *c != delimiter
            })
            .map(|(_, c)| c)
            .collect();

        self.pos = last_pos;

        if last_matched == delimiter {
            Ok(Token::Str(begin_pos, s))
        } else {
            Err(NmpolicyError::new(
                ErrorKind::InvalidArgument,
                format!("invalid string format (missing {delimiter} terminator)"),
            ))
        }
    }
    fn tokenize_number(&mut self, ch: char) -> Result<Token, NmpolicyError> {
        let mut number = String::from(ch);
        let begin_pos = self.pos;
        while let Some((pos, ch)) = self.char_indices.next_if(|(_, ch)| ch.is_numeric()) {
            number.push(ch);
            self.pos = pos;
        }
        match self.char_indices.peek() {
            Some((_, ' ' | '.' | ':' | '+' | '|' | '=')) | None => {
                Ok(Token::Number(begin_pos, number.parse::<u32>().unwrap()))
            }
            Some((pos, ch)) => {
                self.pos = *pos;
                Err(NmpolicyError::new(
                    ErrorKind::InvalidArgument,
                    format!("invalid number format ({ch} is not a digit)"),
                ))
            }
        }
    }
    fn tokenize_identity(&mut self, ch: char) -> Result<Token, NmpolicyError> {
        let mut identity = String::from(ch);
        let begin_pos = self.pos;
        while let Some((pos, ch)) = self
            .char_indices
            .next_if(|(_, ch)| ch.is_alphanumeric() || *ch == '-')
        {
            identity.push(ch);
            self.pos = pos;
        }
        match self.char_indices.peek() {
            Some((_, ' ' | '.' | ':' | '+' | '|' | '=')) | None => {
                Ok(Token::Identity(begin_pos, identity))
            }
            Some((_, ch)) => Err(NmpolicyError::new(
                ErrorKind::InvalidArgument,
                format!("invalid identity format ({ch} is not a digit, letter or -)"),
            )),
        }
    }
}

impl<'a> Iterator for Tokens<'a> {
    type Item = Result<Token, NmpolicyError>;
    fn next(&mut self) -> Option<Self::Item> {
        self.next_token()
    }
}
