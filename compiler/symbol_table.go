package compiler

type SymbolScope string

const (
    GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
    Name string
    Scope SymbolScope
    Index int
}

type SymbolTable struct {
    store map[string]Symbol
    numDefs int
}

func NewSymbolTable() *SymbolTable {
    s := make(map[string]Symbol)
    return &SymbolTable{store: s}
}

func (st *SymbolTable) Define(name string) Symbol {
    s := Symbol{Name: name, Scope: GlobalScope, Index: st.numDefs}
    st.store[name] = s
    st.numDefs++

    return s
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
    s, ok := st.store[name]
    return s, ok
}
