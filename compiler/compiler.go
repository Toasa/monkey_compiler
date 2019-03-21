package compiler

import (
    "fmt"
    "monkey_interpreter/ast"
    "monkey_interpreter/object"
    "monkey_compiler/code"
)

type Compiler struct {
    // kinds of instructions
    // * operator(OpConst: 1byte) + indices of operands(the index(2byte) start without operator)
    // * operator(OpJump: 1byte) + dstAddress(2byte)
    // * operator(OpGetGlobal: 1byte) + index of operand(2byte)
    // * operator(otherwise: 1byte)
    instructions code.Instructions

    // constants pool
    constants []object.Object

    lastInstruction EmitedInstruction
    prevInstruction EmitedInstruction

    symbolTable *SymbolTable
}

type Bytecode struct {
    Instructions code.Instructions
    Constants []object.Object
}

type EmitedInstruction struct {
    Opcode code.Opcode
    Position int
}

func New() *Compiler {
    return &Compiler{
        instructions: code.Instructions{},
        constants: []object.Object{},
        lastInstruction: EmitedInstruction{},
        prevInstruction: EmitedInstruction{},
        symbolTable: NewSymbolTable(),
    }
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
    compiler := New()
    compiler.symbolTable = s
    compiler.constants = constants
    return compiler
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

    case *ast.LetStatement:
        err := c.Compile(node.Value)
        if err != nil {
            return err
        }
        symbol := c.symbolTable.Define(node.Name.Value)
        c.emit(code.OpSetGlobal, symbol.Index)

    case *ast.ExpressionStatement:
        err := c.Compile(node.Expression)
        if err != nil {
            return err
        }
        // Pop are needed immediately following the expressionstatement(es).
        // Because the value that es produce is not reused by definition.
        // (reuse means push a value on the stack)
        c.emit(code.OpPop)

    case *ast.BlockStatement:
        for _, stmt := range node.Statements {
            err := c.Compile(stmt)
            if err != nil {
                return nil
            }
        }

    case *ast.IntegerLiteral:
        integer := &object.Integer{Value: node.Value}
        c.emit(code.OpConst, c.addConstant(integer))

    case *ast.Identifier:
        symbol, ok := c.symbolTable.Resolve(node.Value)
        if !ok {
            return fmt.Errorf("undefined variable %s", node.Value)
        }
        c.emit(code.OpGetGlobal, symbol.Index)

    case *ast.Boolean:
        var opc code.Opcode
        if node.Value {
            opc = code.OpTrue
        } else {
            opc = code.OpFalse
        }
        c.emit(opc)

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

    case *ast.IfExpression:
        err := c.Compile(node.Cond)
        if err != nil {
            return err
        }

        jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

        err = c.Compile(node.Cons)
        if err != nil {
            return err
        }
        if c.lastInstructionIsPop() {
            c.removeLastPop()
        }

        jumpPos := c.emit(code.OpJump, 9999)

        afterConsPos := len(c.instructions)
        c.changeOperand(jumpNotTruthyPos, afterConsPos)

        if node.Alt == nil {
            c.emit(code.OpNull)
        } else {
            err = c.Compile(node.Alt)
            if err != nil {
                return err
            }
            if c.lastInstructionIsPop() {
                c.removeLastPop()
            }
        }

        afterAltPos := len(c.instructions)
        c.changeOperand(jumpPos, afterAltPos)
    }
    return nil
}

func (c *Compiler) Bytecode() *Bytecode {
    return &Bytecode {
        Instructions: c.instructions,
        Constants: c.constants,
    }
}

// Generate an instruction and add it to the result.
// Return value is the index of new added instruction.
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
    // operand... are each index of it
    ins := code.Make(op, operands...)
    pos := c.addInstruction(ins)
    c.setLastInstruction(op, pos)

    return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
    last := EmitedInstruction{Opcode: op, Position: pos}
    c.prevInstruction = c.lastInstruction
    c.lastInstruction = last
}

func (c *Compiler) addInstruction(ins []byte) int {
    posNewInstruction := len(c.instructions)
    c.instructions = append(c.instructions, ins...)
    return posNewInstruction
}

func (c *Compiler) addConstant(obj object.Object) int {
    c.constants = append(c.constants, obj)
    return len(c.constants) - 1
}

func (c *Compiler) lastInstructionIsPop() bool {
    return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
    c.instructions = c.instructions[: c.lastInstruction.Position]
    c.lastInstruction = c.prevInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
    for i := 0; i < len(newInstruction); i++ {
        c.instructions[pos + i] = newInstruction[i]
    }
}

func (c *Compiler) changeOperand(opPos int, operand int) {
    op := code.Opcode(c.instructions[opPos])
    newInstruction := code.Make(op, operand)
    c.replaceInstruction(opPos, newInstruction)
}
