package compiler

import (
    "monkey_interpreter/ast"
    "monkey_interpreter/object"
    "monkey_compiler/code"
)

type Compiler struct {
    insts code.Insts
    // consts pool
    consts []object.Object
}

func New() *Compiler {
    return &Compiler{
        insts: code.Insts{},
        consts: []object.Object{},
    }
}

func (c *Compiler) Compile(node ast.Node) error {
    return nil
}

func (c *Compiler) Bytecode() *Bytecode {
    return &Bytecode {
        Insts: c.insts,
        Consts: c.consts,
    }
}

type Bytecode struct {
    Insts code.Insts
    Consts []object.Object
}
