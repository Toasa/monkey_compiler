package vm

import (
    "fmt"
    "monkey_interpreter/object"
    "monkey_compiler/code"
    "monkey_compiler/compiler"
)

const StackSize = 2048
const GlobalsSize = 65536

type VM struct{
    instructions code.Instructions
    constants []object.Object
    stack []object.Object
    sp int // always points to the next value. Top of stack is stack[sp-1]
    globals []object.Object
}

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null =&object.Null{}

func New(bytecode *compiler.Bytecode) *VM {
    vm := &VM{
        instructions: bytecode.Instructions,
        constants: bytecode.Constants,
        stack: make([]object.Object, StackSize),
        sp: 0,
        globals: make([]object.Object, GlobalsSize),
    }

    return vm
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
    vm := New(bytecode)
    vm.globals = s
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

        case code.OpSetGlobal:
            globalIndex := code.ReadUint16(vm.instructions[ip+1:])
            ip += 2
            vm.globals[globalIndex] = vm.pop()

        case code.OpGetGlobal:
            globalIndex := code.ReadUint16(vm.instructions[ip+1:])
            ip += 2
            err := vm.push(vm.globals[globalIndex])
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

        case code.OpBang:
            err := vm.executeBangOperator()
            if err != nil {
                return nil
            }

        case code.OpMinus:
            err := vm.executeMinusOperator()
            if err != nil {
                return err
            }

        case code.OpEq, code.OpNE, code.OpGT:
            err := vm.executeComparison(op)
            if err != nil {
                return err
            }

        case code.OpPop:
            vm.pop()

        case code.OpJump:
            jumpDst := code.ReadUint16(vm.instructions[ip+1:])
            ip = int(jumpDst) - 1

        case code.OpJumpNotTruthy:
            jumpDst := code.ReadUint16(vm.instructions[ip+1:])
            ip += 2

            cond := vm.pop()
            if !isTruthy(cond) {
                ip = int(jumpDst) - 1
            }

        case code.OpArray:
            len := int(code.ReadUint16(vm.instructions[ip+1:]))
            ip += 2

            arr := vm.buildArray(vm.sp - len, vm.sp)
            vm.sp -= len

            err := vm.push(arr)
            if err != nil {
                return err
            }

        case code.OpHash:
            len := int(code.ReadUint16(vm.instructions[ip+1:]))
            ip += 2

            hash, err := vm.buildHash(vm.sp - len, vm.sp)
            if err != nil {
                return err
            }
            vm.sp -= len

            err = vm.push(hash)
            if err != nil {
                return err
            }

        case code.OpNull:
            err := vm.push(Null)
            if err != nil {
                return err
            }
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

    if ltype == object.STRING_OBJ && rtype == object.STRING_OBJ {
        return vm.executeBinaryStringOperation(op, l, r)
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

func (vm *VM) executeBinaryStringOperation(op code.Opcode, l, r object.Object) error {
    lval := l.(*object.String).Value
    rval := r.(*object.String).Value

    if op != code.OpAdd {
        return fmt.Errorf("String has only `+` operator")
    }
    o := &object.String{Value: lval + rval}
    return vm.push(o)
}

func (vm *VM) executeComparison(op code.Opcode) error {
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

func (vm *VM) executeBangOperator() error {
    operand := vm.pop()

    switch operand {
    case True:
        return vm.push(False)
    case False:
        return vm.push(True)
    case Null:
        return vm.push(True)
    default:
        return vm.push(False)
    }
}

func (vm *VM) executeMinusOperator() error {
    operand := vm.pop()
    if operand.Type() != object.INTEGER_OBJ {
        return fmt.Errorf("invalid operand")
    }

    value := operand.(*object.Integer).Value
    return vm.push(&object.Integer{Value: -value})
}

func nativeBoolToBooleanObject(b bool) *object.Boolean {
    if b {
        return True
    }
    return False
}

func (vm *VM) buildArray(startIndex int, endIndex int) object.Object {
    elems := make([]object.Object, endIndex - startIndex)

    for i := startIndex; i < endIndex; i++ {
        elems[i - startIndex] = vm.stack[i]
    }

    return &object.Array{Elems: elems}
}

func (vm *VM) buildHash(startIndex int, endIndex int) (object.Object, error) {
    hashedPairs := make(map[object.HashKey]object.HashPair)

    for i := startIndex; i < endIndex; i += 2 {
        key := vm.stack[i]
        val := vm.stack[i + 1]

        pair := object.HashPair{Key: key, Value: val}

        hashKey, ok := key.(object.Hashable)
        if !ok {
            return nil, fmt.Errorf("unusable as hash key")
        }

        hashedPairs[hashKey.HashKey()] = pair
    }

    return &object.Hash{Pairs: hashedPairs}, nil
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

func isTruthy(obj object.Object) bool {
    switch obj := obj.(type) {
    case *object.Boolean:
        return obj.Value
    case *object.Null:
        return false
    default:
        return true
    }
}
