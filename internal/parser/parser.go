package parser

import (
	"fmt"
	"strconv"

	"github.com/nmstate/nmpolicy/internal/ast"
	"github.com/nmstate/nmpolicy/internal/lexer"
	"github.com/nmstate/nmpolicy/internal/source"
)

type Parser struct {
	source       source.Source
	lexer        lexer.Lexer
	currentToken *lexer.Token
	readNext     bool
}

func New(source source.Source, lexer lexer.Lexer) Parser {
	return Parser{source: source, lexer: lexer, currentToken: nil, readNext: true}
}

func (p *Parser) Parse() (*ast.Command, error) {
	cmd, err := p.parseCommand()
	if err != nil {
		if p.currentToken != nil {
			return nil, source.NewError(err, p.currentToken.Position).Decorate(p.source)
		}
		return nil, err
	}
	return cmd, nil
}

func (p *Parser) parseCommand() (*ast.Command, error) {
	var cmd *ast.Command
	for {
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		tokenType := p.currentToken.Type
		if tokenType == lexer.EOF {
			if cmd == nil {
				return nil, fmt.Errorf("missing command")
			}
			return cmd, nil
		} else if tokenType == lexer.IDENTITY {
			if cmd == nil {
				cmd = &ast.Command{Node: p.newNodeWithCurrentPosition()}
			}
			path, err := p.parsePath()
			if err != nil {
				return nil, err
			}
			cmd.Path = path
		} else if tokenType.IsOperator() {
			if cmd == nil {
				return nil, badOperationFormat(tokenType, "missing right hand argument")
			}
			argument, err := p.parseArgument()
			if err != nil {
				return nil, badOperationFormat(tokenType, fmt.Sprintf("%v", err))
			}
			if cmd == nil {
				cmd = &ast.Command{Node: p.newNodeWithCurrentPosition()}
			}
			arguments := []ast.Argument{
				{Node: cmd.Node, Path: cmd.Path},
				*argument,
			}

			switch tokenType {
			case lexer.EQUALITY:
				cmd.Equal = arguments
			case lexer.ASSIGN:
				cmd.Assign = arguments
			case lexer.MERGE:
				if argument.Path == nil {
					return nil, badOperationFormat(tokenType, "only paths can be merged")
				}
				cmd.Merge = arguments
			}
			cmd.Path = nil
		} else {
			return nil, fmt.Errorf("parsing token %s not implemented yet", p.currentToken.Type)
		}
	}
}

func (p *Parser) nextToken() error {
	if !p.readNext {
		p.readNext = true
		return nil
	}

	token, err := p.lexer.Lex()
	if err != nil {
		return fmt.Errorf("failed lexing next token: %w", err)
	}
	p.currentToken = token
	return nil
}

func (p *Parser) parseArgument() (*ast.Argument, error) {
	err := p.nextToken()
	if err != nil {
		return nil, err
	}
	argument := &ast.Argument{Node: p.newNodeWithCurrentPosition()}
	switch p.currentToken.Type {
	case lexer.IDENTITY:
		path, err := p.parsePath()
		if err != nil {
			return nil, fmt.Errorf("failed parsing path argument: %w", err)
		}
		argument.Path = path
		return argument, nil
	case lexer.STRING:
		argument.String = &p.currentToken.Literal
		return argument, nil
	case lexer.NUMBER:
		number, err := strconv.Atoi(p.currentToken.Literal)
		if err != nil {
			return nil, fmt.Errorf("failed parsing number argument: %w", err)
		}
		argument.Number = &number
		return argument, nil
	default:
		return nil, fmt.Errorf("supported argument missing")
	}
}

func (p *Parser) parsePath() (ast.Path, error) {
	path := ast.Path{{Node: p.newNodeWithCurrentPosition()}}
	var lastToken *lexer.Token
	for {
		lastToken = p.currentToken
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		switch p.currentToken.Type {
		case lexer.PATH:
			continue
		case lexer.IDENTITY:
			if lastToken.Type != lexer.PATH {
				return nil, badPathFormatError("identities has to be separated by a dot")
			}
			path = append(path, ast.Step{Node: p.newNodeWithCurrentPosition(), Identity: &p.currentToken.Literal})
		case lexer.NUMBER:
			if lastToken.Type != lexer.PATH {
				return nil, badPathFormatError("indexes has to be separated by a dot")
			}
			index, err := strconv.Atoi(p.currentToken.Literal)
			if err != nil {
				return nil, badPathFormatError("number token is not an interger")
			}
			path = append(path, ast.Step{Node: p.newNodeWithCurrentPosition(), Index: &index})
		case lexer.STRING:
			if lastToken.Type == lexer.PATH {
				return nil, badPathFormatError(`string can only be used with brackets, ej foo["bar"]`)
			}
		case lexer.LBRACKET:
			return nil, badPathFormatError("brackets path step not implemented")
		case lexer.RBRACKET:
			return nil, badPathFormatError("brackets path step not implemented")
		default:
			if lastToken.Type == lexer.PATH {
				return nil, badPathFormatError("dot has to be followed by indentity or index")
			}
			p.readNext = false
			return path, nil

		}
	}
}

func (p *Parser) newNodeWithCurrentPosition() ast.Node {
	return ast.Node{Position: p.currentToken.Position}
}
