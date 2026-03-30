package parser

import (
	"fmt"
	"strconv"

	"github.com/robogg133/MoonMS/mscript/interpreter/ast"
	"github.com/robogg133/MoonMS/mscript/interpreter/lexer"
)

type Parser struct {
	l       *lexer.Lexer
	current lexer.Token
	peek    lexer.Token
	errors  []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.advance()
	p.advance()
	return p
}

type PluginMeta struct {
	Identifier string
	Name       string
	Version    string
	Exports    []ast.FnSignature
	Emits      []ast.FnSignature
}

func ExtractMeta(prog *ast.Program) (*PluginMeta, error) {
	for _, node := range prog.Statements {
		decl, ok := node.(*ast.PluginDecl)
		if !ok {
			continue
		}

		return &PluginMeta{
			Name:    decl.Name,
			Version: decl.Version,
			Exports: decl.Exports,
			Emits:   decl.Emits,
		}, nil
	}
	return nil, fmt.Errorf("plugin declaration not found")
}

func (p *Parser) Errors() []string { return p.errors }
func (p *Parser) HasErrors() bool  { return len(p.errors) > 0 }

func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{}
	for !p.is(lexer.TOKEN_EOF) {
		node := p.parseTopLevel()
		if node != nil {
			prog.Statements = append(prog.Statements, node)
		}
	}
	return prog
}

func (p *Parser) advance() {
	p.current = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) expect(t lexer.TokenType) bool {
	if p.current.Type != t {
		p.errors = append(p.errors, fmt.Sprintf(
			"expected token %d, got %d (%q)", t, p.current.Type, p.current.Value,
		))
		return false
	}
	p.advance()
	return true
}

func (p *Parser) is(t lexer.TokenType) bool     { return p.current.Type == t }
func (p *Parser) peekIs(t lexer.TokenType) bool { return p.peek.Type == t }

// ── Top-level ────────────────────────────────────────────

func (p *Parser) parseTopLevel() ast.Node {
	switch p.current.Type {
	case lexer.TOKEN_PLUGIN:
		return p.parsePlugin()
	case lexer.TOKEN_REQUIRE:
		return p.parseRequire()
	case lexer.TOKEN_FN:
		return p.parseFnDecl()
	case lexer.TOKEN_COMMAND:
		return p.parseCommand()
	case lexer.TOKEN_ON:
		return p.parseEvent()
	case lexer.TOKEN_LET:
		return p.parseLet()
	default:
		p.errors = append(p.errors, fmt.Sprintf(
			"unexpected token %q at top level", p.current.Value,
		))
		p.advance()
		return nil
	}
}

// ── plugin "name" version "x" { export { } emits { } } ──

func (p *Parser) parsePlugin() *ast.PluginDecl {
	p.advance() // consume 'plugin'

	name := p.current.Value
	if !p.expect(lexer.TOKEN_STRING) {
		return nil
	}

	if !p.expect(lexer.TOKEN_VERSION) {
		return nil
	}

	version := p.current.Value
	if !p.expect(lexer.TOKEN_STRING) {
		return nil
	}

	if !p.expect(lexer.TOKEN_LBRACE) {
		return nil
	}

	decl := &ast.PluginDecl{Name: name, Version: version}

	for !p.is(lexer.TOKEN_RBRACE) && !p.is(lexer.TOKEN_EOF) {
		switch p.current.Type {
		case lexer.TOKEN_EXPORT:
			p.advance()
			p.expect(lexer.TOKEN_LBRACE)
			for !p.is(lexer.TOKEN_RBRACE) && !p.is(lexer.TOKEN_EOF) {
				decl.Exports = append(decl.Exports, p.parseFnSignature())
			}
			p.expect(lexer.TOKEN_RBRACE)

		case lexer.TOKEN_EMITS:
			p.advance()
			p.expect(lexer.TOKEN_LBRACE)
			for !p.is(lexer.TOKEN_RBRACE) && !p.is(lexer.TOKEN_EOF) {
				decl.Emits = append(decl.Emits, p.parseFnSignature())
			}
			p.expect(lexer.TOKEN_RBRACE)

		default:
			p.errors = append(p.errors, fmt.Sprintf(
				"unexpected %q inside plugin block", p.current.Value,
			))
			p.advance()
		}
	}

	p.expect(lexer.TOKEN_RBRACE)
	return decl
}

