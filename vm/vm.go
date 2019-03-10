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

        case code.OpAdd:
            r := vm.pop()
            l := vm.pop()
            val := l.(*object.Integer).Value + r.(*object.Integer).Value
            o := &object.Integer{Value: val}
            vm.push(o)
            ip++

        case code.OpPop:
            vm.pop()
        }
    }

    return nil
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
