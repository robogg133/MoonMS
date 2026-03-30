package ast

type Node interface{ nodeType() string }

type Program struct {
	Statements []Node
}

// plugin "economy" version "1.0" { export { ... } emits { ... } }
type PluginDecl struct {
	Name    string
	Version string
	Exports []FnSignature
	Emits   []FnSignature
}

// assign func: nome(params...): returnType
type FnSignature struct {
	Name       string
	Params     []Param
	ReturnType string // "" = void
}

type Param struct {
	Name string
	Type string // "" = sem tipo
}

// require teleport as tp
type RequireDecl struct {
	Path  string // "teleport"
	Alias string // "tp"
}

// fn name(params): returnType { body }
type FnDecl struct {
	Name       string
	Params     []Param
	ReturnType string
	Body       []Node
}

// command "pay" (name, amount) { body }
type CommandDecl struct {
	Name   string
	Params []string
	Body   []Node
}

// on event_name(param) { body }
// on event_name { body }
// on plugin.event { body }
type EventDecl struct {
	Source string // "" = local, "tp" = tp.killed
	Event  string
	Param  string // "" = sem param
	Body   []Node
}

// ── Statements ────────────────────────────────────────────

// if cond { then } else { els }
type IfStmt struct {
	Condition Node
	Then      []Node
	Else      []Node
}

// after expr { body }
type AfterStmt struct {
	Duration Node
	Body     []Node
}

// every expr { body }
type EveryStmt struct {
	Interval Node
	Body     []Node
}

// return expr
type ReturnStmt struct {
	Value Node // nil = return void
}

// cancel
type CancelStmt struct{}

// let name = expr
type LetStmt struct {
	Name  string
	Value Node
}

// target = expr / target += expr
type AssignStmt struct {
	Target Node
	Op     string
	Value  Node
}

// emit name(args)
type EmitStmt struct {
	Event string
	Args  []Node
}

// expr statement
type ExprStmt struct {
	Expr Node
}

// ── Expressions ──────────────────────────────────────────

type NumberLit struct{ Value float64 }
type StringLit struct{ Value string }
type BoolLit struct{ Value bool }
type NullLit struct{}

// usual ident
type Ident struct{ Name string }

// member expr
type MemberExpr struct {
	Object Node
	Field  string
}

// call: fn(args) ou obj.fn(args)
type CallExpr struct {
	Callee Node
	Args   []Node
}

type BinaryExpr struct {
	Op    string
	Left  Node
	Right Node
}

type UnaryExpr struct {
	Op      string
	Operand Node
}

// expr is Type
type IsExpr struct {
	Value    Node
	TypeName string
}

// { X: 123, Y: 80 }
type ObjectLit struct {
	Fields []ObjectField
}

type ObjectField struct {
	Key   string
	Value Node
}

// ── nodeType ─────────────────────────────────────────────

func (n *Program) nodeType() string     { return "Program" }
func (n *PluginDecl) nodeType() string  { return "PluginDecl" }
func (n *RequireDecl) nodeType() string { return "RequireDecl" }
func (n *FnDecl) nodeType() string      { return "FnDecl" }
func (n *CommandDecl) nodeType() string { return "CommandDecl" }
func (n *EventDecl) nodeType() string   { return "EventDecl" }
func (n *IfStmt) nodeType() string      { return "IfStmt" }
func (n *AfterStmt) nodeType() string   { return "AfterStmt" }
func (n *EveryStmt) nodeType() string   { return "EveryStmt" }
func (n *ReturnStmt) nodeType() string  { return "ReturnStmt" }
func (n *CancelStmt) nodeType() string  { return "CancelStmt" }
func (n *LetStmt) nodeType() string     { return "LetStmt" }
func (n *AssignStmt) nodeType() string  { return "AssignStmt" }
func (n *EmitStmt) nodeType() string    { return "EmitStmt" }
func (n *ExprStmt) nodeType() string    { return "ExprStmt" }
func (n *NumberLit) nodeType() string   { return "NumberLit" }
func (n *StringLit) nodeType() string   { return "StringLit" }
func (n *BoolLit) nodeType() string     { return "BoolLit" }
func (n *NullLit) nodeType() string     { return "NullLit" }
func (n *Ident) nodeType() string       { return "Ident" }
func (n *MemberExpr) nodeType() string  { return "MemberExpr" }
func (n *CallExpr) nodeType() string    { return "CallExpr" }
func (n *BinaryExpr) nodeType() string  { return "BinaryExpr" }
func (n *UnaryExpr) nodeType() string   { return "UnaryExpr" }
func (n *IsExpr) nodeType() string      { return "IsExpr" }
func (n *ObjectLit) nodeType() string   { return "ObjectLit" }