// name(params...): returnType
func (p *Parser) parseFnSignature() ast.FnSignature {
	name := p.current.Value
	p.expect(lexer.TOKEN_IDENT)
	p.expect(lexer.TOKEN_LPAREN)
	params := p.parseParams()
	p.expect(lexer.TOKEN_RPAREN)

	returnType := ""
	if p.is(lexer.TOKEN_COLON) {
		p.advance()
		returnType = p.current.Value
		p.advance()
	}

	return ast.FnSignature{Name: name, Params: params, ReturnType: returnType}
}

func (p *Parser) parseRequire() *ast.RequireDecl {
	p.advance()

	path := p.current.Value
	p.expect(lexer.TOKEN_STRING)

	alias := path
	if p.is(lexer.TOKEN_AS) {
		p.advance()
		alias = p.current.Value
		p.expect(lexer.TOKEN_IDENT)
	}

	return &ast.RequireDecl{Path: path, Alias: alias}
}

// ── fn name(params): type { body } ───────────────────────

func (p *Parser) parseFnDecl() *ast.FnDecl {
	p.advance() // consume 'fn'

	name := p.current.Value
	p.expect(lexer.TOKEN_IDENT)
	p.expect(lexer.TOKEN_LPAREN)
	params := p.parseParams()
	p.expect(lexer.TOKEN_RPAREN)

	returnType := ""
	if p.is(lexer.TOKEN_COLON) {
		p.advance()
		returnType = p.current.Value
		p.advance()
	}

	body := p.parseBlock()
	return &ast.FnDecl{Name: name, Params: params, ReturnType: returnType, Body: body}
}

// ── command "name" (params) { body } ─────────────────────

func (p *Parser) parseCommand() *ast.CommandDecl {
	p.advance() // consume 'command'

	name := p.current.Value
	p.expect(lexer.TOKEN_STRING)
	p.expect(lexer.TOKEN_LPAREN)

	var params []string
	for !p.is(lexer.TOKEN_RPAREN) && !p.is(lexer.TOKEN_EOF) {
		params = append(params, p.current.Value)
		p.expect(lexer.TOKEN_IDENT)
		if p.is(lexer.TOKEN_COMMA) {
			p.advance()
		}
	}
	p.expect(lexer.TOKEN_RPAREN)

	body := p.parseBlock()
	return &ast.CommandDecl{Name: name, Params: params, Body: body}
}

// ── on event(param) { }  /  on plugin.event { } ──────────

func (p *Parser) parseEvent() *ast.EventDecl {
	p.advance() // consume 'on'

	// first ident
	first := p.current.Value
	p.expect(lexer.TOKEN_IDENT)

	source := ""
	event := first

	// on tp.killed  →  source="tp", event="killed"
	if p.is(lexer.TOKEN_DOT) {
		p.advance()
		source = first
		event = p.current.Value
		p.expect(lexer.TOKEN_IDENT)
	}

	// opcional: on player_join(player)
	param := ""
	if p.is(lexer.TOKEN_LPAREN) {
		p.advance()
		param = p.current.Value
		p.expect(lexer.TOKEN_IDENT)
		p.expect(lexer.TOKEN_RPAREN)
	}

	body := p.parseBlock()
	return &ast.EventDecl{Source: source, Event: event, Param: param, Body: body}
}

// ── Block ─────────────────────────────────────────────────

func (p *Parser) parseBlock() []ast.Node {
	p.expect(lexer.TOKEN_LBRACE)
	var stmts []ast.Node
	for !p.is(lexer.TOKEN_RBRACE) && !p.is(lexer.TOKEN_EOF) {
		stmts = append(stmts, p.parseStatement())
	}
	p.expect(lexer.TOKEN_RBRACE)
	return stmts
}

// ── Statements ────────────────────────────────────────────

func (p *Parser) parseStatement() ast.Node {
	switch p.current.Type {
	case lexer.TOKEN_IF:
		return p.parseIf()
	case lexer.TOKEN_AFTER:
		return p.parseAfter()
	case lexer.TOKEN_EVERY:
		return p.parseEvery()
	case lexer.TOKEN_RETURN:
		return p.parseReturn()
	case lexer.TOKEN_CANCEL:
		p.advance()
		return &ast.CancelStmt{}
	case lexer.TOKEN_LET:
		return p.parseLet()
	case lexer.TOKEN_EMIT:
		return p.parseEmit()
	default:
		return p.parseExprOrAssign()
	}
}

