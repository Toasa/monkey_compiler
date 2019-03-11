package compiler

import (
    "fmt"
    "monkey_interpreter/ast"
    "monkey_interpreter/object"
    "monkey_compiler/code"
)

type Compiler struct {
    // operator + indices of operands(the index start without operator)
    instructions code.Instructions
    // constants pool
    constants []object.Object
}

type Bytecode struct {
    Instructions code.Instructions
    Constants []object.Object
}

func New() *Compiler {
    return &Compiler{
        instructions: code.Instructions{},
        constants: []object.Object{},
    }
}

func (c *Compiler) Compile(node ast.Node) error {
    switch node := node.(type) {
    case *ast.Program:
        for _, s := range node.Statements {
            err := c.Compile(s)
            if err != nil {
                return err
            }
        }

    case *ast.ExpressionStatement:
        err := c.Compile(node.Expression)
        if err != nil {
            return err
        }
        c.emit(code.OpPop)

    case *ast.InfixExpression:

        if node.Operator == "<" {
            err := c.Compile(node.Right)
            if err != nil {
                return err
            }
            err = c.Compile(node.Left)
            if err != nil {
                return err
            }
            c.emit(code.OpGT)
            return nil
        }

        err := c.Compile(node.Left)
        if err != nil {
            return err
        }

        err = c.Compile(node.Right)
        if err != nil {
            return err
        }

        switch node.Operator {
        case "+":
            c.emit(code.OpAdd)
        case "-":
            c.emit(code.OpSub)
        case "*":
            c.emit(code.OpMul)
        case "/":
            c.emit(code.OpDiv)
        case ">":
            c.emit(code.OpGT)
        case "==":
            c.emit(code.OpEq)
        case "!=":
            c.emit(code.OpNE)
        default:
            return fmt.Errorf("unknown operator")
        }

    case *ast.PrefixExpression:
        err := c.Compile(node.Right)
        if err != nil {
            return err
        }
        
        switch node.Operator {
        case "!":
            c.emit(code.OpBang)
        case "-":
            c.emit(code.OpMinus)
        default:
            return fmt.Errorf("unknown operator")
        }

    case *ast.IntegerLiteral:
        integer := &object.Integer{Value: node.Value}
        c.emit(code.OpConst, c.addConstant(integer))

    case *ast.Boolean:
        var opc code.Opcode
        if node.Value {
            opc = code.OpTrue
        } else {
            opc = code.OpFalse
        }
        c.emit(opc)
    }
    return nil
}

func (c *Compiler) Bytecode() *Bytecode {
    return &Bytecode {
        Instructions: c.instructions,
        Constants: c.constants,
    }
}

func (c *Compiler) addConstant(obj object.Object) int {
    c.constants = append(c.constants, obj)
    return len(c.constants) - 1
}

// generate an instruction and add it to the result
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
    // operand... are each index of it
    ins := code.Make(op, operands...)
    pos := c.addInstruction(ins)
    return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
    posNewInstruction := len(c.instructions)
    c.instructions = append(c.instructions, ins...)
    return posNewInstruction
}
