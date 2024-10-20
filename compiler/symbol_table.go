package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
	Outer          *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol), Outer: nil}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	scope := GlobalScope
	if s.Outer != nil {
		scope = LocalScope
	}
	symbol := Symbol{Name: name, Scope: scope, Index: s.numDefinitions}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	if !ok && s.Outer != nil {
		return s.Outer.Resolve(name)
	}
	return symbol, ok
}