// if cond { } else { }
func (p *Parser) parseIf() *ast.IfStmt {
	p.advance() // consume 'if'
	cond := p.parseExpr()
	then := p.parseBlock()

	var els []ast.Node
	if p.is(lexer.TOKEN_ELSE) {
		p.advance()
		els = p.parseBlock()
	}

	return &ast.IfStmt{Condition: cond, Then: then, Else: els}
}

// after expr { }
func (p *Parser) parseAfter() *ast.AfterStmt {
	p.advance() // consume 'after'
	duration := p.parseExpr()
	body := p.parseBlock()
	return &ast.AfterStmt{Duration: duration, Body: body}
}

// every expr { }
func (p *Parser) parseEvery() *ast.EveryStmt {
	p.advance() // consume 'every'
	interval := p.parseExpr()
	body := p.parseBlock()
	return &ast.EveryStmt{Interval: interval, Body: body}
}

// return / return expr
func (p *Parser) parseReturn() *ast.ReturnStmt {
	p.advance() // consume 'return'

	if p.is(lexer.TOKEN_RBRACE) || p.is(lexer.TOKEN_EOF) {
		return &ast.ReturnStmt{Value: nil}
	}

	return &ast.ReturnStmt{Value: p.parseExpr()}
}

// let name = expr
func (p *Parser) parseLet() *ast.LetStmt {
	p.advance() // consume 'let'
	name := p.current.Value
	p.expect(lexer.TOKEN_IDENT)
	p.expect(lexer.TOKEN_ASSIGN)
	val := p.parseExpr()
	return &ast.LetStmt{Name: name, Value: val}
}

// emit event(args)
func (p *Parser) parseEmit() *ast.EmitStmt {
	p.advance() // consume 'emit'
	name := p.current.Value
	p.expect(lexer.TOKEN_IDENT)
	p.expect(lexer.TOKEN_LPAREN)
	args := p.parseArgs()
	p.expect(lexer.TOKEN_RPAREN)
	return &ast.EmitStmt{Event: name, Args: args}
}

func (p *Parser) parseExprOrAssign() ast.Node {
	expr := p.parseExpr()

	switch p.current.Type {
	case lexer.TOKEN_ASSIGN, lexer.TOKEN_PLUS_ASSIGN, lexer.TOKEN_MINUS_ASSIGN:
		op := p.current.Value
		p.advance()
		val := p.parseExpr()
		return &ast.AssignStmt{Target: expr, Op: op, Value: val}
	}

	return &ast.ExprStmt{Expr: expr}
}

// ── Expressões ────────────────────────────────────────────
//
// higher to lower
//   1. is
//   2. == != > < >= <=
//   3. + -
//   4. * / %
//   5. unary: not !
//   6. postfix: call . member
//   7. primary

func (p *Parser) parseExpr() ast.Node { return p.parseIs() }

func (p *Parser) parseIs() ast.Node {
	left := p.parseComparison()
	if p.is(lexer.TOKEN_IS) {
		p.advance()
		typeName := p.current.Value
		p.advance()
		return &ast.IsExpr{Value: left, TypeName: typeName}
	}
	return left
}

