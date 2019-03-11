package vm

import (
    "fmt"
    "monkey_interpreter/object"
    "monkey_compiler/code"
    "monkey_compiler/compiler"
)

const StackSize = 2048

type VM struct{
    instructions code.Instructions
    constants []object.Object
    stack []object.Object
    sp int // always points to the next value. Top of stack is stack[sp-1]
}

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

func New(bytecode *compiler.Bytecode) *VM {
    vm := &VM{
        instructions: bytecode.Instructions,
        constants: bytecode.Constants,
        stack: make([]object.Object, StackSize),
        sp: 0,
    }

    return vm
}

func (vm *VM) Run() error {
    for ip := 0; ip < len(vm.instructions); ip++ {
        op := code.Opcode(vm.instructions[ip])

        switch op {
        case code.OpConst:
            constIndex := code.ReadUint16(vm.instructions[ip+1:])
            ip += 2

            err := vm.push(vm.constants[constIndex])
            if err != nil {
                return err
            }

        case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
            err := vm.executeBinaryOperation(op)
            if err != nil {
                return err
            }

        case code.OpTrue:
            err := vm.push(True)
            if err != nil {
                return err
            }

        case code.OpFalse:
            err := vm.push(False)
            if err != nil {
                return err
            }

        case code.OpEq, code.OpNE, code.OpGT:
            err := vm.executeCompatison(op)
            if err != nil {
                return err
            }

        case code.OpPop:
            vm.pop()
        }
    }

    return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
    r := vm.pop()
    l := vm.pop()
    ltype := l.Type()
    rtype := r.Type()

    if ltype == object.INTEGER_OBJ && rtype == object.INTEGER_OBJ {
        return vm.executeBinaryIntegerOperation(op, l, r)
    }

    return fmt.Errorf("invalid ltype or rtype")
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, l, r object.Object) error {
    lval := l.(*object.Integer).Value
    rval := r.(*object.Integer).Value

    var val int64
    switch op {
    case code.OpAdd:
        val = lval + rval
    case code.OpSub:
        val = lval - rval
    case code.OpMul:
        val = lval * rval
    case code.OpDiv:
        val = lval / rval
    default:
        return fmt.Errorf("invalid operator")
    }
    o := &object.Integer{Value: val}
    return vm.push(o)
}

func (vm *VM) executeCompatison(op code.Opcode) error {
    rExp := vm.pop()
    lExp := vm.pop()

    if lExp.Type() == object.INTEGER_OBJ || rExp.Type() == object.INTEGER_OBJ {
        return vm.executeIntegerComparison(op, lExp, rExp)
    }

    switch op {
    case code.OpEq:
        return vm.push(nativeBoolToBooleanObject(lExp == rExp))
    case code.OpNE:
        return vm.push(nativeBoolToBooleanObject(lExp != rExp))
    default:
        return fmt.Errorf("unknown operator")
    }
}

func (vm *VM) executeIntegerComparison(op code.Opcode, l, r object.Object) error {
    lval := l.(*object.Integer).Value
    rval := r.(*object.Integer).Value

    switch op {
    case code.OpEq:
        return vm.push(nativeBoolToBooleanObject(lval == rval))
    case code.OpNE:
        return vm.push(nativeBoolToBooleanObject(lval != rval))
    case code.OpGT:
        return vm.push(nativeBoolToBooleanObject(lval > rval))
    default:
        return fmt.Errorf("unknown operator")
    }
}

func nativeBoolToBooleanObject(b bool) *object.Boolean {
    if b {
        return True
    }
    return False
}

func (vm *VM) StackTop() object.Object {
    if vm.sp == 0 {
        return nil
    }
    return vm.stack[vm.sp-1]
}

func (vm *VM) push(ob object.Object) error {
    if vm.sp >= StackSize {
        return fmt.Errorf("stack overflow")
    }
    vm.stack[vm.sp] = ob
    vm.sp++
    return nil
}

func (vm *VM) pop() object.Object {
    if vm.sp < 1 {
        fmt.Errorf("cannot pop")
    }

    ob := vm.stack[vm.sp-1]
    vm.sp--
    return ob
}

func (vm *VM) LastPoppedStackElem() object.Object {
    return vm.stack[vm.sp]
}