func (p *Parser) parseComparison() ast.Node {
	left := p.parseAddSub()
	for p.current.Type == lexer.TOKEN_EQUALS ||
		p.current.Type == lexer.TOKEN_NOT_EQ ||
		p.current.Type == lexer.TOKEN_GREATER ||
		p.current.Type == lexer.TOKEN_LESS ||
		p.current.Type == lexer.TOKEN_GTE ||
		p.current.Type == lexer.TOKEN_LTE {
		op := p.current.Value
		p.advance()
		right := p.parseAddSub()
		left = &ast.BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseAddSub() ast.Node {
	left := p.parseMulDiv()
	for p.is(lexer.TOKEN_PLUS) || p.is(lexer.TOKEN_MINUS) {
		op := p.current.Value
		p.advance()
		right := p.parseMulDiv()
		left = &ast.BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseMulDiv() ast.Node {
	left := p.parseUnary()
	for p.is(lexer.TOKEN_ASTERISK) || p.is(lexer.TOKEN_SLASH) || p.is(lexer.TOKEN_REST) {
		op := p.current.Value
		p.advance()
		right := p.parseUnary()
		left = &ast.BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseUnary() ast.Node {
	if p.is(lexer.TOKEN_NOT) || p.is(lexer.TOKEN_NOT_KW) {
		op := p.current.Value
		p.advance()
		return &ast.UnaryExpr{Op: op, Operand: p.parseUnary()}
	}
	return p.parsePostfix()
}

func (p *Parser) parsePostfix() ast.Node {
	node := p.parsePrimary()
	for {
		if p.is(lexer.TOKEN_DOT) {
			p.advance()
			field := p.current.Value
			p.advance()
			if p.is(lexer.TOKEN_LPAREN) {
				p.advance()
				args := p.parseArgs()
				p.expect(lexer.TOKEN_RPAREN)
				node = &ast.CallExpr{
					Callee: &ast.MemberExpr{Object: node, Field: field},
					Args:   args,
				}
			} else {
				node = &ast.MemberExpr{Object: node, Field: field}
			}
		} else if p.is(lexer.TOKEN_LPAREN) {
			p.advance()
			args := p.parseArgs()
			p.expect(lexer.TOKEN_RPAREN)
			node = &ast.CallExpr{Callee: node, Args: args}
		} else {
			break
		}
	}
	return node
}

func (p *Parser) parsePrimary() ast.Node {
	switch p.current.Type {
	case lexer.TOKEN_NUMBER:
		v, _ := strconv.ParseFloat(p.current.Value, 64)
		p.advance()
		return &ast.NumberLit{Value: v}

	case lexer.TOKEN_STRING:
		v := p.current.Value
		p.advance()
		return &ast.StringLit{Value: v}

	case lexer.TOKEN_TRUE:
		p.advance()
		return &ast.BoolLit{Value: true}

	case lexer.TOKEN_FALSE:
		p.advance()
		return &ast.BoolLit{Value: false}

	case lexer.TOKEN_NULL:
		p.advance()
		return &ast.NullLit{}

	case lexer.TOKEN_LPAREN:
		p.advance()
		expr := p.parseExpr()
		p.expect(lexer.TOKEN_RPAREN)
		return expr

	case lexer.TOKEN_LBRACE:
		return p.parseObjectLit()

	case lexer.TOKEN_IDENT:
		name := p.current.Value
		p.advance()
		return &ast.Ident{Name: name}

	default:
		p.errors = append(p.errors, fmt.Sprintf(
			"unexpected token %q in expression", p.current.Value,
		))
		p.advance()
		return &ast.NullLit{}
	}
}

// { X: 123, Y: 80, Z: -234 }
func (p *Parser) parseObjectLit() *ast.ObjectLit {
	p.expect(lexer.TOKEN_LBRACE)
	var fields []ast.ObjectField
	for !p.is(lexer.TOKEN_RBRACE) && !p.is(lexer.TOKEN_EOF) {
		key := p.current.Value
		p.expect(lexer.TOKEN_IDENT)
		p.expect(lexer.TOKEN_COLON)
		val := p.parseExpr()
		fields = append(fields, ast.ObjectField{Key: key, Value: val})
		if p.is(lexer.TOKEN_COMMA) {
			p.advance()
		}
	}
	p.expect(lexer.TOKEN_RBRACE)
	return &ast.ObjectLit{Fields: fields}
}

// ── Helpers ───────────────────────────────────────────────

func (p *Parser) parseArgs() []ast.Node {
	var args []ast.Node
	for !p.is(lexer.TOKEN_RPAREN) && !p.is(lexer.TOKEN_EOF) {
		args = append(args, p.parseExpr())
		if p.is(lexer.TOKEN_COMMA) {
			p.advance()
		}
	}
	return args
}

// name / name: type
func (p *Parser) parseParams() []ast.Param {
	var params []ast.Param
	for !p.is(lexer.TOKEN_RPAREN) && !p.is(lexer.TOKEN_EOF) {
		name := p.current.Value
		p.expect(lexer.TOKEN_IDENT)
		typeName := ""
		if p.is(lexer.TOKEN_COLON) {
			p.advance()
			typeName = p.current.Value
			p.advance()
		}
		params = append(params, ast.Param{Name: name, Type: typeName})
		if p.is(lexer.TOKEN_COMMA) {
			p.advance()
		}
	}
	return params
}
